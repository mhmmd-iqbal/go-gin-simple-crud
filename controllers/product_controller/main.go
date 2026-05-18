package product_controller

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mhmmd-iqbal/go-rest-gin/models"
	"github.com/mhmmd-iqbal/go-rest-gin/requests/product_request"
	"github.com/mhmmd-iqbal/go-rest-gin/services/product_service"
)

func Index(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	search := c.Query("search")

	if page < 1 {
		page = 1
	}

	if limit < 1 {
		limit = 10
	}

	products, total, err := product_service.GetProducts(search, page, limit)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var productResponses []models.ListProductResponse

	for _, product := range products {
		productResponses = append(productResponses, models.ListProductResponse{
			SKU:   product.SKU,
			Name:  product.Name,
			Price: product.Price,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Products retrieved successfully",
		"data":    productResponses,
		"meta": gin.H{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

func Show(c *gin.Context) {
	sku := c.Param("sku")

	if sku == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "SKU is required"})
		return
	}

	product, err := product_service.GetDetailProduct(sku)

	if err != nil {
		if errors.Is(err, errors.New("product not found")) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var productResponse = models.DetailProductResponse{
		SKU:         product.SKU,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
	}

	c.JSON(http.StatusOK, gin.H{"data": productResponse, "message": "Product retrieved successfully"})
}

func Create(c *gin.Context) {
	var input product_request.CreateProductRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.Name == "" || input.SKU == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name and SKU are required"})
		return
	}

	if input.Price <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "price must be greater than zero"})
		return
	}

	product, err := product_service.CreateProduct(input)

	if err != nil {
		if errors.Is(err, errors.New("product already exists")) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Product with the SKU : " + input.SKU + " already exists"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var productResponse = models.DetailProductResponse{
		SKU:         product.SKU,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
	}

	c.JSON(http.StatusCreated, gin.H{
		"data":    productResponse,
		"message": "Product created successfully",
	})
}

func Update(c *gin.Context) {
	var input product_request.UpdateProductRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sku := c.Param("sku")

	if sku == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "SKU is required"})
		return
	}

	product, err := product_service.UpdateProduct(sku, input)

	if err != nil {
		if errors.Is(err, errors.New("product not found")) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var productResponse = models.DetailProductResponse{
		SKU:         product.SKU,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
	}

	c.JSON(http.StatusOK, gin.H{"data": productResponse, "message": "Product updated successfully"})
}

func Delete(c *gin.Context) {
	sku := c.Param("sku")

	if sku == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "SKU is required"})
		return
	}

	if err := product_service.DeleteProduct(sku); err != nil {
		if errors.Is(err, errors.New("product not found")) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
}

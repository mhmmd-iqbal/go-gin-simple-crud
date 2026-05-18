package product_controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mhmmd-iqbal/go-rest-gin/models"
	"github.com/mhmmd-iqbal/go-rest-gin/repositories/product_repositories"
	"github.com/mhmmd-iqbal/go-rest-gin/requests/product_request"
)

func Index(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	search := c.Query("search")

	// Prevent invalid values
	if page < 1 {
		page = 1
	}

	if limit < 1 {
		limit = 10
	}

	products, total, err := product_repositories.GetProducts(search, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
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

	product, err := product_repositories.GetProductBySKU(sku)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
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

	_, err := product_repositories.GetProductBySKU(input.SKU)

	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Product with the SKU : " + input.SKU + " already exists"})
		return
	}

	mappingProduct := models.Product{
		Name:        input.Name,
		Description: input.Description,
		Price:       input.Price,
		SKU:         input.SKU,
	}

	product, err := product_repositories.CreateProduct(&mappingProduct)

	if err != nil {
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

	product, err := product_repositories.GetProductBySKU(sku)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	if input.Name != nil {
		product.Name = *input.Name
	}

	if input.Description != nil {
		product.Description = *input.Description
	}

	if input.Price != nil {
		product.Price = *input.Price
	}

	product, err = product_repositories.UpdateProductBySKU(sku, product)

	if err != nil {
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

	_, err := product_repositories.GetProductBySKU(sku)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	if err := product_repositories.DeleteProductBySKU(sku); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
}

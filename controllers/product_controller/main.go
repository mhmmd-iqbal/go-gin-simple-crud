package product_controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mhmmd-iqbal/go-rest-gin/models"
	"github.com/mhmmd-iqbal/go-rest-gin/requests/product_request"
)

func Index(c *gin.Context) {
	var products []models.Product
	var total int64

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

	offset := (page - 1) * limit

	query := models.DB.Model(&models.Product{})

	// Search by name or sku
	if search != "" {
		query = query.Where(
			"name LIKE ? OR sku LIKE ?",
			"%"+search+"%",
			"%"+search+"%",
		)
	}

	if err := query.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Get paginated data
	if err := query.
		Limit(limit).
		Offset(offset).
		Order("id DESC").
		Find(&products).Error; err != nil {

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
	var Product models.Product
	sku := c.Param("sku")

	if sku == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "SKU is required"})
		return
	}

	if err := models.DB.Where("sku = ?", sku).First(&Product).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	var productResponse = models.DetailProductResponse{
		SKU:         Product.SKU,
		Name:        Product.Name,
		Description: Product.Description,
		Price:       Product.Price,
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

	var Product models.Product

	if err := models.DB.Where("sku = ?", input.SKU).First(&Product).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Product with the SKU : " + input.SKU + " already exists"})
		return
	}

	Product.Name = input.Name
	Product.Description = input.Description
	Product.Price = input.Price
	Product.SKU = input.SKU

	if err := models.DB.Create(&Product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var productResponse = models.DetailProductResponse{
		SKU:         Product.SKU,
		Name:        Product.Name,
		Description: Product.Description,
		Price:       Product.Price,
	}

	c.JSON(http.StatusCreated, map[string]any{
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

	var Product models.Product

	if err := models.DB.Where("sku = ?", sku).First(&Product).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	if input.Name != "" {
		Product.Name = input.Name
	}

	if input.Description != "" {
		Product.Description = input.Description
	}

	if input.Price > 0 {
		Product.Price = input.Price
	}

	if err := models.DB.Updates(&Product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var productResponse = models.DetailProductResponse{
		SKU:         Product.SKU,
		Name:        Product.Name,
		Description: Product.Description,
		Price:       Product.Price,
	}

	c.JSON(http.StatusOK, gin.H{"data": productResponse, "message": "Product updated successfully"})
}

func Delete(c *gin.Context) {
	var product models.Product
	sku := c.Param("sku")

	if sku == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "SKU is required"})
		return
	}

	if err := models.DB.Where("sku = ?", sku).First(&product).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	if err := models.DB.Delete(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
}

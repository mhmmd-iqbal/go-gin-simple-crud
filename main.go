package main

import (
	"github.com/gin-gonic/gin"
	"github.com/mhmmd-iqbal/go-rest-gin/models"
	"github.com/mhmmd-iqbal/go-rest-gin/routes/product_route"
)

func main() {
	r := gin.Default()
	models.ConnectDatabase()

	// Routes
	// r.GET("/products", product_controller.Index)
	// r.GET("/products/:sku", product_controller.Show)
	// r.POST("/products", product_controller.Create)
	// r.PUT("/products/:sku", product_controller.Update)
	// r.DELETE("/products/:sku", product_controller.Delete)

	// Register product routes
	product_route.RegisterProductRoutes(r)

	// Start server
	r.Run(":5555")
}

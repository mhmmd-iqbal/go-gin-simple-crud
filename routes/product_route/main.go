package product_route

import (
	"github.com/gin-gonic/gin"
	"github.com/mhmmd-iqbal/go-rest-gin/controllers/product_controller"
)

func RegisterProductRoutes(r *gin.Engine) {
	r.GET("/products", product_controller.Index)
	r.GET("/products/:sku", product_controller.Show)
	r.POST("/products", product_controller.Create)
	r.PUT("/products/:sku", product_controller.Update)
	r.DELETE("/products/:sku", product_controller.Delete)
}

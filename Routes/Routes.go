package Routes

import (
	"ecommerce_api/Controllers"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	grp1 := r.Group("/customer")

	{
		grp1.GET("/products", Controllers.GetProducts)
		grp1.GET("/:id/cart", Controllers.GetCartByID)
		grp1.PUT("/:id/cart", Controllers.ModifyCart)
		grp1.POST("/:id/order", Controllers.CreateOrder)
		grp1.GET("/:id/orders", Controllers.GetAllOrders)
		grp1.POST("/:id/checkout", Controllers.CheckOut)
	}
	grp2 := r.Group("/merchant")
	{
		grp2.GET("/products", Controllers.GetProducts)
		grp2.GET("/orders", Controllers.GetOrdersByStatus)

	}
	return r
}

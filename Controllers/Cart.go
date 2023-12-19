package Controllers

import (
	"ecommerce_api/Models"
	"encoding/json"
	"fmt"

	"net/http"

	"github.com/gin-gonic/gin"
)

type Stock struct {
	Id       int `json:"id"`
	Quantity int `json:"quantity"`
}
type CartInfo struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Quantity    int    `json:"quantity"`
}

// Get Cart Products by Id
func GetCartByID(c *gin.Context) {
	id := c.Params.ByName("id")
	var cart Models.Cart
	err := Models.GetCartById(&cart, id)
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
	} else {
		var list []int
		var items []Stock
		json.Unmarshal(cart.Products, &items)

		for _, product := range items {
			list = append(list, product.Id)
		}
		var product []Models.Product
		err := Models.GetProductsById(&product, list)
		if err != nil {
			c.AbortWithStatus(http.StatusNotFound)
		}
		CartProducts := make(map[uint]Models.Product)
		var total float64
		for _, item := range product {
			CartProducts[item.Id] = item
			total = total + item.Price
		}
		var ProductInfo []CartInfo
		for _, ele := range items {
			ProductInfo = append(ProductInfo, CartInfo{
				Id:          ele.Id,
				Name:        CartProducts[uint(ele.Id)].Name,
				Description: CartProducts[uint(ele.Id)].Description,
				Quantity:    ele.Quantity,
			})
		}
		c.JSON(http.StatusOK, ProductInfo)
	}
}
func ModifyCart(c *gin.Context) {
	var cart Models.Cart
	id := c.Params.ByName("id")
	err := Models.GetCartById(&cart, id)
	if err != nil {
		c.JSON(http.StatusNotFound, cart)
	}
	fmt.Print(cart)
	c.BindJSON(&cart)

	err = Models.UpdateCart(&cart, id)
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
	} else {
		c.JSON(http.StatusOK, cart)
	}

}

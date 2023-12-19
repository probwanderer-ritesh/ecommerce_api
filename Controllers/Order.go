package Controllers

import (
	"ecommerce_api/Models"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	razorpay "github.com/razorpay/razorpay-go"
	"gorm.io/datatypes"
)

type OrderInfo struct {
	CustomerId      string         `json:"customer_id"`
	ShippingAddress string         `json:"shipping_address"`
	BillingAddress  string         `json:"billing_address"`
	Products        datatypes.JSON `json:"products"`
	Mobile          string         `json:"mobile"`
}

func GetAllOrders(c *gin.Context) {
	var orders []Models.Order
	err := Models.GetAllOrders(&orders)
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
	} else {
		c.JSON(http.StatusOK, orders)
	}
}
func GetOrdersByStatus(c *gin.Context) {
	var orders []Models.Order
	status := c.Query("status")
	err := Models.GetOrdersByStatus(&orders, status)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusNotFound)
	} else {
		c.JSON(http.StatusOK, orders)
	}
}

// Create Order
func CreateOrder(c *gin.Context) {
	id := c.Params.ByName("id")
	var list []int
	var items []Stock
	var orderInfo OrderInfo
	c.BindJSON(&orderInfo)
	json.Unmarshal(orderInfo.Products, &items)

	for _, product := range items {
		list = append(list, product.Id)
	}
	var product []Models.Product
	err := Models.GetProductsById(&product, list)
	productMap := make(map[string]Models.Product)
	for _, ele := range product {
		x := fmt.Sprint(ele.Id)
		productMap[x] = ele
	}
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
	}
	var total float64
	for _, item := range items {
		x := fmt.Sprint(item.Id)
		total = total + productMap[x].Price*float64(item.Quantity)
	}

	client := razorpay.NewClient("rzp_test_5Cgy7HOjqS4cSC", "PLiAaCmjrhwFlB39ReKoZMVE")
	data := map[string]interface{}{
		"amount":   int(total * 100),
		"currency": "INR",
		"notes": map[string]interface{}{
			"mobile": orderInfo.Mobile,
		},
	}
	body, err := client.Order.Create(data, nil)

	if err != nil {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusNotFound)
	}
	CustomerId, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusNotFound)
	}

	newOrder := &Models.Order{
		Id:              body["id"].(string),
		CustomerId:      uint(CustomerId),
		ShippingAddress: orderInfo.ShippingAddress,
		BillingAddress:  orderInfo.BillingAddress,
		Status:          "pending",
		Products:        orderInfo.Products,
		Mobile:          body["notes"].(map[string]interface{})["mobile"].(string),
		Total:           total,
	}
	err = Models.CreateOrder(newOrder)
	if err != nil {
		fmt.Println(err.Error())
		c.AbortWithStatus(http.StatusNotFound)
	} else {
		c.JSON(http.StatusOK, newOrder)
	}
}

type Payment struct {
	PaymentId        string `json:"razorpay_payment_id"`
	OrderId          string `json:"razorpay_order_id"`
	PaymentSignature string `json:"razorpay_signature"`
}

func CheckOut(c *gin.Context) {
	var paymentInfo Payment
	c.BindJSON(&paymentInfo)
	err := Models.UpdatePayment(paymentInfo.OrderId, paymentInfo.PaymentId)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusNotFound)
	}

	c.JSON(http.StatusOK, paymentInfo.OrderId)
}

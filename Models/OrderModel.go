package Models

import (
	_ "github.com/go-sql-driver/mysql"
	"gorm.io/datatypes"
)

type Order struct {
	Id              string         `json:"id"`
	CustomerId      uint           `json:"customer_id"`
	ShippingAddress string         `json:"shipping_address"`
	BillingAddress  string         `json:"billing_address"`
	Status          string         `json:"status"`
	PaymentId       string         `json:"payment_id"`
	Products        datatypes.JSON `json:"products"`
	Mobile          string         `json:"mobile"`
	Total           float64        `json:"total"`
}

func (o *Order) TableName() string {
	return "orders"
}

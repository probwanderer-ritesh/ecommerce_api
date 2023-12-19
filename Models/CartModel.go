package Models

import (
	_ "github.com/go-sql-driver/mysql"
	"gorm.io/datatypes"
)

type Cart struct {
	Id       uint           `json:"id"`
	Products datatypes.JSON `json:"products"`
}

func (c *Cart) TableName() string {
	return "cart"
}

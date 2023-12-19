package Models

import (
	"ecommerce_api/Config"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func GetCartById(c *Cart, id string) (err error) {
	if err = Config.DB.Where("id=?", id).First(c).Error; err != nil {
		return err
	}

	return nil
}

func UpdateCart(c *Cart, id string) (err error) {
	fmt.Println(c)
	if err = Config.DB.Save(c).Error; err != nil {
		return err
	}
	return nil
}

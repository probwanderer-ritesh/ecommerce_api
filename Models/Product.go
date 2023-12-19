package Models

import (
	"ecommerce_api/Config"

	_ "github.com/go-sql-driver/mysql"
)

func GetAllProducts(product *[]Product) (err error) {
	if err = Config.DB.Find(product).Error; err != nil {
		return err
	}
	return nil
}
func GetProductsById(product *[]Product, ids []int) (err error) {
	if err = Config.DB.Find(&product, ids).Error; err != nil {
		return err
	}
	return nil
}

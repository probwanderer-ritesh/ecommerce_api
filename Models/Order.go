package Models

import (
	"ecommerce_api/Config"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

type Stock struct {
	Id       int `json:"id"`
	Quantity int `json:"quantity"`
}

func GetAllOrders(order *[]Order) (err error) {
	if err = Config.DB.Find(order).Error; err != nil {
		return err
	}
	return nil
}

func GetOrdersByStatus(order *[]Order, filterValue string) (err error) {
	if err = Config.DB.Find(order, "status=?", filterValue).Error; err != nil {
		return err
	}
	return nil
}

func CreateOrder(order *Order) (err error) {
	// if err = Config.DB.Create(order).Error; err != nil {
	// 	return err
	// }
	// return nil
	Config.DB.Transaction(func(tx *gorm.DB) error {

		var products []Stock
		json.Unmarshal(order.Products, &products)
		if err := tx.Create(order).Error; err != nil {
			return err
		}
		for _, item := range products {
			product := &Product{}
			if err = tx.Find(&product, "id=?", item.Id).Error; err != nil {
				return err
			}
			product.Quantity = product.Quantity - item.Quantity
			if err = tx.Save(product).Error; err != nil {
				return err
			}

		}
		return nil
	})
	go func() {

		select {
		case <-time.After(5 * 60 * time.Second):
			Config.DB.Transaction(func(tx *gorm.DB) error {
				fmt.Print("rolled back")
				order_new := &Order{}
				if err = tx.Find(&order_new, "id=?", order.Id).Error; err != nil {
					fmt.Println(err)
					return err
				}
				var products []Stock
				json.Unmarshal(order.Products, &products)
				fmt.Println(order_new.Status)
				if order_new.Status != "processed" {
					fmt.Println("not processed")
					order_new.PaymentId = ""
					order_new.Status = "failed"
					tx.Save(order_new)
					for _, item := range products {
						product := &Product{}
						if err = tx.Find(&product, "id=?", item.Id).Error; err != nil {
							fmt.Println(err)
							return err
						}
						product.Quantity = product.Quantity + item.Quantity
						if err = tx.Save(product).Error; err != nil {
							fmt.Println(err)
							return err
						}

					}
					return nil
				}
				return errors.New("transaction rollback fail")
			})
		}

	}()
	return nil
}
func UpdatePayment(orderId string, paymentId string) (err error) {
	order := &Order{}
	if err = Config.DB.Find(&order, "id=?", orderId).Error; err != nil {
		return err
	}
	order.PaymentId = paymentId
	order.Status = "processed"

	Config.DB.Save(&order)
	return nil
}

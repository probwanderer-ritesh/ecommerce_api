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

var ch chan string = make(chan string, 100)

func CreateOrder(order *Order) (err error) {
	// if err = Config.DB.Create(order).Error; err != nil {
	// 	return err
	// }
	// return nil
	if err := Config.DB.Create(order).Error; err != nil {
		return err
	}
	go func() {
		Config.DB.Transaction(func(tx *gorm.DB) error {

			var products []Stock
			json.Unmarshal(order.Products, &products)

			for _, item := range products {
				product := &Product{}
				if err = tx.Find(&product, "id=?", item.Id).Error; err != nil {
					return err
				}
				product.Quantity = product.Quantity - item.Quantity
				if err = tx.Save(product).Error; err != nil {
					return nil
				}

			}
			select {
			case x := <-ch:
				fmt.Println(x)
				return nil
			case <-time.After(5 * 60 * time.Second):
				FailedOrder := &Order{}
				if err = Config.DB.Find(&FailedOrder, "id=?", order.Id).Error; err != nil {
					return err
				}
				FailedOrder.PaymentId = "0"
				FailedOrder.Status = "failed"
				Config.DB.Save(&FailedOrder)
				return errors.New("transaction failed")
			}

		})
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
	ch <- orderId
	Config.DB.Save(&order)
	return nil
}

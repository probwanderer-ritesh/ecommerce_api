package main

import (
	"ecommerce_api/Config"
	"ecommerce_api/Models"
	"ecommerce_api/Routes"
	"fmt"

	"github.com/jinzhu/gorm"
)

var err error

func main() {
	Config.DB, err = gorm.Open("mysql",
		Config.DbURL(Config.BuildDBConfig()))
	if err != nil {
		fmt.Println("Status:", err)
	}
	defer Config.DB.Close()
	Config.DB.AutoMigrate(&Models.Product{}, &Models.Cart{}, &Models.Order{})
	r := Routes.SetupRouter()

	r.Run()
}

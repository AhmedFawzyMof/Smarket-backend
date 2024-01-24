package models

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"sync"

	"alwadi_markets/middleware"
)

type OrderProducts struct {
	Product     int `json:"id"`
	Quantity    int `json:"quantity"`
	ProductType int `json:"typeId"`
	Id          int
	Order       string
	Products    []Product
}

type OP struct {
	Quantity int
	Order    string
	Name     string
	Image    string
	Portion  int
	Price    int
	Uint     string
}

func (op OrderProducts) Add(db *sql.DB) bool {
	_, err := db.Exec("INSERT INTO `OrderProducts`(`product`, `producttype`, `order`, `quantity`) VALUES (?, ?, ?, ?)", op.Product, op.ProductType, op.Order, op.Quantity)

	if err != nil {
		return false
	}

	return true
}

func (op OrderProducts) GetByOrder(db *sql.DB, response chan []byte, wg *sync.WaitGroup) {
	var Products []OP

	OrderProduct, err := db.Prepare("SELECT OrderProducts.quantity, OrderProducts.`order`, Products.`name`, Products.image, producttype.portion, producttype.price, producttype.uint FROM OrderProducts INNER JOIN Products ON OrderProducts.product=Products.`id` INNER JOIN producttype ON OrderProducts.producttype=producttype.id WHERE OrderProducts.`order` = ?")

	if err != nil {
		panic(err.Error())
	}

	rows, err := OrderProduct.Query(op.Order)
	if err != nil {
		panic(err.Error())
	}

	defer rows.Close()

	for rows.Next() {
		var Product OP

		if err := rows.Scan(&Product.Quantity, &Product.Order, &Product.Name, &Product.Image, &Product.Portion, &Product.Price, &Product.Uint); err != nil {
			panic(err.Error())
		}

		Products = append(Products, Product)
	}

	res, err := json.Marshal(Products)

	if err != nil {
		fmt.Println(err)
	}

	response <- res
	wg.Done()
}

func (op OrderProducts) Delete(db *sql.DB, response chan []byte, wg *sync.WaitGroup) {
	defer wg.Done()
	Response := make(map[string]interface{})
	_, err := db.Exec("DELETE FROM OrderProducts WHERE id = ?", op.Id)

	if err != nil {
		Response["Error"] = true

		middleware.SendResponse(response, Response)
		return
	}

	Response["Error"] = false

	middleware.SendResponse(response, Response)
}

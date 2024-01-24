package models

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"alwadi_markets/middleware"

	"github.com/google/uuid"
)

type Orders struct {
	Id        string
	Date      time.Time
	User      string
	Delivered int
	Paid      int
	Method    string
	Confirmed int
}

func (o Orders) Add(db *sql.DB) (string, error) {
	id := uuid.New()
	o.Id = fmt.Sprintf("%v", id)
	o.Date = time.Now()
	_, err := db.Exec("INSERT INTO `Orders`(`id`, `date`, `user`, `delivered`, `paid`, `method`, `confirmed`) VALUES (?, ?, ?, ?, ?, ?, ?)", o.Id, o.Date, o.User, o.Delivered, o.Paid, o.Method, o.Confirmed)
	if err != nil {
		return "", err
	}

	return o.Id, nil
}

func (o Orders) GetForUser(db *sql.DB, response chan []byte, wg *sync.WaitGroup) {
	defer wg.Done()
	var SliceOrder []Orders
	TOrders, err := db.Query("SELECT * FROM `Orders` WHERE user = ?", o.User)

	if err != nil {
		fmt.Println(err.Error())
	}

	defer TOrders.Close()

	for TOrders.Next() {
		var Order Orders

		if err := TOrders.Scan(&Order.Id, &Order.Date, &Order.User, &Order.Delivered, &Order.Paid, &Order.Method, &Order.Confirmed); err != nil {
			fmt.Println(err.Error())
		}

		SliceOrder = append(SliceOrder, Order)
	}

	res, err := json.Marshal(SliceOrder)

	if err != nil {
		fmt.Println(err)
	}

	response <- res
}

func (o Orders) OrderDitails(db *sql.DB, response chan []byte, wg *sync.WaitGroup) {
	defer wg.Done()
	Response := make(map[string]interface{})
	TOrders := db.QueryRow("SELECT * FROM `Orders` WHERE id = ?", o.Id)

	var Order Orders

	if err := TOrders.Scan(&Order.Id, &Order.Date, &Order.User, &Order.Delivered, &Order.Paid, &Order.Method, &Order.Confirmed); err != nil {
		fmt.Println(err.Error())
	}

	FindUser := db.QueryRow("SELECT username, email, phone, spare_phone FROM Users WHERE id = ?", Order.User)

	var User Users

	err := FindUser.Scan(&User.Username, &User.Email, &User.Phone, &User.Spare_phone)

	if err != nil {
		Response["Error"] = true

		middleware.SendResponse(response, Response)
		return
	}

	var Products []OP

	OrderProduct, err := db.Prepare("SELECT OrderProducts.quantity, OrderProducts.`order`, Products.`name`, Products.image, producttype.portion, producttype.price, producttype.uint FROM OrderProducts INNER JOIN Products ON OrderProducts.product=Products.`id` INNER JOIN producttype ON OrderProducts.producttype=producttype.id WHERE OrderProducts.`order` = ?")

	if err != nil {
		panic(err.Error())
	}

	rows, err := OrderProduct.Query(o.Id)
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

	Response["Order"] = Order
	Response["User"] = User
	Response["Products"] = Products

	middleware.SendResponse(response, Response)
}

func (o Orders) Delete(db *sql.DB, response chan []byte, wg *sync.WaitGroup) {
	defer wg.Done()
	Response := make(map[string]interface{})

	Confirmed := db.QueryRow("SELECT confirmed FROM Orders WHERE id = ?", o.Id)

	Err := Confirmed.Scan(&o.Confirmed)
	if Err == nil {

		if o.Confirmed == 0 {
			_, err := db.Exec("DELETE FROM Orders WHERE (id, user)=(?, ?)", o.Id, o.User)

			if err != nil {
				Response["Error"] = true
				middleware.SendResponse(response, Response)
				return
			}

			Response["Error"] = false
			middleware.SendResponse(response, Response)
			return
		}
	}

	Response["Error"] = true

	middleware.SendResponse(response, Response)
}

func (o Orders) GetAll(db *sql.DB, response chan []byte, wg *sync.WaitGroup) {
	defer wg.Done()
	var SliceOrder []Orders
	TOrders, err := db.Query("SELECT * FROM `Orders`")

	if err != nil {
		fmt.Println(err.Error())
	}

	defer TOrders.Close()

	for TOrders.Next() {
		var Order Orders

		if err := TOrders.Scan(&Order.Id, &Order.Date, &Order.User, &Order.Delivered, &Order.Paid, &Order.Method, &Order.Confirmed); err != nil {
			fmt.Println(err.Error())
		}

		SliceOrder = append(SliceOrder, Order)
	}
	res, err := json.Marshal(SliceOrder)

	if err != nil {
		fmt.Println(err)
	}

	response <- res
}

func (o Orders) Update(db *sql.DB, response chan []byte, wg *sync.WaitGroup) {
	defer wg.Done()
	Response := make(map[string]interface{})

	if o.Paid == 1 {
		_, err := db.Exec("UPDATE Orders set paid=? WHERE id =  ?", o.Paid, o.Id)
		if err != nil {
			Response["Error"] = true
			middleware.SendResponse(response, Response)
			return
		}

		Response["Error"] = false
		middleware.SendResponse(response, Response)
		return
	}
	if o.Delivered == 1 {
		_, err := db.Exec("UPDATE Orders set delivered=? WHERE id =  ?", o.Delivered, o.Id)
		if err != nil {
			Response["Error"] = true
			middleware.SendResponse(response, Response)
			return
		}

		Response["Error"] = false
		middleware.SendResponse(response, Response)
		return
	}
	if o.Confirmed == 1 {
		_, err := db.Exec("UPDATE Orders set confirmed=? WHERE id =  ?", o.Confirmed, o.Id)
		if err != nil {
			Response["Error"] = true
			middleware.SendResponse(response, Response)
			return
		}

		Response["Error"] = false
		middleware.SendResponse(response, Response)
		return
	}

	Response["Error"] = true
	middleware.SendResponse(response, Response)
}

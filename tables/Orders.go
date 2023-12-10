package tables

import (
	"alwadi_markets/middleware"
	"database/sql"
	"encoding/json"
	"fmt"
	"sync"
	"time"

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

func (o Orders) Get(db *sql.DB, response chan []byte, wg *sync.WaitGroup) {
	var SOrders []Orders
	TOrders, err := db.Query("SELECT * FROM `Orders`")

	if err != nil {
		panic(err.Error())
	}

	defer TOrders.Close()

	for TOrders.Next() {
		var Order Orders

		if err := TOrders.Scan(&Order.Id, &Order.Date, &Order.User, &Order.Delivered, &Order.Paid, &Order.Method, &Order.Confirmed); err != nil {
			panic(err.Error())
		}

		SOrders = append(SOrders, Order)
	}

	res, err := json.Marshal(SOrders)

	if err != nil {
		fmt.Println(err)
	}

	response <- res
	wg.Done()
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
	TOrders := db.QueryRow("SELECT * FROM `Orders` WHERE user = ?", o.Id)

	var Order Orders

	if err := TOrders.Scan(&Order.Id, &Order.Date, &Order.User, &Order.Delivered, &Order.Paid, &Order.Method, &Order.Confirmed); err != nil {
		fmt.Println(err.Error())
	}

	FindUser := db.QueryRow("SELECT username, email FROM Users WHERE id = ?", o.User)
	var User Users

	err := FindUser.Scan(&User.Username, &User.Email)

	if err != nil {
		Response["Error"] = true
		Response["Message"] = "لا يمكن تسجيل الدخول ببيانات المقدمة"

		middleware.SendResponse(response, Response)
		return
	}

	var Products []OrderProducts
	OP, err := db.Query("SELECT * FROM `OrderProducts` WHERE `order` = ?", o.Id)

	if err != nil {
		fmt.Println(err.Error())
	}

	defer OP.Close()

	for OP.Next() {
		var OrderProducts OrderProducts

		if err := OP.Scan(&OrderProducts.Id, &OrderProducts.Product, &OrderProducts.ProductType, &OrderProducts.Order, &OrderProducts.Quantity); err != nil {
			fmt.Println(err.Error())
		}

		Products = append(Products, OrderProducts)
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

	if Err != nil {

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
	Response := make(map[string]interface{})
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

	Response["Orders"] = SliceOrder
	middleware.SendResponse(response, Response)
}

package controller

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Order struct {
	id        string
	date      time.Time
	user      string
	delivered int
	paid      int
	method    string
	confirmed int
}

type orderProducts struct {
	order    string
	name     string
	image    string
	price    int
	quantity int
}

type Orders map[string]interface{}
type OrdersProducts map[string]interface{}

func CreateOrder(db *sql.DB, order map[string]interface{}) map[string]interface{} {
	token, deleverd, paid, method := fmt.Sprintf("%v", order["token"]), order["deleverd"], order["paid"], fmt.Sprintf("%v", order["method"])
	id := uuid.New()
	var OrderRes = make(map[string]interface{})
	claims, err := verifyToken(token)
	if err != nil {
		OrderRes["Error"] = true
		return OrderRes
	}
	var tm time.Time
	switch iat := claims["exp"].(type) {
	case float64:
		tm = time.Unix(int64(iat), 0)
	case json.Number:
		v, _ := iat.Int64()
		tm = time.Unix(v, 0)
	}

	if tm == time.Now() {
		OrderRes["Error"] = true
		return OrderRes
	}

	FindEmail := db.QueryRow("SELECT id FROM Users WHERE email = ?", claims["email"])

	var User user

	Err := FindEmail.Scan(&User.id)
	if Err == nil {
		_, err := db.Exec("INSERT INTO Orders (id, date, user, delivered, paid, method) VALUES (?, ?, ?, ?, ?, ?)", id, time.Now(), User.id, deleverd, paid, method)
		if err != nil {
			panic(err.Error())
		}

	}
	OrderRes["Error"] = false
	OrderRes["OrderId"] = id
	OrderRes["ID"] = User.id

	return OrderRes
}

func GetOrders(db *sql.DB, token string) map[string]interface{} {
	var Orders []Orders
	var OrderRes = make(map[string]interface{})
	var orderProduct []OrdersProducts

	claims, err := verifyToken(token)
	if err != nil {
		OrderRes["Error"] = true
		return OrderRes
	}
	var tm time.Time
	switch iat := claims["exp"].(type) {
	case float64:
		tm = time.Unix(int64(iat), 0)
	case json.Number:
		v, _ := iat.Int64()
		tm = time.Unix(v, 0)
	}

	if tm == time.Now() {
		OrderRes["Error"] = true
		return OrderRes
	}

	FindEmail := db.QueryRow("SELECT id FROM Users WHERE email = ?", claims["email"])

	var User user

	Err := FindEmail.Scan(&User.id)

	if Err == nil {

		Select, err := db.Query("SELECT * FROM `Orders` WHERE user = ?", User.id)

		if err != nil {
			panic(err.Error())
		}

		defer Select.Close()

		for Select.Next() {
			var Order Order

			if err := Select.Scan(&Order.id, &Order.date, &Order.user, &Order.delivered, &Order.paid, &Order.method, &Order.confirmed); err != nil {
				panic(err.Error())
			}

			TheOrder := map[string]interface{}{
				"id":        Order.id,
				"date":      Order.date,
				"user":      Order.user,
				"delivered": Order.delivered,
				"paid":      Order.paid,
				"method":    Order.method,
				"confirmed": Order.confirmed,
			}

			Orders = append(Orders, TheOrder)
		}

		ordersIds := ""

		for i, order := range Orders {
			if i == 0 {
				ordersIds += fmt.Sprintf("'%v'", order["id"])
			}
			if i > 0 {
				ordersIds += fmt.Sprintf(",'%v'", order["id"])
			}
		}

		sql := fmt.Sprintf("SELECT OrderProducts.quantity, OrderProducts.`order`, Products.`name`, Products.image, Products.price  FROM OrderProducts INNER JOIN Products ON OrderProducts.product=Products.`id`  WHERE OrderProducts.`order` IN (%s)", ordersIds)
		Products, err := db.Prepare(sql)

		if err != nil {
			fmt.Println(err)
		}

		rows, err := Products.Query()

		if err != nil {
			fmt.Println(err)
			OrderRes["Error"] = true
			return OrderRes
		}

		defer rows.Close()

		for rows.Next() {

			var product orderProducts

			if err := rows.Scan(&product.quantity, &product.order, &product.name, &product.image, &product.price); err != nil {
				panic(err.Error())
			}

			TheOrderProduct := map[string]interface{}{
				"order":    product.order,
				"name":     product.name,
				"image":    product.image,
				"price":    product.price,
				"quantity": product.quantity,
			}

			orderProduct = append(orderProduct, TheOrderProduct)
		}
	}

	OrderRes["Order"] = Orders
	OrderRes["Products"] = orderProduct
	return OrderRes
}

func CancelOrder(db *sql.DB, token, order string, confirmed int) map[string]bool {
	var OrderRes = make(map[string]bool)

	claims, err := verifyToken(token)
	if err != nil {
		OrderRes["Error"] = true
		return OrderRes
	}
	var tm time.Time
	switch iat := claims["exp"].(type) {
	case float64:
		tm = time.Unix(int64(iat), 0)
	case json.Number:
		v, _ := iat.Int64()
		tm = time.Unix(v, 0)
	}

	if tm == time.Now() {
		OrderRes["Error"] = true
		return OrderRes
	}

	userEmail := claims["email"]

	FindEmail := db.QueryRow("SELECT id FROM Users WHERE email = ?", userEmail)

	var User user

	Err := FindEmail.Scan(&User.id)

	if Err == nil {
		if confirmed == 0 {
			_, err := db.Exec("DELETE FROM Orders WHERE (id, user, confirmed)=(?, ?, ?)", order, User.id, 0)

			if err != nil {
				panic(err.Error())
			}

			return map[string]bool{
				"Error": false,
			}
		}

		return map[string]bool{
			"Error": true,
		}
	}

	return map[string]bool{
		"Error": true,
	}
}

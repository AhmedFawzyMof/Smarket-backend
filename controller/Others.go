package controller

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

func ForYou(db *sql.DB, token string, resChan chan any, wg *sync.WaitGroup) {
	var Orders []Orders
	var products []Products

	claims, err := verifyToken(token)
	if err != nil {
		UserRes := map[string]interface{}{
			"Error": true,
		}
		resChan <- UserRes
		wg.Done()
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
		UserRes := map[string]interface{}{
			"Error": true,
		}

		resChan <- UserRes
		wg.Done()
	}

	FindEmail := db.QueryRow("SELECT id FROM Users WHERE email = ?", claims["email"])

	var User user

	Err := FindEmail.Scan(&User.id)

	if Err == nil {
		Products, err := db.Query("SELECT * FROM Products WHERE offer > 0")

		if err != nil {
			panic(err.Error())
		}

		defer Products.Close()

		for Products.Next() {
			var Product product

			if err := Products.Scan(&Product.id, &Product.name, &Product.description, &Product.price, &Product.company, &Product.subcategories, &Product.category, &Product.image, &Product.unit, &Product.available, &Product.offer, &Product.inStock, &Product.pricePerUint, &Product.unitNumber); err != nil {
				productRes := make(map[string]interface{})
				productRes["Error"] = true

				products = append(products, productRes)

				resChan <- products
				wg.Done()
			}

			TheProduct := map[string]interface{}{
				"id":            Product.id,
				"name":          Product.name,
				"description":   Product.description,
				"price":         Product.price,
				"company":       Product.company,
				"subcategories": Product.subcategories,
				"category":      Product.category,
				"image":         Product.image,
				"unit":          Product.unit,
				"available":     Product.available,
				"offer":         Product.offer,
				"pricePerUint":  Product.pricePerUint,
				"unitNumber":    Product.unitNumber,
			}
			products = append(products, TheProduct)
		}

		resChan <- products
		wg.Done()
		return
	}

	Select, err := db.Query("SELECT id FROM `Orders` WHERE user = ?", User.id)
	if err != nil {
		panic(err.Error())
	}

	defer Select.Close()

	for Select.Next() {
		var Order Order

		if err := Select.Scan(&Order.id); err != nil {
			panic(err.Error())
		}

		TheOrder := map[string]interface{}{
			"id": Order.id,
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
	if ordersIds != "" {
		productCategory := getProducts(db, ordersIds)

		categroies := ""

		for i, cate := range productCategory {
			if i == 0 {
				categroies += fmt.Sprintf("('%s'", cate["category"])
			}
			if i > 0 {
				categroies += fmt.Sprintf(",'%s'", cate["category"])
			}
			if i == len(productCategory)-1 {
				categroies += ")"
			}
		}

		stmt := fmt.Sprintf("SELECT * FROM `Products` WHERE category IN %s", categroies)

		Select, err := db.Query(stmt)
		if err != nil {
			panic(err.Error())
		}

		defer Select.Close()

		for Select.Next() {
			var Product product

			if err := Select.Scan(&Product.id, &Product.name, &Product.description, &Product.price, &Product.company, &Product.subcategories, &Product.category, &Product.image, &Product.unit, &Product.available, &Product.offer, &Product.inStock, &Product.pricePerUint, &Product.unitNumber); err != nil {
				productRes := make(map[string]interface{})
				productRes["Error"] = true

				products = append(products, productRes)

				resChan <- products
				wg.Done()
			}

			TheProduct := map[string]interface{}{
				"id":            Product.id,
				"name":          Product.name,
				"description":   Product.description,
				"price":         Product.price,
				"company":       Product.company,
				"subcategories": Product.subcategories,
				"category":      Product.category,
				"image":         Product.image,
				"unit":          Product.unit,
				"available":     Product.available,
				"offer":         Product.offer,
				"pricePerUint":  Product.pricePerUint,
				"unitNumber":    Product.unitNumber,
			}
			products = append(products, TheProduct)

		}
		resChan <- products
		wg.Done()
	}
}

func getProducts(db *sql.DB, ordersIds string) []Products {
	var products []Products

	sql := fmt.Sprintf("SELECT Products.`category`, COUNT(Products.`id`) AS `Count` FROM OrderProducts INNER JOIN Products ON OrderProducts.product=Products.`id` WHERE OrderProducts.`order` IN (%s) GROUP BY Products.`category` ORDER BY Count DESC", ordersIds)
	Product, err := db.Prepare(sql)

	if err != nil {
		panic(err.Error())
	}

	rows, err := Product.Query()

	if err != nil {
		panic(err.Error())
	}

	defer rows.Close()

	for rows.Next() {

		var product product
		var count int

		if err := rows.Scan(&product.category, &count); err != nil {
			panic(err.Error())
		}

		TheProduct := map[string]interface{}{
			"category": product.category,
			"count":    count,
		}

		products = append(products, TheProduct)
	}
	return products
}

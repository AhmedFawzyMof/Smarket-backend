package controller

import (
	"database/sql"
	"sync"
)

type product struct {
	id            int
	name          string
	description   string
	price         int
	company       string
	subcategories sql.NullString
	category      string
	image         string
	unit          string
	available     int
	offer         int
	inStock       int
}

type Products map[string]interface{}

func ProductGetAll(db *sql.DB, responseChan chan []Products, wg *sync.WaitGroup) {
	var Products []Products

	Select, err := db.Query("SELECT * FROM `Products`")

	if err != nil {
		panic(err.Error())
	}

	defer Select.Close()

	for Select.Next() {
		var Product product

		if err := Select.Scan(&Product.id, &Product.name, &Product.description, &Product.price, &Product.company, &Product.subcategories, &Product.category, &Product.image, &Product.unit, &Product.available, &Product.offer, &Product.inStock); err != nil {
			panic(err.Error())
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
			"inStock":         Product.inStock,
		}

		Products = append(Products, TheProduct)

	}

	responseChan <- Products
	wg.Done()
}

func ProductGetId(db *sql.DB, responseChan chan Products, wg *sync.WaitGroup, ID int) {
	var Products Products

	Select := db.QueryRow("SELECT * FROM Products WHERE id=?", ID)

	var Product product

	err := Select.Scan(&Product.id, &Product.name, &Product.description, &Product.price, &Product.company, &Product.subcategories, &Product.category, &Product.image, &Product.unit, &Product.available, &Product.offer)

	if err == nil {
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
		}

		Products = TheProduct
	}

	responseChan <- Products
	wg.Done()
}

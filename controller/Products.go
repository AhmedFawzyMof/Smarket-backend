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
	subcategories string
	category      string
	image         []byte
	unit          string
	available     int
	offer         int
	inStock       int
	pricePerUint  float32
	unitNumber    int
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

		if err := Select.Scan(&Product.id, &Product.name, &Product.description, &Product.price, &Product.company, &Product.subcategories, &Product.category, &Product.image, &Product.unit, &Product.available, &Product.offer, &Product.inStock, &Product.pricePerUint, &Product.unitNumber); err != nil {
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
			"inStock":       Product.inStock,
			"pricePerUint":  Product.pricePerUint,
			"unitNumber":    Product.unitNumber,
		}

		Products = append(Products, TheProduct)

	}

	responseChan <- Products
	wg.Done()
}

func ProductGetId(db *sql.DB, ID int) Products {

	Select := db.QueryRow("SELECT * FROM Products WHERE id=?", ID)

	var Product product

	if err := Select.Scan(&Product.id, &Product.name, &Product.description, &Product.price, &Product.company, &Product.subcategories, &Product.category, &Product.image, &Product.unit, &Product.available, &Product.offer, &Product.inStock, &Product.pricePerUint, &Product.unitNumber); err != nil {
		return map[string]interface{}{}
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
		"inStock":       Product.inStock,
		"pricePerUint":  Product.pricePerUint,
		"unitNumber":    Product.unitNumber,
	}
	return TheProduct

}

func ProductOffers(db *sql.DB, resChan chan []Products, wg *sync.WaitGroup) {
	var products []Products
	Select, err := db.Query("SELECT * FROM Products WHERE offer > 0")

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

func GetallProduct(db *sql.DB, productChan chan []Products, wg *sync.WaitGroup) {
	var Products []Products

	Select, err := db.Query("SELECT * FROM `Products`")

	if err != nil {
		panic(err.Error())
	}

	defer Select.Close()

	for Select.Next() {
		var Product product

		if err := Select.Scan(&Product.id, &Product.name, &Product.description, &Product.price, &Product.company, &Product.subcategories, &Product.category, &Product.image, &Product.unit, &Product.available, &Product.offer, &Product.inStock, &Product.pricePerUint, &Product.unitNumber); err != nil {
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
			"inStock":       Product.inStock,
			"pricePerUint":  Product.pricePerUint,
			"unitNumber":    Product.unitNumber,
		}

		Products = append(Products, TheProduct)

	}

	productChan <- Products
	wg.Done()
}

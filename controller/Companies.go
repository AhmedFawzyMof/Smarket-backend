package controller

import (
	"database/sql"
	"sync"
)

type company struct {
	name string
}

type Companies map[string]interface{}

func CompanyGetAllProducts(db *sql.DB, responseChan chan []Companies, wg *sync.WaitGroup, slug string) {
	var Products []Companies

	Select, err := db.Query("SELECT * FROM `Products` WHERE company=?", slug)

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

func GetAllComponies(db *sql.DB, responseChan chan []Companies, wg *sync.WaitGroup) {
	var Companies []Companies

	Select, err := db.Query("SELECT * FROM `Companies`")

	if err != nil {
		panic(err.Error())
	}

	defer Select.Close()

	for Select.Next() {
		var Company company

		if err := Select.Scan(&Company.name); err != nil {
			panic(err.Error())
		}
		TheCompany := map[string]interface{}{
			"name": Company.name,
		}

		Companies = append(Companies, TheCompany)

	}

	responseChan <- Companies
	wg.Done()
}

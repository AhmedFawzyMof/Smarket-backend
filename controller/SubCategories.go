package controller

import (
	"database/sql"
	"sync"
)

type subcategory struct {
	name     string
	category string
	image    []byte
}

type SubCategories map[string]interface{}

func GetAllSubs(db *sql.DB, subCates chan []SubCategories, wg *sync.WaitGroup) {
	var SubCategories []SubCategories

	Select, err := db.Query("SELECT * FROM `SubCategories`")

	if err != nil {
		panic(err.Error())
	}

	defer Select.Close()

	for Select.Next() {
		var subcategory subcategory

		if err := Select.Scan(&subcategory.name, &subcategory.category, &subcategory.image); err != nil {
			panic(err.Error())
		}

		theCategory := map[string]interface{}{
			"name":     subcategory.name,
			"category": subcategory.category,
			"image":    subcategory.image,
		}

		SubCategories = append(SubCategories, theCategory)
	}

	subCates <- SubCategories
	wg.Done()
}

func CategoriesGetAllSubs(db *sql.DB, responseChan chan []SubCategories, wg *sync.WaitGroup, category string) {
	var Products []SubCategories

	Select, err := db.Query("SELECT * FROM `SubCategories` WHERE category=?", category)

	if err != nil {
		panic(err.Error())
	}

	defer Select.Close()

	for Select.Next() {
		var Product subcategory

		if err := Select.Scan(&Product.name, &Product.category, &Product.image); err != nil {
			panic(err.Error())
		}

		TheProduct := map[string]interface{}{
			"name":     Product.name,
			"category": Product.category,
			"image":    Product.image,
		}

		Products = append(Products, TheProduct)

	}

	responseChan <- Products
	wg.Done()
}

func SubCategoriesProduct(db *sql.DB, responseChan chan []Products, wg *sync.WaitGroup, subcategory string) {
	var Products []Products
	Select, err := db.Query("SELECT * FROM `Products` WHERE `subcategories`=?", subcategory)

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

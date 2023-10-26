package controller

import (
	"database/sql"
	"sync"
)



type category struct {
	name  string
	image string
}

type Categories map[string]interface{}

func CategoryGetAll(db *sql.DB, responseChan chan []Categories, wg *sync.WaitGroup) {
	var Categories []Categories

	Select, err := db.Query("SELECT * FROM `Categories`")

	if err != nil {
		panic(err.Error())
	}

	defer Select.Close()

	for Select.Next() {
		var category category

		if err := Select.Scan(&category.name, &category.image); err != nil {
			panic(err.Error())
		}

		theCategory := map[string]interface{}{
			"name":  category.name,
			"image": category.image,
		}

		Categories = append(Categories, theCategory)
	}

	responseChan <- Categories
	wg.Done()
}

func CategoriesGetAllProducts(db *sql.DB, responseChan chan []Categories, wg *sync.WaitGroup, slug string) {
	var Products []Categories

	Select, err := db.Query("SELECT * FROM `Products` WHERE category=?", slug)

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
			"inStock":		 Product.inStock,
		}

		Products = append(Products, TheProduct)

	}

	responseChan <- Products
	wg.Done()
}

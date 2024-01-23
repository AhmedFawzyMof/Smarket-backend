package models

import (
	"alwadi_markets/middleware"
	"database/sql"
	"encoding/json"
	"fmt"
	"sync"
)

type Product struct {
	Id            int
	Name          string
	Description   string
	Price         int
	Offer         int
	Company       string
	Subcategories string
	Category      string
	Available     int
	Image         string
	Types         []ProductType
}

func (p Product) Add(db *sql.DB, response chan []byte, wg *sync.WaitGroup) {
	defer wg.Done()
	Response := make(map[string]interface{})

	_, err := db.Exec("INSERT INTO `Products`(`name`, `description`, `company`, `subcategories`, `category`, `image`, `available`, `price`, `offer`) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)", p.Name, p.Description, p.Company, p.Subcategories, p.Category, p.Image, p.Available, p.Price, p.Offer)

	if err != nil {
		Response["Error"] = true

		middleware.SendResponse(response, Response)
		return
	}
	Response["Error"] = false

	middleware.SendResponse(response, Response)
}

func (p Product) Update(db *sql.DB, response chan []byte, wg *sync.WaitGroup) {
	defer wg.Done()
	Response := make(map[string]interface{})

	_, err := db.Exec("UPDATE `Products` SET `name`=?, `description`=?, `company`=?, `subcategories`=?, `category`=?, `available`=?, `price`=?, `offer`=? WHERE id = ?", p.Name, p.Description, p.Company, p.Subcategories, p.Category, p.Available, p.Price, p.Offer, p.Id)

	if err != nil {
		Response["Error"] = true
		middleware.SendResponse(response, Response)
		return
	}
	Response["Error"] = false

	middleware.SendResponse(response, Response)
}

func (p Product) Get(db *sql.DB, response chan []byte, wg *sync.WaitGroup) {
	var Products []Product
	products, err := db.Query("SELECT * FROM `Products`")

	if err != nil {
		panic(err.Error())
	}

	defer products.Close()

	for products.Next() {
		var product Product

		if err := products.Scan(&product.Id, &product.Name, &product.Description, &product.Company, &product.Subcategories, &product.Category, &product.Image, &product.Available, &product.Price, &product.Offer); err != nil {
			panic(err.Error())
		}

		Products = append(Products, product)
	}

	res, err := json.Marshal(Products)

	if err != nil {
		fmt.Println(err)
	}

	response <- res
	wg.Done()
}

func (p Product) GetByCategory(db *sql.DB, response chan []byte, wg *sync.WaitGroup) {
	var Products []Product
	products, err := db.Query("SELECT * FROM `Products` WHERE category =?", p.Category)

	if err != nil {
		panic(err.Error())
	}

	defer products.Close()

	for products.Next() {
		var product Product

		if err := products.Scan(&product.Id, &product.Name, &product.Description, &product.Company, &product.Subcategories, &product.Category, &product.Image, &product.Available, &product.Price, &product.Offer); err != nil {
			panic(err.Error())
		}

		Products = append(Products, product)
	}

	res, err := json.Marshal(Products)

	if err != nil {
		fmt.Println(err)
	}

	response <- res
	wg.Done()
}

func (p Product) GetByCompany(db *sql.DB, response chan []byte, wg *sync.WaitGroup) {
	var Products []Product
	products, err := db.Query("SELECT * FROM `Products` WHERE company=?", p.Company)

	if err != nil {
		panic(err.Error())
	}

	defer products.Close()

	for products.Next() {
		var product Product

		if err := products.Scan(&product.Id, &product.Name, &product.Description, &product.Company, &product.Subcategories, &product.Category, &product.Image, &product.Available, &product.Price, &product.Offer); err != nil {
			panic(err.Error())
		}

		Products = append(Products, product)
	}

	res, err := json.Marshal(Products)

	if err != nil {
		fmt.Println(err)
	}

	response <- res
	wg.Done()
}

func (p Product) GetBySubCategory(db *sql.DB, response chan []byte, wg *sync.WaitGroup) {
	var Products []Product
	products, err := db.Query("SELECT * FROM `Products` WHERE subCategories = ?", p.Subcategories)

	if err != nil {
		panic(err.Error())
	}

	defer products.Close()

	for products.Next() {
		var product Product

		if err := products.Scan(&product.Id, &product.Name, &product.Description, &product.Company, &product.Subcategories, &product.Category, &product.Image, &product.Available, &product.Price, &product.Offer); err != nil {
			panic(err.Error())
		}

		Products = append(Products, product)
	}

	res, err := json.Marshal(Products)

	if err != nil {
		fmt.Println(err)
	}

	response <- res
	wg.Done()
}

func (p Product) GetById(db *sql.DB, response chan []byte, wg *sync.WaitGroup) {
	product := db.QueryRow("SELECT * FROM `Products` WHERE id = ?", p.Id)

	var Product Product
	var types []ProductType

	if err := product.Scan(&Product.Id, &Product.Name, &Product.Description, &Product.Company, &Product.Subcategories, &Product.Category, &Product.Image, &Product.Available, &Product.Price, &Product.Offer); err != nil {
		panic(err.Error())
	}

	products, err := db.Query("SELECT * FROM `ProductType` WHERE product = ?", p.Id)

	if err != nil {
		panic(err.Error())
	}

	defer products.Close()

	for products.Next() {
		var product ProductType

		if err := products.Scan(&product.Id, &product.Product, &product.Portion, &product.Uint, &product.Price, &product.Offer); err != nil {
			panic(err.Error())
		}

		types = append(types, product)
	}

	Product.Types = types

	res, err := json.Marshal(Product)

	if err != nil {
		fmt.Println(err)
	}

	response <- res
	wg.Done()
}

func (p Product) GetByOffers(db *sql.DB, response chan []byte, wg *sync.WaitGroup) {
	var Products []Product
	products, err := db.Query("SELECT * FROM `Products` WHERE offer > 0")

	if err != nil {
		panic(err.Error())
	}

	defer products.Close()

	for products.Next() {
		var product Product

		if err := products.Scan(&product.Id, &product.Name, &product.Description, &product.Company, &product.Subcategories, &product.Category, &product.Image, &product.Available, &product.Price, &product.Offer); err != nil {
			panic(err.Error())
		}

		Products = append(Products, product)
	}

	res, err := json.Marshal(Products)

	if err != nil {
		fmt.Println(err)
	}

	response <- res
	wg.Done()
}

func (p Product) Delete(db *sql.DB, response chan []byte, wg *sync.WaitGroup) {
	defer wg.Done()
	Response := make(map[string]interface{})
	_, err := db.Exec("DELETE FROM `Products` WHERE id = ?", p.Id)

	if err != nil {
		Response["Error"] = true

		middleware.SendResponse(response, Response)
		return
	}

	Response["Error"] = true
	middleware.SendResponse(response, Response)
}

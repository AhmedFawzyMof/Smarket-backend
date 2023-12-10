package admin

import (
	DB "alwadi_markets/db"
	"alwadi_markets/tables"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"sync"
)

func GetProducts(res http.ResponseWriter, req *http.Request, params map[string]string) {
	db := DB.Connect()
	res.WriteHeader(http.StatusOK)

	defer db.Close()

	Products := make(chan []byte, 1)

	wg := &sync.WaitGroup{}

	wg.Add(1)

	go tables.Product.Get(tables.Product{}, db, Products, wg)

	wg.Wait()

	close(Products)

	var products []tables.Product

	err := json.Unmarshal(<-Products, &products)

	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}

	Response := make(map[string]interface{}, 1)

	Response["Products"] = products

	if err := json.NewEncoder(res).Encode(Response); err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}
}

func AddProduct(res http.ResponseWriter, req *http.Request, params map[string]string) {
	db := DB.Connect()
	res.WriteHeader(http.StatusOK)

	defer db.Close()

	body, err := io.ReadAll(req.Body)

	if err != nil {
		panic(err.Error())
	}
	var productmap map[string]interface{}

	var Product tables.Product

	Error := json.Unmarshal(body, &productmap)

	if Error != nil {
		http.Error(res, Error.Error(), http.StatusInternalServerError)
	}

	// body data
	price, _ := strconv.Atoi(fmt.Sprintf("%s", productmap["price"]))
	available, _ := strconv.Atoi(fmt.Sprintf("%s", productmap["available"]))
	Product.Name = fmt.Sprintf("%s", productmap["name"])
	Product.Description = fmt.Sprintf("%s", productmap["description"])
	Product.Price = price
	Product.Image = fmt.Sprintf("%s", productmap["image"])
	Product.Subcategories = fmt.Sprintf("%s", productmap["subcategories"])
	Product.Category = fmt.Sprintf("%s", productmap["category"])
	Product.Company = fmt.Sprintf("%s", productmap["company"])
	Product.Available = available

	Products := make(chan []byte, 1)

	wg := &sync.WaitGroup{}

	wg.Add(1)
	go tables.Product.Add(Product, db, Products, wg)
	wg.Wait()

	close(Products)
	var products map[string]interface{}

	errors := json.Unmarshal(<-Products, &products)

	if errors != nil {
		http.Error(res, errors.Error(), http.StatusInternalServerError)
	}

	if err := json.NewEncoder(res).Encode(products); err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}
}

func UpdateProduct(res http.ResponseWriter, req *http.Request, params map[string]string) {
	if req.Method == "PUT" {
		db := DB.Connect()
		res.WriteHeader(http.StatusOK)

		defer db.Close()

		body, err := io.ReadAll(req.Body)

		if err != nil {
			panic(err.Error())
		}
		var productmap map[string]interface{}
		var Product tables.Product

		Error := json.Unmarshal(body, &productmap)
		if Error != nil {
			http.Error(res, Error.Error(), http.StatusInternalServerError)
		}

		// body data
		price, _ := strconv.Atoi(fmt.Sprintf("%s", productmap["price"]))
		available, _ := strconv.Atoi(fmt.Sprintf("%s", productmap["available"]))
		id, _ := strconv.Atoi(fmt.Sprintf("%s", productmap["id"]))
		Product.Name = fmt.Sprintf("%s", productmap["name"])
		Product.Description = fmt.Sprintf("%s", productmap["description"])
		Product.Price = price
		Product.Subcategories = fmt.Sprintf("%s", productmap["subcategories"])
		Product.Category = fmt.Sprintf("%s", productmap["category"])
		Product.Company = fmt.Sprintf("%s", productmap["company"])
		Product.Available = available
		Product.Id = id

		fmt.Println(Product.Available)
		Products := make(chan []byte, 1)

		wg := &sync.WaitGroup{}

		wg.Add(1)

		go tables.Product.Update(Product, db, Products, wg)

		wg.Wait()

		close(Products)

		var products map[string]interface{}

		errors := json.Unmarshal(<-Products, &products)

		if errors != nil {
			http.Error(res, errors.Error(), http.StatusInternalServerError)
		}

		if err := json.NewEncoder(res).Encode(products); err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
		}
	}
}

func DeleteProduct(res http.ResponseWriter, req *http.Request, params map[string]string) {
	if req.Method == "DELETE" {
		db := DB.Connect()
		res.WriteHeader(http.StatusOK)

		defer db.Close()
		body, err := io.ReadAll(req.Body)

		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
		}

		var Product tables.Product

		Error := json.Unmarshal(body, &Product)
		if Error != nil {
			http.Error(res, Error.Error(), http.StatusInternalServerError)
		}

		Products := make(chan []byte, 1)

		wg := &sync.WaitGroup{}

		wg.Add(1)

		go tables.Product.Delete(Product, db, Products, wg)

		wg.Wait()

		close(Products)

		var products map[string]interface{}

		errors := json.Unmarshal(<-Products, &products)

		if errors != nil {
			http.Error(res, errors.Error(), http.StatusInternalServerError)
		}

		if err := json.NewEncoder(res).Encode(products); err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
		}
	}
}

package admin

import (
	DB "alwadi_markets/db"
	"alwadi_markets/middleware"
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

	body, err := io.ReadAll(req.Body)

	if err != nil {
		middleware.SendError(err, res)
	}
	var productmap map[string]interface{}

	if err := json.Unmarshal(body, &productmap); err != nil {
		middleware.SendError(err, res)
	}

	var token string = fmt.Sprintf("%s", productmap["auth-token"])

	admin := middleware.CheckIsAdmin(token, db)

	if !admin {
		err := fmt.Errorf("user is not admin")
		middleware.SendError(err, res)
	}

	Products := make(chan []byte, 1)

	wg := &sync.WaitGroup{}

	wg.Add(1)

	go tables.Product.Get(tables.Product{}, db, Products, wg)

	wg.Wait()

	close(Products)

	var products []tables.Product

	if err := json.Unmarshal(<-Products, &products); err != nil {
		middleware.SendError(err, res)
	}

	Response := make(map[string]interface{}, 1)

	Response["Products"] = products

	if err := json.NewEncoder(res).Encode(Response); err != nil {
		middleware.SendError(err, res)
	}
}

func AddProduct(res http.ResponseWriter, req *http.Request, params map[string]string) {
	db := DB.Connect()
	res.WriteHeader(http.StatusOK)

	defer db.Close()

	body, err := io.ReadAll(req.Body)

	if err != nil {
		middleware.SendError(err, res)
	}
	var productmap map[string]interface{}

	var Product tables.Product

	if err := json.Unmarshal(body, &productmap); err != nil {
		middleware.SendError(err, res)
	}

	var token string = fmt.Sprintf("%s", productmap["auth-token"])

	admin := middleware.CheckIsAdmin(token, db)

	if !admin {
		err := fmt.Errorf("user is not admin")
		middleware.SendError(err, res)
	}

	// body data

	Product.Name = fmt.Sprintf("%s", productmap["name"])
	Product.Description = fmt.Sprintf("%s", productmap["description"])
	Product.Price = middleware.ConvertToInt(productmap["price"], res)
	Product.Image = fmt.Sprintf("%s", productmap["image"])
	Product.Subcategories = fmt.Sprintf("%s", productmap["subcategories"])
	Product.Category = fmt.Sprintf("%s", productmap["category"])
	Product.Company = fmt.Sprintf("%s", productmap["company"])
	Product.Available = middleware.ConvertToInt(productmap["available"], res)

	Products := make(chan []byte, 1)

	wg := &sync.WaitGroup{}

	wg.Add(1)
	go tables.Product.Add(Product, db, Products, wg)
	wg.Wait()

	close(Products)
	var products map[string]interface{}

	if err := json.Unmarshal(<-Products, &products); err != nil {
		middleware.SendError(err, res)
	}

	if err := json.NewEncoder(res).Encode(products); err != nil {
		middleware.SendError(err, res)
	}
}

func UpdateProduct(res http.ResponseWriter, req *http.Request, params map[string]string) {
	if req.Method == "PUT" {
		db := DB.Connect()
		res.WriteHeader(http.StatusOK)

		defer db.Close()

		body, err := io.ReadAll(req.Body)

		if err != nil {
			middleware.SendError(err, res)
		}
		var productmap map[string]interface{}
		var Product tables.Product

		if err := json.Unmarshal(body, &productmap); err != nil {
			middleware.SendError(err, res)
		}

		var token string = fmt.Sprintf("%s", productmap["auth-token"])

		admin := middleware.CheckIsAdmin(token, db)

		if !admin {
			err := fmt.Errorf("user is not admin")
			middleware.SendError(err, res)
		}

		// body data
		Product.Name = fmt.Sprintf("%s", productmap["name"])
		Product.Description = fmt.Sprintf("%s", productmap["description"])
		Product.Price = middleware.ConvertToInt(productmap["price"], res)
		Product.Subcategories = fmt.Sprintf("%s", productmap["subcategories"])
		Product.Category = fmt.Sprintf("%s", productmap["category"])
		Product.Company = fmt.Sprintf("%s", productmap["company"])
		Product.Available = middleware.ConvertToInt(productmap["available"], res)
		Product.Id = middleware.ConvertToInt(productmap["id"], res)

		Products := make(chan []byte, 1)

		wg := &sync.WaitGroup{}

		wg.Add(1)

		go tables.Product.Update(Product, db, Products, wg)

		wg.Wait()

		close(Products)

		var products map[string]interface{}

		if err := json.Unmarshal(<-Products, &products); err != nil {
			middleware.SendError(err, res)
		}

		if err := json.NewEncoder(res).Encode(products); err != nil {
			middleware.SendError(err, res)
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
			middleware.SendError(err, res)
		}

		var productmap map[string]interface{}
		var Product tables.Product

		if err := json.Unmarshal(body, &Product); err != nil {
			middleware.SendError(err, res)
		}

		var token string = fmt.Sprintf("%s", productmap["auth-token"])

		admin := middleware.CheckIsAdmin(token, db)

		if !admin {
			err := fmt.Errorf("user is not admin")
			middleware.SendError(err, res)
		}

		id, _ := strconv.Atoi(fmt.Sprintf("%s", productmap["id"]))

		Product.Id = id

		Products := make(chan []byte, 1)

		wg := &sync.WaitGroup{}

		wg.Add(1)

		go tables.Product.Delete(Product, db, Products, wg)

		wg.Wait()

		close(Products)

		var products map[string]interface{}

		if err := json.Unmarshal(<-Products, &products); err != nil {
			middleware.SendError(err, res)
		}

		if err := json.NewEncoder(res).Encode(products); err != nil {
			middleware.SendError(err, res)
		}
	}
}

package admin

import (
	"alwadi/controller"
	DB "alwadi/db"
	"encoding/json"
	"net/http"
	"sync"
)

func GetProducts(res http.ResponseWriter, req *http.Request) {
	db := DB.Connect()

	defer db.Close()
	res.Header().Set("Access-Control-Allow-Origin", "*")
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)

	products := make(chan []controller.Products, 1)
	wg := &sync.WaitGroup{}

	wg.Add(1)

	go controller.GetallProduct(db, products, wg)

	wg.Wait()

	close(products)

	var Products = map[string]interface{}{
		"Products": <-products,
	}

	json.NewEncoder(res).Encode(Products)
}

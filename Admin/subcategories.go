package admin

import (
	"alwadi/controller"
	DB "alwadi/db"
	"encoding/json"
	"net/http"
	"sync"
)

func GetSubCategories(res http.ResponseWriter, req *http.Request) {
	db := DB.Connect()

	defer db.Close()
	res.Header().Set("Access-Control-Allow-Origin", "*")
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)

	products := make(chan []controller.SubCategories, 1)
	wg := &sync.WaitGroup{}

	wg.Add(1)

	go controller.GetAllSubs(db, products, wg)

	wg.Wait()

	close(products)

	var Products = map[string]interface{}{
		"SubCategories": <-products,
	}

	json.NewEncoder(res).Encode(Products)
}

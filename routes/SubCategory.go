package routes

import (
	"Smarket/controller"
	DB "Smarket/db"
	"encoding/json"
	"net/http"
	"sync"
)

type SubCategorySlug struct {
	Slug string
}

func (c SubCategorySlug) GetBySlug(res http.ResponseWriter, req *http.Request) {
	db := DB.Connect()

	defer db.Close()
	res.Header().Set("Access-Control-Allow-Origin", "*")
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	Products := make(chan []controller.Products, 1)
	wg := &sync.WaitGroup{}

	wg.Add(1)
	go controller.SubCategoriesProduct(db, Products, wg, c.Slug)
	wg.Wait()
	close(Products)

	products := <-Products

	var data = map[string]interface{}{
		"Products": products,
	}

	json.NewEncoder(res).Encode(data)
}

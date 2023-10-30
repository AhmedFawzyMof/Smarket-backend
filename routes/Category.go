package routes

import (
	"Smarket/controller"
	DB "Smarket/db"
	"encoding/json"
	"net/http"
	"sync"
)

type CategorySlug struct {
	Slug string
}

func (c CategorySlug) GetBySlug(res http.ResponseWriter, req *http.Request) {
	db := DB.Connect()

	defer db.Close()
	res.Header().Set("Access-Control-Allow-Origin", "*")
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)

	Category := make(chan []controller.Categories, 1)
	wg := &sync.WaitGroup{}

	wg.Add(1)
	go controller.CategoriesGetAllProducts(db, Category, wg, c.Slug)

	wg.Wait()
	close(Category)

	products := <-Category

	var data = map[string]interface{}{
		"Products": products,
	}

	json.NewEncoder(res).Encode(data)
}

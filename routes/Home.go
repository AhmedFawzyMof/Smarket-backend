package routes

import (
	DB "alwadi_markets/db"
	"alwadi_markets/middleware"
	"alwadi_markets/tables"
	"encoding/json"
	"net/http"
	"sync"
)

func Home(res http.ResponseWriter, req *http.Request, params map[string]string) {
	res.WriteHeader(http.StatusOK)
	db := DB.Connect()
	defer db.Close()

	Categories := make(chan []byte, 1)
	Products := make(chan []byte, 1)
	Offers := make(chan []byte, 1)

	wg := &sync.WaitGroup{}

	wg.Add(3)

	go tables.Category.Get(tables.Category{}, db, Categories, wg)
	go tables.Product.Get(tables.Product{}, db, Products, wg)
	go tables.Offer.Get(tables.Offer{}, db, Offers, wg)

	wg.Wait()
	close(Products)
	close(Categories)
	close(Offers)

	var category []tables.Category

	err := json.Unmarshal(<-Categories, &category)

	if err != nil {
		middleware.SendError(err, res)
	}

	var product []tables.Product

	if err := json.Unmarshal(<-Products, &product); err != nil {
		middleware.SendError(err, res)
	}

	var offer []tables.Offer

	if err := json.Unmarshal(<-Offers, &offer); err != nil {
		middleware.SendError(err, res)
	}

	Respones := make(map[string]interface{}, 3)
	Respones["Products"] = product
	Respones["Categories"] = category
	Respones["Offers"] = offer

	if err := json.NewEncoder(res).Encode(Respones); err != nil {
		middleware.SendError(err, res)
	}
}

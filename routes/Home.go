package routes

import (
	"encoding/json"
	"net/http"
	"sync"

	DB "alwadi_markets/db"
	"alwadi_markets/middleware"
	"alwadi_markets/models"
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

	go models.Category.Get(models.Category{}, db, Categories, wg)
	go models.Product.Get(models.Product{}, db, Products, wg)
	go models.Offer.Get(models.Offer{}, db, Offers, wg)

	wg.Wait()
	close(Products)
	close(Categories)
	close(Offers)

	var category []models.Category

	err := json.Unmarshal(<-Categories, &category)

	if err != nil {
		middleware.SendError(err, res)
	}

	var product []models.Product

	if err := json.Unmarshal(<-Products, &product); err != nil {
		middleware.SendError(err, res)
	}

	var offer []models.Offer

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

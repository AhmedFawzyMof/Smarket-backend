package routes

import (
	DB "alwadi_markets/db"
	"alwadi_markets/tables"
	"encoding/json"
	"fmt"
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
		fmt.Println(err.Error())
	}

	var product []tables.Product

	Err := json.Unmarshal(<-Products, &product)

	if Err != nil {
		http.Error(res, Err.Error(), http.StatusInternalServerError)
	}

	var offer []tables.Offer

	ERR := json.Unmarshal(<-Offers, &offer)

	if ERR != nil {
		fmt.Println(ERR.Error())
	}

	Respones := make(map[string]interface{}, 3)
	Respones["Products"] = product
	Respones["Categories"] = category
	Respones["Offers"] = offer

	json.NewEncoder(res).Encode(Respones)
}

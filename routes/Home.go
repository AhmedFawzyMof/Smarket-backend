package routes

import (
	controller "Smarket/controller"
	DB "Smarket/db"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

func Home(res http.ResponseWriter, req *http.Request) {
	start := time.Now()
	db := DB.Connect()

	defer db.Close()
	res.Header().Set("Access-Control-Allow-Origin", "https://wild-pear-millipede.cyclic.app/")
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)

	Categories := make(chan []controller.Categories, 1)
	Products := make(chan []controller.Products, 1)
	Offers := make(chan []controller.Offers, 1)
	wg := &sync.WaitGroup{}

	wg.Add(3)

	go controller.ProductGetAll(db, Products, wg)
	go controller.CategoryGetAll(db, Categories, wg)
	go controller.OfferGetAll(db, Offers, wg)

	wg.Wait()
	close(Products)
	close(Categories)
	close(Offers)

	Product, Category, Offer := <-Products, <-Categories, <-Offers

	var data = map[string]interface{}{
		"Products":   Product,
		"Categories": Category,
		"Offers":     Offer,
	}

	ip := req.RemoteAddr

	fmt.Println(ip)

	json.NewEncoder(res).Encode(data)
	excuteTime := time.Since(start)
	fmt.Println(excuteTime)
}

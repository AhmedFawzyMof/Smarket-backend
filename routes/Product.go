package routes

import (
	"Smarket/controller"
	DB "Smarket/db"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

type ProductId struct {
	Id int
}

func (p ProductId) GetById(res http.ResponseWriter, req *http.Request) {
	db := DB.Connect()

	defer db.Close()
	res.Header().Set("Access-Control-Allow-Origin", "*")
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)

	Product := controller.ProductGetId(db, p.Id)

	fmt.Println(Product, p.Id)

	var data = map[string]interface{}{
		"product": Product,
	}

	json.NewEncoder(res).Encode(data)
}

func GetProductsOffers(res http.ResponseWriter, req *http.Request) {
	db := DB.Connect()

	defer db.Close()
	res.Header().Set("Access-Control-Allow-Origin", "*")
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)

	productChan := make(chan []controller.Products, 1)
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go controller.ProductOffers(db, productChan, wg)
	wg.Wait()
	close(productChan)

	products := <-productChan

	json.NewEncoder(res).Encode(products)
}

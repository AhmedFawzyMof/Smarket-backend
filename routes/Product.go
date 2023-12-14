package routes

import (
	DB "alwadi_markets/db"
	"alwadi_markets/middleware"
	"alwadi_markets/tables"
	"encoding/json"
	"net/http"
	"sync"
)

func ProductId(res http.ResponseWriter, req *http.Request, params map[string]string) {
	res.WriteHeader(http.StatusOK)
	db := DB.Connect()
	defer db.Close()

	var Product tables.Product

	Product.Id = middleware.ConvertToInt(params["id"], res)

	ProductChan := make(chan []byte, 1)

	wg := &sync.WaitGroup{}

	wg.Add(1)
	go tables.Product.GetById(Product, db, ProductChan, wg)
	wg.Wait()

	close(ProductChan)

	var Products tables.Product

	if err := json.Unmarshal(<-ProductChan, &Products); err != nil {
		middleware.SendError(err, res)
	}

	Response := make(map[string]interface{}, 1)
	Response["Product"] = Products

	if err := json.NewEncoder(res).Encode(Response); err != nil {
		middleware.SendError(err, res)
	}
}

func ProductInOffers(res http.ResponseWriter, req *http.Request, params map[string]string) {
	res.WriteHeader(http.StatusOK)
	db := DB.Connect()
	defer db.Close()

	var Product tables.Product

	ProductChan := make(chan []byte, 1)

	wg := &sync.WaitGroup{}

	wg.Add(1)
	go tables.Product.GetByOffers(Product, db, ProductChan, wg)
	wg.Wait()

	close(ProductChan)

	var Products []tables.Product

	if err := json.Unmarshal(<-ProductChan, &Products); err != nil {
		middleware.SendError(err, res)
	}

	Response := make(map[string]interface{}, 1)
	Response["Products"] = Products

	if err := json.NewEncoder(res).Encode(Response); err != nil {
		middleware.SendError(err, res)
	}
}

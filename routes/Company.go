package routes

import (
	DB "alwadi_markets/db"
	"alwadi_markets/middleware"
	"alwadi_markets/tables"
	"encoding/json"
	"net/http"
	"sync"
)

func Company(res http.ResponseWriter, req *http.Request, params map[string]string) {
	res.WriteHeader(http.StatusOK)

	db := DB.Connect()

	defer db.Close()

	var Product tables.Product

	Product.Company = params["name"]

	ProductChan := make(chan []byte, 1)

	wg := &sync.WaitGroup{}

	wg.Add(1)
	go tables.Product.GetByCompany(Product, db, ProductChan, wg)
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

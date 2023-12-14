package routes

import (
	DB "alwadi_markets/db"
	"alwadi_markets/middleware"
	"alwadi_markets/tables"
	"encoding/json"
	"net/http"
	"sync"
)

func Category(res http.ResponseWriter, req *http.Request, params map[string]string) {
	res.WriteHeader(http.StatusOK)

	db := DB.Connect()

	defer db.Close()

	var Product tables.Product
	var SubCategory tables.SubCategory

	Product.Category = params["name"]
	SubCategory.Category = params["name"]

	ProductChan := make(chan []byte, 1)
	SubCategoryChan := make(chan []byte, 1)

	wg := &sync.WaitGroup{}

	wg.Add(2)
	go tables.Product.GetByCategory(Product, db, ProductChan, wg)
	go tables.SubCategory.GetByCategory(SubCategory, db, SubCategoryChan, wg)
	wg.Wait()

	close(ProductChan)
	close(SubCategoryChan)

	var Products []tables.Product

	if err := json.Unmarshal(<-ProductChan, &Products); err != nil {
		middleware.SendError(err, res)
	}

	var SubCategories []tables.SubCategory

	if err := json.Unmarshal(<-SubCategoryChan, &SubCategories); err != nil {
		middleware.SendError(err, res)
	}

	Response := make(map[string]interface{}, 1)
	Response["Products"] = Products
	Response["SubCategories"] = SubCategories

	if err := json.NewEncoder(res).Encode(Response); err != nil {
		middleware.SendError(err, res)
	}
}

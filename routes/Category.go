package routes

import (
	DB "alwadi_markets/db"
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

	errors := json.Unmarshal(<-ProductChan, &Products)

	if errors != nil {
		http.Error(res, errors.Error(), http.StatusInternalServerError)
	}
	var SubCategories []tables.SubCategory

	erroR := json.Unmarshal(<-SubCategoryChan, &SubCategories)

	if erroR != nil {
		http.Error(res, erroR.Error(), http.StatusInternalServerError)
	}

	Response := make(map[string]interface{}, 1)
	Response["Products"] = Products
	Response["SubCategories"] = SubCategories

	if err := json.NewEncoder(res).Encode(Response); err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}
}

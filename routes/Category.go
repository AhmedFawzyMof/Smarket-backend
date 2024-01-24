package routes

import (
	"encoding/json"
	"net/http"
	"sync"

	DB "alwadimarkets/db"
	"alwadimarkets/middleware"
	"alwadimarkets/models"
)

func Category(res http.ResponseWriter, req *http.Request, params map[string]string) {
	res.WriteHeader(http.StatusOK)

	db := DB.Connect()

	defer db.Close()

	var Product models.Product
	var SubCategory models.SubCategory

	Product.Category = params["name"]
	SubCategory.Category = params["name"]

	ProductChan := make(chan []byte, 1)
	SubCategoryChan := make(chan []byte, 1)

	wg := &sync.WaitGroup{}

	wg.Add(2)
	go models.Product.GetByCategory(Product, db, ProductChan, wg)
	go models.SubCategory.GetByCategory(SubCategory, db, SubCategoryChan, wg)
	wg.Wait()

	close(ProductChan)
	close(SubCategoryChan)

	var Products []models.Product

	if err := json.Unmarshal(<-ProductChan, &Products); err != nil {
		middleware.SendError(err, res)
	}

	var SubCategories []models.SubCategory

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

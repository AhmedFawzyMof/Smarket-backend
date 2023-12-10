package admin

import (
	DB "alwadi_markets/db"
	"alwadi_markets/tables"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"sync"
)

func GetTypes(res http.ResponseWriter, req *http.Request, params map[string]string) {
	db := DB.Connect()
	res.WriteHeader(http.StatusOK)

	defer db.Close()

	Types := make(chan []byte, 1)

	wg := &sync.WaitGroup{}

	wg.Add(1)

	go tables.ProductType.Get(tables.ProductType{}, db, Types, wg)

	wg.Wait()

	close(Types)

	var Type []tables.ProductType

	err := json.Unmarshal(<-Types, &Type)

	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}

	Response := make(map[string]interface{}, 1)

	Response["Types"] = Type

	if err := json.NewEncoder(res).Encode(Response); err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}
}


func AddTypes(res http.ResponseWriter, req *http.Request, params map[string]string) {
	db := DB.Connect()
	res.WriteHeader(http.StatusOK)

	defer db.Close()

	body, err := io.ReadAll(req.Body)

	if err != nil {
		panic(err.Error())
	}
	var productmap map[string]interface{}

	var ProductType tables.ProductType

	Error := json.Unmarshal(body, &productmap)

	if Error != nil {
		http.Error(res, Error.Error(), http.StatusInternalServerError)
	}

	// body data
	price, _ := strconv.Atoi(fmt.Sprintf("%s", productmap["price"]))
	product, _ := strconv.Atoi(fmt.Sprintf("%s", productmap["product"]))
	offer, _ := strconv.Atoi(fmt.Sprintf("%s", productmap["offer"]))
	portion, _ := strconv.Atoi(fmt.Sprintf("%s", productmap["portion"]))
	ProductType.Product = product
	ProductType.Portion = portion
	ProductType.Price = price
	ProductType.Offer = offer
	ProductType.Uint = fmt.Sprintf("%s", productmap["uint"])
	

	Types := make(chan []byte, 1)

	wg := &sync.WaitGroup{}

	wg.Add(1)
	go tables.ProductType.Add(ProductType, db, Types, wg)
	wg.Wait()

	close(Types)
	var products map[string]interface{}

	errors := json.Unmarshal(<-Types, &products)

	if errors != nil {
		http.Error(res, errors.Error(), http.StatusInternalServerError)
	}

	if err := json.NewEncoder(res).Encode(products); err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}
}

func DeleteTypes(res http.ResponseWriter, req *http.Request, params map[string]string) {
	if req.Method == "DELETE"{
	db := DB.Connect()
	res.WriteHeader(http.StatusOK)

	defer db.Close()

	body, err := io.ReadAll(req.Body)

	if err != nil {
		panic(err.Error())
	}
	var productmap map[string]interface{}

	var ProductType tables.ProductType

	Error := json.Unmarshal(body, &productmap)

	if Error != nil {
		http.Error(res, Error.Error(), http.StatusInternalServerError)
	}

	// body data
	id, _ := strconv.Atoi(fmt.Sprintf("%s", productmap["id"]))
	ProductType.Id = id
	

	Types := make(chan []byte, 1)

	wg := &sync.WaitGroup{}

	wg.Add(1)
	go tables.ProductType.Delete(ProductType, db, Types, wg)
	wg.Wait()

	close(Types)
	var products map[string]interface{}

	errors := json.Unmarshal(<-Types, &products)

	if errors != nil {
		http.Error(res, errors.Error(), http.StatusInternalServerError)
	}

	if err := json.NewEncoder(res).Encode(products); err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}
	}
}

package routes

import (
	"Smarket/controller"
	DB "Smarket/db"
	"encoding/json"
	"net/http"
	"sync"
)

type ProductId struct {
	Id int
}

func (p ProductId) GetById(res http.ResponseWriter, req *http.Request) {
	db := DB.Connect()

	defer db.Close()
	res.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)

	Product := make(chan controller.Products, 1)
	wg := &sync.WaitGroup{}

	wg.Add(1)
	go controller.ProductGetId(db, Product, wg, p.Id)

	wg.Wait()
	close(Product)

	product := <-Product

	var data = map[string]interface{}{
		"product": product,
	}

	json.NewEncoder(res).Encode(data)
}

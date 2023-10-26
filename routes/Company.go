package routes

import (
	"Smarket/controller"
	DB "Smarket/db"
	"encoding/json"
	"net/http"
	"sync"
)

type CompanySlug struct {
	Slug string
}

func (c CompanySlug) GetBySlug(res http.ResponseWriter, req *http.Request) {
	db := DB.Connect()

	defer db.Close()
	res.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)

	Company := make(chan []controller.Companies, 1)
	wg := &sync.WaitGroup{}

	wg.Add(1)
	go controller.CompanyGetAllProducts(db, Company, wg, c.Slug)

	wg.Wait()
	close(Company)

	product := <-Company

	var data = map[string]interface{}{
		"Products": product,
	}

	json.NewEncoder(res).Encode(data)
}

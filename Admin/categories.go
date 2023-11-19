package admin

import (
	"Smarket/controller"
	DB "Smarket/db"
	"encoding/json"
	"net/http"
	"sync"
)

func GetCategories(res http.ResponseWriter, req *http.Request) {
	db := DB.Connect()

	defer db.Close()
	res.Header().Set("Access-Control-Allow-Origin", "*")
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)

	category := make(chan []controller.Categories, 1)
	wg := &sync.WaitGroup{}

	wg.Add(1)

	go controller.GetAllCategories(db, category, wg)

	wg.Wait()

	close(category)

	var Categories = map[string]interface{}{
		"Categories": <-category,
	}

	json.NewEncoder(res).Encode(Categories)
}

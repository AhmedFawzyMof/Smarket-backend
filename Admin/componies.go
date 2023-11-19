package admin

import (
	"Smarket/controller"
	DB "Smarket/db"
	"encoding/json"
	"net/http"
	"sync"
)

func GetComponies(res http.ResponseWriter, req *http.Request) {
	db := DB.Connect()

	defer db.Close()
	res.Header().Set("Access-Control-Allow-Origin", "*")
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)

	company := make(chan []controller.Companies, 1)
	wg := &sync.WaitGroup{}

	wg.Add(1)

	go controller.GetAllComponies(db, company, wg)

	wg.Wait()

	close(company)

	var Companies = map[string]interface{}{
		"Companies": <-company,
	}

	json.NewEncoder(res).Encode(Companies)
}

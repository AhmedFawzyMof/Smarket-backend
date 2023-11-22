package admin

import (
	"alwadi/controller"
	DB "alwadi/db"
	"encoding/json"
	"net/http"
	"sync"
)

func GetOffers(res http.ResponseWriter, req *http.Request) {
	db := DB.Connect()

	defer db.Close()
	res.Header().Set("Access-Control-Allow-Origin", "*")
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)

	offers := make(chan []controller.Offers, 1)
	wg := &sync.WaitGroup{}

	wg.Add(1)

	go controller.GetallOffers(db, offers, wg)

	wg.Wait()

	close(offers)

	var Offers = map[string]interface{}{
		"Offers": <-offers,
	}

	json.NewEncoder(res).Encode(Offers)
}

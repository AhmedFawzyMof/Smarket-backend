package admin

import (
	DB "alwadi_markets/db"
	"alwadi_markets/tables"
	"encoding/json"
	"net/http"
	"sync"
)

func GetOrders(res http.ResponseWriter, req *http.Request, params map[string]string) {
	res.WriteHeader(http.StatusOK)
	db := DB.Connect()
	defer db.Close()

	var Order tables.Orders

	orderChan := make(chan []byte, 1)

	wg := &sync.WaitGroup{}

	wg.Add(1)
	go tables.Orders.GetAll(Order, db, orderChan, wg)
	wg.Wait()

	close(orderChan)

	var OrderResponse []tables.Orders

	Response := make(map[string]interface{})

	errors := json.Unmarshal(<-orderChan, &OrderResponse)

	if errors != nil {
		http.Error(res, errors.Error(), http.StatusInternalServerError)
	}

	if err := json.NewEncoder(res).Encode(Response); err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}
}

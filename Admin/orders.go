package admin

import (
	"Smarket/controller"
	DB "Smarket/db"
	"encoding/json"
	"net/http"
	"sync"
)

func GetOrders(res http.ResponseWriter, req *http.Request) {
	db := DB.Connect()

	defer db.Close()
	res.Header().Set("Access-Control-Allow-Origin", "*")
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)

	orders := make(chan controller.Orders, 1)
	wg := &sync.WaitGroup{}

	wg.Add(1)

	go controller.GetAllOrders(db, orders, wg)

	wg.Wait()

	close(orders)

	var Orders = <-orders

	json.NewEncoder(res).Encode(Orders)
}

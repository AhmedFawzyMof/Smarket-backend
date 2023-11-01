package routes

import (
	"Smarket/controller"
	DB "Smarket/db"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

type TheOrder struct {
	Address  map[string]string
	Method   string
	Token    string
	Products []controller.OrderProducts
}

func MakeOrders(res http.ResponseWriter, req *http.Request) {
	db := DB.Connect()

	defer db.Close()
	res.Header().Set("Access-Control-Allow-Origin", "*")
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)

	body, err := io.ReadAll(req.Body)

	if err != nil {
		panic(err.Error())
	}

	var Order TheOrder

	mapData := json.Unmarshal(body, &Order)
	if mapData != nil {
		panic(mapData.Error())
	}
	var method string = Order.Method
	token := Order.Token
	var deleverd int
	var paid int
	var OrderId string

	if method == "cash" {
		deleverd = 0
		paid = 0
	}

	orderForm := map[string]interface{}{
		"deleverd": deleverd,
		"paid":     paid,
		"method":   method,
		"token":    token,
	}

	CheckErrorOrder := controller.CreateOrder(db, orderForm)

	OrderId = fmt.Sprintf("%v", CheckErrorOrder["OrderId"])

	address := Order.Address
	address["UserId"] = fmt.Sprintf("%v", CheckErrorOrder["ID"])

	controller.CreateOrderProduct(db, OrderId, Order.Products)
	controller.CheckAddress(db, address)

	json.NewEncoder(res).Encode(CheckErrorOrder)
}

func OrdersHistory(res http.ResponseWriter, req *http.Request) {
	db := DB.Connect()

	defer db.Close()
	res.Header().Set("Access-Control-Allow-Origin", "*")
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)

	body, err := io.ReadAll(req.Body)

	if err != nil {
		panic(err.Error())
	}

	dataForm := make(map[string]string)
	TheToken := json.Unmarshal(body, &dataForm)
	if TheToken != nil {
		panic(TheToken.Error())
	}

	token := dataForm["authToken"]

	Orders := controller.GetOrders(db, token)

	json.NewEncoder(res).Encode(Orders)
}

func CancelOrder(res http.ResponseWriter, req *http.Request) {
	db := DB.Connect()

	defer db.Close()
	res.Header().Set("Access-Control-Allow-Origin", "*")
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	body, err := io.ReadAll(req.Body)

	if err != nil {
		panic(err.Error())
	}

	dataForm := make(map[string]string)
	TheToken := json.Unmarshal(body, &dataForm)
	if TheToken != nil {
		panic(TheToken.Error())
	}

	order := dataForm["order"]
	token := dataForm["token"]
	confirmed := dataForm["confirmed"]
	i, err := strconv.Atoi(confirmed)
	if err != nil {
		panic(err.Error())
	}

	orderRes := controller.CancelOrder(db, token, order, i)

	json.NewEncoder(res).Encode(orderRes)
}

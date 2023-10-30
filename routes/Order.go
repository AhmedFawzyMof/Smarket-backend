package routes

import (
	"Smarket/controller"
	DB "Smarket/db"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Order struct {
	Token string
}

type TheOrder struct {
	Address  map[string]string
	Method   string
	Token    string
	Products []controller.OrderProducts
}

func MakeOrders(res http.ResponseWriter, req *http.Request) {
	db := DB.Connect()

	defer db.Close()
	res.Header().Set("Access-Control-Allow-Origin", "https://beautiful-ruby-sheath-dress.cyclic.app/")
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

func (o Order) OrdersHistory(res http.ResponseWriter, req *http.Request) {
	db := DB.Connect()

	defer db.Close()
	res.Header().Set("Access-Control-Allow-Origin", "https://beautiful-ruby-sheath-dress.cyclic.app/")
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)

	Orders := controller.GetOrders(db, o.Token)
	var OrderRes = make(map[string]interface{})

	if len(Orders) == 0 {
		OrderRes["Error"] = false
		OrderRes["Message"] = "لا توجد طلبات حتى الآن"
	}
	OrderRes = Orders

	json.NewEncoder(res).Encode(Orders)
}

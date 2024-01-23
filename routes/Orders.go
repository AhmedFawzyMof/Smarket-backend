package routes

import (
	DB "alwadi_markets/db"
	"alwadi_markets/middleware"
	"alwadi_markets/models"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
)

type order struct {
	Products []interface{}     `json:"products"`
	Address  map[string]string `json:"address"`
	Method   string            `json:"method"`
	Token    string            `json:"token"`
}

func OrderHistory(res http.ResponseWriter, req *http.Request, params map[string]string) {
	if req.Method == "POST" {
		res.WriteHeader(http.StatusOK)
		db := DB.Connect()
		defer db.Close()

		body, err := io.ReadAll(req.Body)

		if err != nil {
			middleware.SendError(err, res)
		}

		var orderMap map[string]interface{}

		var Order models.Orders

		orderChan := make(chan []byte, 1)

		wg := &sync.WaitGroup{}

		if err := json.Unmarshal(body, &orderMap); err != nil {
			middleware.SendError(err, res)
		}

		var token string = fmt.Sprintf("%s", orderMap["token"])

		id, err := middleware.VerifyToken(token)

		if err != nil {
			middleware.SendError(err, res)
		}

		Order.User = id

		wg.Add(1)
		go models.Orders.GetForUser(Order, db, orderChan, wg)
		wg.Wait()

		close(orderChan)

		var OrderResponse []models.Orders

		Response := make(map[string]interface{})

		if err := json.Unmarshal(<-orderChan, &OrderResponse); err != nil {
			middleware.SendError(err, res)
		}
		Response["Orders"] = OrderResponse

		if err := json.NewEncoder(res).Encode(Response); err != nil {
			middleware.SendError(err, res)
		}
	}
}

func CancelOrder(res http.ResponseWriter, req *http.Request, params map[string]string) {
	if req.Method == "DELETE" {
		res.WriteHeader(http.StatusOK)
		db := DB.Connect()
		defer db.Close()

		body, err := io.ReadAll(req.Body)

		if err != nil {
			middleware.SendError(err, res)
		}

		var orderMap map[string]string

		var Order models.Orders

		orderChan := make(chan []byte, 1)

		wg := &sync.WaitGroup{}

		if err := json.Unmarshal(body, &orderMap); err != nil {
			middleware.SendError(err, res)
		}

		var token string = orderMap["token"]
		var order string = orderMap["order"]

		id, err := middleware.VerifyToken(token)

		if err != nil {
			middleware.SendError(err, res)
		}

		Order.User = id
		Order.Id = order

		wg.Add(1)
		go models.Orders.Delete(Order, db, orderChan, wg)
		wg.Wait()

		close(orderChan)

		Response := make(map[string]interface{})

		if err := json.Unmarshal(<-orderChan, &Response); err != nil {
			middleware.SendError(err, res)
		}

		if err := json.NewEncoder(res).Encode(Response); err != nil {
			middleware.SendError(err, res)
		}
	}
}

func OrderPage(res http.ResponseWriter, req *http.Request, params map[string]string) {
	if req.Method == "POST" {
		res.WriteHeader(http.StatusOK)
		db := DB.Connect()
		defer db.Close()

		body, err := io.ReadAll(req.Body)

		if err != nil {
			middleware.SendError(err, res)
		}

		var orderMap map[string]string

		var Order models.Orders

		orderChan := make(chan []byte, 1)

		wg := &sync.WaitGroup{}

		if err := json.Unmarshal(body, &orderMap); err != nil {
			middleware.SendError(err, res)
		}

		var token string = orderMap["token"]
		var order string = params["id"]

		id, err := middleware.VerifyToken(token)

		if err != nil {
			middleware.SendError(err, res)
		}

		Order.User = id
		Order.Id = order

		wg.Add(1)
		go models.Orders.OrderDitails(Order, db, orderChan, wg)
		wg.Wait()

		close(orderChan)

		Response := make(map[string]interface{})

		if err := json.Unmarshal(<-orderChan, &Response); err != nil {
			middleware.SendError(err, res)
		}

		if err := json.NewEncoder(res).Encode(Response); err != nil {
			middleware.SendError(err, res)
		}
	}
}

func Order(res http.ResponseWriter, req *http.Request, params map[string]string) {
	if req.Method == "POST" {
		res.WriteHeader(http.StatusOK)
		db := DB.Connect()
		defer db.Close()
		body, err := io.ReadAll(req.Body)

		if err != nil {
			middleware.SendError(err, res)
		}

		var OrderReq map[string]interface{}

		var Orders models.Orders

		if err := json.Unmarshal(body, &OrderReq); err != nil {
			middleware.SendError(err, res)
		}

		var Order order
		var Address models.AddressTable

		middleware.FillStruct(OrderReq, &Order)

		id, err := middleware.VerifyToken(Order.Token)

		if err != nil {
			middleware.SendError(err, res)
		}

		Address.User = id
		Orders.User = id
		Orders.Method = Order.Method

		wg := &sync.WaitGroup{}

		addressChan := make(chan []byte, 1)

		orderId, err := models.Orders.Add(Orders, db)

		if err != nil {
			middleware.SendError(err, res)
		}

		middleware.FillStructInterface(Order.Address, &Address)

		if err := createOrderProducts(Order.Products, orderId, db); err != nil {
			middleware.SendError(err, res)
		}

		wg.Add(1)

		go models.AddressTable.Add(Address, db, addressChan, wg)

		wg.Wait()

		close(addressChan)

		Response := make(map[string]interface{})

		if err := json.Unmarshal(<-addressChan, &Response); err != nil {
			middleware.SendError(err, res)
		}

		if err := json.NewEncoder(res).Encode(Response); err != nil {
			middleware.SendError(err, res)
		}

	}
}

func createOrderProducts(orderProducts []interface{}, orderId string, db *sql.DB) error {
	var Orderproducts []models.OrderProducts
	var Orderproduct models.OrderProducts
	for _, product := range orderProducts {
		middleware.FillStructInterface(product, &Orderproduct)
		Orderproduct.Order = orderId
		Orderproducts = append(Orderproducts, Orderproduct)
	}
	for _, op := range Orderproducts {
		sucsess := models.OrderProducts.Add(op, db)
		if !sucsess {
			return fmt.Errorf("the product didnt inserted")
		}
	}

	return nil
}

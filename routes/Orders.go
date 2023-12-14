package routes

import (
	DB "alwadi_markets/db"
	"alwadi_markets/middleware"
	"alwadi_markets/tables"
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

		var Order tables.Orders

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
		go tables.Orders.GetForUser(Order, db, orderChan, wg)
		wg.Wait()

		close(orderChan)

		var OrderResponse []tables.Orders

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

		var Order tables.Orders

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
		go tables.Orders.Delete(Order, db, orderChan, wg)
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

		var Order tables.Orders
		var OrderProduct tables.OrderProducts

		orderChan := make(chan []byte, 1)
		orderProductChan := make(chan []byte, 1)

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
		OrderProduct.Order = order

		wg.Add(2)
		go tables.Orders.GetForUser(Order, db, orderChan, wg)
		go tables.OrderProducts.GetByOrder(OrderProduct, db, orderProductChan, wg)
		wg.Wait()

		close(orderChan)

		var OrderResponse []tables.Orders
		var OrderProductsResponse []tables.OP

		Response := make(map[string]interface{})

		if err := json.Unmarshal(<-orderChan, &OrderResponse); err != nil {
			middleware.SendError(err, res)
		}

		if err := json.Unmarshal(<-orderProductChan, &OrderProductsResponse); err != nil {
			middleware.SendError(err, res)
		}

		Response["Order"] = OrderResponse
		Response["Products"] = OrderProductsResponse

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

		var Orders tables.Orders

		if err := json.Unmarshal(body, &OrderReq); err != nil {
			middleware.SendError(err, res)
		}

		var Order order
		var Address tables.AddressTable

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

		orderId, err := tables.Orders.Add(Orders, db)

		if err != nil {
			middleware.SendError(err, res)
		}

		middleware.FillStructInterface(Order.Address, &Address)

		if err := createOrderProducts(Order.Products, orderId, db); err != nil {
			middleware.SendError(err, res)
		}

		wg.Add(1)

		go tables.AddressTable.Add(Address, db, addressChan, wg)

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
	var Orderproducts []tables.OrderProducts
	var Orderproduct tables.OrderProducts
	for _, product := range orderProducts {
		middleware.FillStructInterface(product, &Orderproduct)
		Orderproduct.Order = orderId
		Orderproducts = append(Orderproducts, Orderproduct)
	}
	for _, op := range Orderproducts {
		sucsess := tables.OrderProducts.Add(op, db)
		if !sucsess {
			return fmt.Errorf("the product didnt inserted")
		}
	}

	return nil
}

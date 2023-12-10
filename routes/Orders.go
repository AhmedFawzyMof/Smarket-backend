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
			panic(err.Error())
		}

		var orderMap map[string]interface{}

		var Order tables.Orders

		orderChan := make(chan []byte, 1)

		wg := &sync.WaitGroup{}

		Error := json.Unmarshal(body, &orderMap)

		if Error != nil {
			http.Error(res, Error.Error(), http.StatusInternalServerError)
		}

		var token string = fmt.Sprintf("%s", orderMap["token"])

		id, e := middleware.VerifyToken(token)

		if e != nil {
			http.Error(res, e.Error(), http.StatusInternalServerError)
		}

		Order.User = id

		wg.Add(1)
		go tables.Orders.GetForUser(Order, db, orderChan, wg)
		wg.Wait()

		close(orderChan)

		var OrderResponse []tables.Orders

		Response := make(map[string]interface{})

		errors := json.Unmarshal(<-orderChan, &OrderResponse)

		if errors != nil {
			http.Error(res, errors.Error(), http.StatusInternalServerError)
		}
		Response["Orders"] = OrderResponse

		if err := json.NewEncoder(res).Encode(Response); err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
		}
	}
}

func CancelOrder(res http.ResponseWriter, req *http.Request, params map[string]string) {
	if req.Method == "POST" {
		res.WriteHeader(http.StatusOK)
		db := DB.Connect()
		defer db.Close()

		body, err := io.ReadAll(req.Body)

		if err != nil {
			panic(err.Error())
		}

		var orderMap map[string]string

		var Order tables.Orders

		orderChan := make(chan []byte, 1)

		wg := &sync.WaitGroup{}

		Error := json.Unmarshal(body, &orderMap)

		if Error != nil {
			http.Error(res, Error.Error(), http.StatusInternalServerError)
		}

		var token string = orderMap["token"]
		var order string = orderMap["order"]

		id, e := middleware.VerifyToken(token)

		if e != nil {
			http.Error(res, e.Error(), http.StatusInternalServerError)
		}

		Order.User = id
		Order.Id = order

		wg.Add(1)
		go tables.Orders.Delete(Order, db, orderChan, wg)
		wg.Wait()

		close(orderChan)

		Response := make(map[string]interface{})

		errors := json.Unmarshal(<-orderChan, &Response)

		if errors != nil {
			http.Error(res, errors.Error(), http.StatusInternalServerError)
		}

		if err := json.NewEncoder(res).Encode(Response); err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
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
			http.Error(res, err.Error(), http.StatusInternalServerError)
		}

		var orderMap map[string]string

		var Order tables.Orders
		var OrderProduct tables.OrderProducts

		orderChan := make(chan []byte, 1)
		orderProductChan := make(chan []byte, 1)

		wg := &sync.WaitGroup{}

		Error := json.Unmarshal(body, &orderMap)

		if Error != nil {
			http.Error(res, Error.Error(), http.StatusInternalServerError)
		}

		var token string = orderMap["token"]
		var order string = params["id"]

		id, e := middleware.VerifyToken(token)

		if e != nil {
			http.Error(res, e.Error(), http.StatusInternalServerError)
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

		errors := json.Unmarshal(<-orderChan, &OrderResponse)
		erro := json.Unmarshal(<-orderProductChan, &OrderProductsResponse)

		if errors != nil {
			http.Error(res, errors.Error(), http.StatusInternalServerError)
		}

		if erro != nil {
			http.Error(res, erro.Error(), http.StatusInternalServerError)
		}

		Response["Order"] = OrderResponse
		Response["Products"] = OrderProductsResponse

		if err := json.NewEncoder(res).Encode(Response); err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
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
			http.Error(res, err.Error(), http.StatusInternalServerError)
		}

		var OrderReq map[string]interface{}

		var Orders tables.Orders

		Error := json.Unmarshal(body, &OrderReq)

		if Error != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
		}

		var Order order
		var Address tables.AddressTable

		middleware.FillStruct(OrderReq, &Order)

		id, e := middleware.VerifyToken(Order.Token)

		if e != nil {
			http.Error(res, e.Error(), http.StatusInternalServerError)
		}

		Address.User = id
		Orders.User = id
		Orders.Method = Order.Method

		wg := &sync.WaitGroup{}

		addressChan := make(chan []byte, 1)

		orderId, Er := tables.Orders.Add(Orders, db)

		if Er != nil {
			http.Error(res, Er.Error(), http.StatusInternalServerError)
		}

		middleware.FillStructInterface(Order.Address, &Address)

		Err := createOrderProducts(Order.Products, orderId, db)

		if Err != nil {
			http.Error(res, Err.Error(), http.StatusInternalServerError)
		}

		wg.Add(1)

		go tables.AddressTable.Add(Address, db, addressChan, wg)

		wg.Wait()

		close(addressChan)

		Response := make(map[string]interface{})

		errors := json.Unmarshal(<-addressChan, &Response)

		if errors != nil {
			http.Error(res, errors.Error(), http.StatusInternalServerError)
		}

		if err := json.NewEncoder(res).Encode(Response); err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
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

package mange

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"

	DB "alwadimarkets/db"
	"alwadimarkets/middleware"
	"alwadimarkets/models"
)

func GetOrders(res http.ResponseWriter, req *http.Request, params map[string]string) {
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

		if err := json.Unmarshal(body, &orderMap); err != nil {
			middleware.SendError(err, res)
		}

		var token string = fmt.Sprintf("%s", orderMap["auth-token"])

		admin := middleware.CheckIsAdmin(token, db)

		if !admin {
			err := fmt.Errorf("user is not admin")
			middleware.SendError(err, res)
		}

		orderChan := make(chan []byte, 1)

		wg := &sync.WaitGroup{}

		wg.Add(1)
		go models.Orders.GetAll(Order, db, orderChan, wg)
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

func OrderPage(res http.ResponseWriter, req *http.Request, params map[string]string) {
	if req.Method == "POST" {
		res.WriteHeader(http.StatusOK)
		db := DB.Connect()
		defer db.Close()

		body, err := io.ReadAll(req.Body)

		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
		}

		var orderMap map[string]interface{}

		var Order models.Orders

		orderChan := make(chan []byte, 1)

		wg := &sync.WaitGroup{}

		if err := json.Unmarshal(body, &orderMap); err != nil {
			middleware.SendError(err, res)
		}

		var token string = fmt.Sprintf("%s", orderMap["auth-token"])

		admin := middleware.CheckIsAdmin(token, db)

		if !admin {
			err := fmt.Errorf("user is not admin")
			middleware.SendError(err, res)
		}
		var order string = params["id"]

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

func EditOrder(res http.ResponseWriter, req *http.Request, params map[string]string) {
	if req.Method == "PUT" {
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

		var token string = fmt.Sprintf("%s", orderMap["auth-token"])

		admin := middleware.CheckIsAdmin(token, db)

		if !admin {
			err := fmt.Errorf("user is not admin")
			middleware.SendError(err, res)
		}

		Order.Id = fmt.Sprintf("%v", orderMap["order"])
		_, paid := orderMap["paid"]
		if paid {
			var paidFloat float64 = orderMap["paid"].(float64)
			var paidint int = int(paidFloat)
			Order.Paid = paidint
		}
		_, delivered := orderMap["delivered"]
		if delivered {
			var deliveredFloat float64 = orderMap["delivered"].(float64)
			var deliveredint int = int(deliveredFloat)
			Order.Delivered = deliveredint
		}
		_, confirm := orderMap["confirm"]
		if confirm {
			var confirmFloat float64 = orderMap["confirm"].(float64)
			var confirmint int = int(confirmFloat)
			Order.Confirmed = confirmint
		}
		wg.Add(1)
		go models.Orders.Update(Order, db, orderChan, wg)
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

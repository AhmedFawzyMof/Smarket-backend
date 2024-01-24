package adminHandler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"sync"

	DB "alwadimarkets/db"
	"alwadimarkets/middleware"
	"alwadimarkets/models"
)

func GetTypes(res http.ResponseWriter, req *http.Request, params map[string]string) {
	db := DB.Connect()
	res.WriteHeader(http.StatusOK)

	defer db.Close()

	body, err := io.ReadAll(req.Body)

	if err != nil {
		middleware.SendError(err, res)
	}
	var productmap map[string]interface{}

	if err := json.Unmarshal(body, &productmap); err != nil {
		middleware.SendError(err, res)
	}

	var token string = fmt.Sprintf("%s", productmap["auth-token"])

	admin := middleware.CheckIsAdmin(token, db)

	if !admin {
		err := fmt.Errorf("user is not admin")
		middleware.SendError(err, res)
	}

	Types := make(chan []byte, 1)

	wg := &sync.WaitGroup{}

	wg.Add(1)

	go models.ProductType.Get(models.ProductType{}, db, Types, wg)

	wg.Wait()

	close(Types)

	var Type []models.ProductType

	if err := json.Unmarshal(<-Types, &Type); err != nil {
		middleware.SendError(err, res)
	}

	Response := make(map[string]interface{}, 1)

	Response["Types"] = Type

	if err := json.NewEncoder(res).Encode(Response); err != nil {
		middleware.SendError(err, res)
	}
}

func AddTypes(res http.ResponseWriter, req *http.Request, params map[string]string) {
	db := DB.Connect()
	res.WriteHeader(http.StatusOK)

	defer db.Close()

	body, err := io.ReadAll(req.Body)

	if err != nil {
		panic(err.Error())
	}
	var productmap map[string]interface{}

	var ProductType models.ProductType

	if err := json.Unmarshal(body, &productmap); err != nil {
		middleware.SendError(err, res)
	}

	var token string = fmt.Sprintf("%s", productmap["auth-token"])

	admin := middleware.CheckIsAdmin(token, db)

	if !admin {
		err := fmt.Errorf("user is not admin")
		middleware.SendError(err, res)
	}

	// body data
	ProductType.Product = middleware.ConvertToInt(productmap["product"], res)
	ProductType.Portion = middleware.ConvertToInt(productmap["portion"], res)
	ProductType.Price = middleware.ConvertToInt(productmap["price"], res)
	ProductType.Offer = middleware.ConvertToInt(productmap["offer"], res)
	ProductType.Uint = fmt.Sprintf("%s", productmap["uint"])

	Types := make(chan []byte, 1)

	wg := &sync.WaitGroup{}

	wg.Add(1)
	go models.ProductType.Add(ProductType, db, Types, wg)
	wg.Wait()

	close(Types)
	var products map[string]interface{}

	if err := json.Unmarshal(<-Types, &products); err != nil {
		middleware.SendError(err, res)
	}

	if err := json.NewEncoder(res).Encode(products); err != nil {
		middleware.SendError(err, res)
	}
}

func DeleteTypes(res http.ResponseWriter, req *http.Request, params map[string]string) {
	if req.Method == "DELETE" {
		db := DB.Connect()
		res.WriteHeader(http.StatusOK)

		defer db.Close()

		body, err := io.ReadAll(req.Body)

		if err != nil {
			panic(err.Error())
		}
		var productmap map[string]interface{}

		var ProductType models.ProductType

		if err := json.Unmarshal(body, &productmap); err != nil {
			middleware.SendError(err, res)
		}

		var token string = fmt.Sprintf("%s", productmap["auth-token"])

		admin := middleware.CheckIsAdmin(token, db)

		if !admin {
			err := fmt.Errorf("user is not admin")
			middleware.SendError(err, res)
		}

		// body data
		id, _ := strconv.Atoi(fmt.Sprintf("%s", productmap["id"]))
		ProductType.Id = id

		Types := make(chan []byte, 1)

		wg := &sync.WaitGroup{}

		wg.Add(1)
		go models.ProductType.Delete(ProductType, db, Types, wg)
		wg.Wait()

		close(Types)
		var products map[string]interface{}

		errors := json.Unmarshal(<-Types, &products)

		if errors != nil {
			http.Error(res, errors.Error(), http.StatusInternalServerError)
		}

		if err := json.NewEncoder(res).Encode(products); err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
		}
	}
}

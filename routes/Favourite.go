package routes

import (
	DB "alwadi_markets/db"
	"alwadi_markets/middleware"
	"alwadi_markets/tables"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
)

func AddFavourite(res http.ResponseWriter, req *http.Request, params map[string]string) {
	if req.Method == "POST" {
		db := DB.Connect()
		res.WriteHeader(http.StatusOK)

		defer db.Close()

		body, err := io.ReadAll(req.Body)

		if err != nil {
			panic(err.Error())
		}
		var favMap map[string]interface{}

		var Favourite tables.Favourite

		Error := json.Unmarshal(body, &favMap)

		if Error != nil {
			http.Error(res, Error.Error(), http.StatusInternalServerError)
		}

		var token string = fmt.Sprintf("%s", favMap["token"])

		id, e := middleware.VerifyToken(token)
		if e != nil {
			http.Error(res, e.Error(), http.StatusInternalServerError)
		}

		var productfloat64 float64 = favMap["product"].(float64)

		var product int = int(productfloat64)

		Favourite.Product = product
		Favourite.User = id

		FavChan := make(chan []byte, 1)

		wg := &sync.WaitGroup{}

		wg.Add(1)
		go tables.Favourite.Add(Favourite, db, FavChan, wg)
		wg.Wait()

		close(FavChan)

		var FavResponse map[string]interface{}

		errors := json.Unmarshal(<-FavChan, &FavResponse)

		if errors != nil {
			http.Error(res, errors.Error(), http.StatusInternalServerError)
		}

		if err := json.NewEncoder(res).Encode(FavResponse); err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
		}
	}
}

func GetFavourite(res http.ResponseWriter, req *http.Request, params map[string]string) {
	if req.Method == "POST" {
		db := DB.Connect()
		res.WriteHeader(http.StatusOK)

		defer db.Close()

		body, err := io.ReadAll(req.Body)

		if err != nil {
			panic(err.Error())
		}
		var favMap map[string]interface{}

		var Favourite tables.Favourite

		Error := json.Unmarshal(body, &favMap)

		if Error != nil {
			http.Error(res, Error.Error(), http.StatusInternalServerError)
		}

		var token string = fmt.Sprintf("%s", favMap["token"])

		id, e := middleware.VerifyToken(token)
		if e != nil {
			http.Error(res, e.Error(), http.StatusInternalServerError)
		}

		Favourite.User = id

		FavChan := make(chan []byte, 1)

		wg := &sync.WaitGroup{}

		wg.Add(1)
		go tables.Favourite.Get(Favourite, db, FavChan, wg)
		wg.Wait()

		close(FavChan)

		var FavProducts []tables.Product

		errors := json.Unmarshal(<-FavChan, &FavProducts)

		if errors != nil {
			http.Error(res, errors.Error(), http.StatusInternalServerError)
		}

		Response := make(map[string]interface{}, 1)

		Response["Products"] = FavProducts

		if err := json.NewEncoder(res).Encode(Response); err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
		}
	}
}

func DeleteFavourite(res http.ResponseWriter, req *http.Request, params map[string]string) {
	if req.Method == "DELETE" {
		db := DB.Connect()
		res.WriteHeader(http.StatusOK)

		defer db.Close()

		body, err := io.ReadAll(req.Body)

		if err != nil {
			panic(err.Error())
		}
		var favMap map[string]interface{}

		var Favourite tables.Favourite

		Error := json.Unmarshal(body, &favMap)

		if Error != nil {
			http.Error(res, Error.Error(), http.StatusInternalServerError)
		}

		var token string = fmt.Sprintf("%s", favMap["token"])

		id, e := middleware.VerifyToken(token)
		if e != nil {
			http.Error(res, e.Error(), http.StatusInternalServerError)
		}

		var productfloat64 float64 = favMap["product"].(float64)

		var product int = int(productfloat64)

		Favourite.Product = product
		Favourite.User = id

		FavChan := make(chan []byte, 1)

		wg := &sync.WaitGroup{}

		wg.Add(1)
		go tables.Favourite.Delete(Favourite, db, FavChan, wg)
		wg.Wait()

		close(FavChan)

		var FavResponse map[string]interface{}

		errors := json.Unmarshal(<-FavChan, &FavResponse)

		if errors != nil {
			http.Error(res, errors.Error(), http.StatusInternalServerError)
		}

		if err := json.NewEncoder(res).Encode(FavResponse); err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
		}
	}
}

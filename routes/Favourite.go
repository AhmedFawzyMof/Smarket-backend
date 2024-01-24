package routes

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

func AddFavourite(res http.ResponseWriter, req *http.Request, params map[string]string) {
	if req.Method == "POST" {
		db := DB.Connect()
		res.WriteHeader(http.StatusOK)

		defer db.Close()

		body, err := io.ReadAll(req.Body)

		if err != nil {
			middleware.SendError(err, res)
		}
		var favMap map[string]interface{}

		var Favourite models.Favourite

		if err := json.Unmarshal(body, &favMap); err != nil {
			middleware.SendError(err, res)
		}

		var token string = fmt.Sprintf("%s", favMap["token"])

		id, err := middleware.VerifyToken(token)

		if err != nil {
			middleware.SendError(err, res)
		}

		var productfloat64 float64 = favMap["product"].(float64)

		var product int = int(productfloat64)

		Favourite.Product = product
		Favourite.User = id

		FavChan := make(chan []byte, 1)

		wg := &sync.WaitGroup{}

		wg.Add(1)
		go models.Favourite.Add(Favourite, db, FavChan, wg)
		wg.Wait()

		close(FavChan)

		var FavResponse map[string]interface{}

		if err := json.Unmarshal(<-FavChan, &FavResponse); err != nil {
			middleware.SendError(err, res)
		}

		if err := json.NewEncoder(res).Encode(FavResponse); err != nil {
			middleware.SendError(err, res)
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
			middleware.SendError(err, res)
		}

		var favMap map[string]interface{}

		var Favourite models.Favourite

		if err := json.Unmarshal(body, &favMap); err != nil {
			middleware.SendError(err, res)
		}

		var token string = fmt.Sprintf("%s", favMap["token"])

		id, err := middleware.VerifyToken(token)
		if err != nil {
			middleware.SendError(err, res)
		}

		Favourite.User = id

		FavChan := make(chan []byte, 1)

		wg := &sync.WaitGroup{}

		wg.Add(1)
		go models.Favourite.Get(Favourite, db, FavChan, wg)
		wg.Wait()

		close(FavChan)

		var FavProducts []models.Product

		if err := json.Unmarshal(<-FavChan, &FavProducts); err != nil {
			middleware.SendError(err, res)
		}

		Response := make(map[string]interface{}, 1)

		Response["Products"] = FavProducts

		if err := json.NewEncoder(res).Encode(Response); err != nil {
			middleware.SendError(err, res)
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
			middleware.SendError(err, res)
		}
		var favMap map[string]interface{}

		var Favourite models.Favourite

		if err := json.Unmarshal(body, &favMap); err != nil {
			middleware.SendError(err, res)
		}

		var token string = fmt.Sprintf("%s", favMap["token"])

		id, err := middleware.VerifyToken(token)

		if err != nil {
			middleware.SendError(err, res)
		}

		var productfloat64 float64 = favMap["product"].(float64)

		var product int = int(productfloat64)

		Favourite.Product = product
		Favourite.User = id

		FavChan := make(chan []byte, 1)

		wg := &sync.WaitGroup{}

		wg.Add(1)
		go models.Favourite.Delete(Favourite, db, FavChan, wg)
		wg.Wait()

		close(FavChan)

		var FavResponse map[string]interface{}

		if err := json.Unmarshal(<-FavChan, &FavResponse); err != nil {
			middleware.SendError(err, res)
		}

		if err := json.NewEncoder(res).Encode(FavResponse); err != nil {
			middleware.SendError(err, res)
		}
	}
}

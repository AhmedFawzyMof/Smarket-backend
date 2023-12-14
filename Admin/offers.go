package admin

import (
	DB "alwadi_markets/db"
	"alwadi_markets/middleware"
	"alwadi_markets/tables"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"sync"
)

func GetOffers(res http.ResponseWriter, req *http.Request, params map[string]string) {
	db := DB.Connect()
	res.WriteHeader(http.StatusOK)

	defer db.Close()
	body, err := io.ReadAll(req.Body)

	if err != nil {
		panic(err.Error())
	}

	var offersMap map[string]interface{}

	if err := json.Unmarshal(body, &offersMap); err != nil {
		middleware.SendError(err, res)
	}

	var token string = fmt.Sprintf("%s", offersMap["auth-token"])

	admin := middleware.CheckIsAdmin(token, db)

	if !admin {
		err := fmt.Errorf("user is not admin")
		middleware.SendError(err, res)
	}

	Offers := make(chan []byte, 1)

	wg := &sync.WaitGroup{}

	wg.Add(1)

	go tables.Offer.Get(tables.Offer{}, db, Offers, wg)

	wg.Wait()

	close(Offers)

	var offers []tables.Offer

	if err := json.Unmarshal(<-Offers, &offers); err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}

	Response := make(map[string]interface{}, 1)

	Response["Offers"] = offers

	if err := json.NewEncoder(res).Encode(Response); err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}
}

func AddOffers(res http.ResponseWriter, req *http.Request, params map[string]string) {
	db := DB.Connect()
	res.WriteHeader(http.StatusOK)

	defer db.Close()

	body, err := io.ReadAll(req.Body)

	if err != nil {
		panic(err.Error())
	}

	var offersMap map[string]interface{}

	var Offer tables.Offer

	if err := json.Unmarshal(body, &offersMap); err != nil {
		middleware.SendError(err, res)
	}

	var token string = fmt.Sprintf("%s", offersMap["auth-token"])

	admin := middleware.CheckIsAdmin(token, db)

	if !admin {
		err := fmt.Errorf("user is not admin")
		middleware.SendError(err, res)
	}

	productId, Err := strconv.Atoi(fmt.Sprintf("%s", offersMap["product"]))

	if Err != nil {
		middleware.SendError(Err, res)
	}

	Offer.Product = productId
	Offer.Image = fmt.Sprintf("%s", offersMap["image"])

	Offers := make(chan []byte, 1)

	wg := &sync.WaitGroup{}

	wg.Add(1)

	go tables.Offer.Add(Offer, db, Offers, wg)

	wg.Wait()

	close(Offers)

	var offers []tables.Offer

	if err := json.Unmarshal(<-Offers, &offers); err != nil {
		middleware.SendError(err, res)
	}

	Response := make(map[string]interface{}, 1)

	Response["Offers"] = offers

	if err := json.NewEncoder(res).Encode(Response); err != nil {
		middleware.SendError(err, res)
	}
}

func DeleteOffers(res http.ResponseWriter, req *http.Request, params map[string]string) {
	if req.Method == "DELETE" {
		db := DB.Connect()
		res.WriteHeader(http.StatusOK)

		defer db.Close()
		body, err := io.ReadAll(req.Body)

		if err != nil {
			panic(err.Error())
		}

		var offersMap map[string]interface{}

		var Offer tables.Offer

		if err := json.Unmarshal(body, &offersMap); err != nil {
			middleware.SendError(err, res)
		}

		var token string = fmt.Sprintf("%s", offersMap["auth-token"])

		admin := middleware.CheckIsAdmin(token, db)

		if !admin {
			err := fmt.Errorf("user is not admin")
			middleware.SendError(err, res)
		}

		idFloat, ok := offersMap["id"].(float64)
		if !ok {
			fmt.Println("Error: id is not a float64")
			return
		}

		idInt := int(idFloat)
		Offer.Id = idInt

		Offers := make(chan []byte, 1)
		wg := &sync.WaitGroup{}

		wg.Add(1)

		go tables.Offer.Delete(Offer, db, Offers, wg)

		wg.Wait()

		close(Offers)

		var offers []tables.Offer

		Err := json.Unmarshal(<-Offers, &offers)

		if Err != nil {
			http.Error(res, Err.Error(), http.StatusInternalServerError)
		}

		Response := make(map[string]interface{}, 1)

		Response["Offers"] = offers

		if err := json.NewEncoder(res).Encode(Response); err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
		}
	}
}

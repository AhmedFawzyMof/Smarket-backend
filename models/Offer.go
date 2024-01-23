package models

import (
	"alwadi_markets/middleware"
	"database/sql"
	"encoding/json"
	"fmt"
	"sync"
)

type Offer struct {
	Id      int
	Product int
	Image   string
}

func (o Offer) Add(db *sql.DB, response chan []byte, wg *sync.WaitGroup) {
	defer wg.Done()
	Response := make(map[string]interface{})
	_, err := db.Exec("INSERT INTO `Offer`(`product`, `image`) VALUES (?, ?)", o.Product, o.Image)
	if err != nil {
		Response["Error"] = true

		middleware.SendResponse(response, Response)
		return
	}

	Response["Error"] = false

	middleware.SendResponse(response, Response)
}

func (o Offer) Get(db *sql.DB, response chan []byte, wg *sync.WaitGroup) {
	var Offers []Offer
	offers, err := db.Query("SELECT * FROM Offer")

	if err != nil {
		panic(err.Error())
	}

	defer offers.Close()

	for offers.Next() {
		var Offer Offer

		if err = offers.Scan(&Offer.Id, &Offer.Product, &Offer.Image); err != nil {
			panic(err.Error())
		}

		Offers = append(Offers, Offer)
	}

	offer, err := json.Marshal(Offers)
	if err != nil {
		fmt.Println(err.Error())
	}

	response <- offer
	wg.Done()
}

func (o Offer) Delete(db *sql.DB, response chan []byte, wg *sync.WaitGroup) {
	defer wg.Done()
	Response := make(map[string]interface{})
	_, err := db.Exec("DELETE FROM Offer WHERE id = ?", o.Id)

	if err != nil {
		Response["Error"] = true

		middleware.SendResponse(response, Response)
		return
	}

	Response["Error"] = false

	middleware.SendResponse(response, Response)
}

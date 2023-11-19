package controller

import (
	"database/sql"
	"sync"
)

type offer struct {
	id      int
	product int
	image   []byte
}

type Offers map[string]interface{}

func OfferGetAll(db *sql.DB, responseChan chan []Offers, wg *sync.WaitGroup) {
	var Offers []Offers

	Select, err := db.Query("SELECT * FROM `Offer`")

	if err != nil {
		panic(err.Error())
	}

	defer Select.Close()

	for Select.Next() {
		var offer offer

		if err := Select.Scan(&offer.id, &offer.product, &offer.image); err != nil {
			panic(err.Error())
		}

		theOffer := map[string]interface{}{
			"id":      offer.id,
			"product": offer.product,
			"image":   offer.image,
		}

		Offers = append(Offers, theOffer)
	}

	responseChan <- Offers
	wg.Done()

}

func GetallOffers(db *sql.DB, offerChan chan []Offers, wg *sync.WaitGroup) {
	var Offers []Offers

	Select, err := db.Query("SELECT * FROM `Offer`")

	if err != nil {
		panic(err.Error())
	}

	defer Select.Close()

	for Select.Next() {
		var offer offer

		if err := Select.Scan(&offer.id, &offer.product, &offer.image); err != nil {
			panic(err.Error())
		}

		theOffer := map[string]interface{}{
			"id":      offer.id,
			"product": offer.product,
			"image":   offer.image,
		}

		Offers = append(Offers, theOffer)
	}

	offerChan <- Offers
	wg.Done()
}

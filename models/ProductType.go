package models

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"sync"

	"alwadi_markets/middleware"
)

type ProductType struct {
	Id      int
	Product int
	Portion int
	Uint    string
	Price   int
	Offer   int
}

func (pt ProductType) Get(db *sql.DB, response chan []byte, wg *sync.WaitGroup) {
	var ProductTypes []ProductType
	types, err := db.Query("SELECT * FROM `ProductType`")

	if err != nil {
		panic(err.Error())
	}

	defer types.Close()

	for types.Next() {
		var producttype ProductType

		if err := types.Scan(&producttype.Id, &producttype.Product, &producttype.Portion, &producttype.Uint, &producttype.Price, &producttype.Offer); err != nil {
			panic(err.Error())
		}

		ProductTypes = append(ProductTypes, producttype)
	}

	res, err := json.Marshal(ProductTypes)

	if err != nil {
		fmt.Println(err)
	}

	response <- res
	wg.Done()
}

func (pt ProductType) Add(db *sql.DB, response chan []byte, wg *sync.WaitGroup) {
	defer wg.Done()
	Response := make(map[string]interface{})
	_, err := db.Exec("INSERT INTO `ProductType`(`product`, `portion`, `uint`, `price`, `offer`) VALUES (?, ?, ?, ?, ?)", pt.Product, pt.Portion, pt.Uint, pt.Price, pt.Offer)
	if err != nil {
		Response["Error"] = true

		res, err := json.Marshal(Response)

		if err != nil {
			fmt.Println(err)
		}

		response <- res
		wg.Done()
	}

	Response["Error"] = false

	middleware.SendResponse(response, Response)

}

func (pt ProductType) Delete(db *sql.DB, response chan []byte, wg *sync.WaitGroup) {
	defer wg.Done()
	Response := make(map[string]interface{})
	_, err := db.Exec("DELETE FROM ProductType WHERE id = ?", pt.Id)

	if err != nil {
		Response["Error"] = true

		middleware.SendResponse(response, Response)
		return
	}

	Response["Error"] = false

	middleware.SendResponse(response, Response)
}

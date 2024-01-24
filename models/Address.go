package models

import (
	"database/sql"
	"fmt"
	"sync"

	"alwadi_markets/middleware"
)

type AddressTable struct {
	Building    string `json:"building"`
	City        string `json:"city"`
	Floor       string `json:"floor"`
	Governorate string `json:"governorate"`
	Street      string `json:"street"`
	User        string `json:"user"`
	Id          int    `json:"id"`
}

func (a AddressTable) Add(db *sql.DB, response chan []byte, wg *sync.WaitGroup) {
	defer wg.Done()
	Response := make(map[string]interface{})

	FindEmail := db.QueryRow("SELECT user FROM Address WHERE user = ?", a.User)

	var User string

	Err := FindEmail.Scan(&User)
	if Err != nil {

		_, err := db.Exec("INSERT INTO `Address`(`user`, `governorate`, `city`, `street`, `building`, `floor`) VALUES (?, ?, ?, ?, ?, ?)", a.User, a.Governorate, a.City, a.Street, a.Building, a.Floor)
		if err != nil {
			Response["Error"] = true
			fmt.Println(err, "err")
			middleware.SendResponse(response, Response)
			return
		}

		Response["Error"] = false
		fmt.Println(err, "err")

		middleware.SendResponse(response, Response)
		return
	}

	Response["Error"] = false

	middleware.SendResponse(response, Response)
}

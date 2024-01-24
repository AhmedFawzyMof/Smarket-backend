package models

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"sync"

	"alwadi_markets/middleware"
)

type Company struct {
	Name string
}

func (c Company) Add(db *sql.DB, response chan []byte, wg *sync.WaitGroup) {
	defer wg.Done()
	Response := make(map[string]interface{})
	_, err := db.Exec("INSERT INTO `Companies`(`name`) VALUES (?)", c.Name)
	if err != nil {
		Response["Error"] = true

		middleware.SendResponse(response, Response)
		return
	}

	Response["Error"] = false

	middleware.SendResponse(response, Response)
}

func (c Company) Get(db *sql.DB, response chan []byte, wg *sync.WaitGroup) {
	defer wg.Done()
	var Response []Company
	companies, err := db.Query("SELECT * FROM Companies")
	if err != nil {
		fmt.Println(err.Error())
	}

	defer companies.Close()

	for companies.Next() {
		var company Company

		if err := companies.Scan(&company.Name); err != nil {
			panic(err.Error())
		}

		Response = append(Response, company)
	}
	res, err := json.Marshal(Response)
	if err != nil {
		fmt.Println(err.Error())
	}
	response <- res
}

func (c Company) Update(db *sql.DB, response chan []byte, wg *sync.WaitGroup, name string) {
	defer wg.Done()
	Response := make(map[string]interface{})
	_, err := db.Exec("UPDATE `Companies` SET `name`=? WHERE name=?", c.Name, name)
	if err != nil {
		Response["Error"] = true

		middleware.SendResponse(response, Response)
		return
	}

	Response["Error"] = false

	middleware.SendResponse(response, Response)
}

func (c Company) Delete(db *sql.DB, response chan []byte, wg *sync.WaitGroup) {
	defer wg.Done()
	Response := make(map[string]interface{})
	_, err := db.Exec("DELETE FROM `Companies` WHERE name=?", c.Name)
	if err != nil {
		Response["Error"] = true

		middleware.SendResponse(response, Response)
		return
	}

	Response["Error"] = false

	middleware.SendResponse(response, Response)
}

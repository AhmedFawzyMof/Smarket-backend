package models

import (
	"alwadi_markets/middleware"
	"database/sql"
	"encoding/json"
	"fmt"
	"sync"
)

type Category struct {
	Name  string
	Image string
}

func (c Category) Add(db *sql.DB, response chan []byte, wg *sync.WaitGroup) {
	defer wg.Done()
	Response := make(map[string]interface{})
	_, err := db.Exec("INSERT INTO `Categories`(`name`, `image`) VALUES (?, ?)", c.Name, c.Image)
	if err != nil {
		Response["Error"] = true

		middleware.SendResponse(response, Response)
		return
	}

	Response["Error"] = false

	middleware.SendResponse(response, Response)
}

func (c Category) Get(db *sql.DB, response chan []byte, wg *sync.WaitGroup) {
	defer wg.Done()
	var Response []Category
	categories, err := db.Query("SELECT * FROM Categories")
	if err != nil {
		fmt.Println(err.Error())
	}

	defer categories.Close()

	for categories.Next() {
		var category Category

		if err := categories.Scan(&category.Name, &category.Image); err != nil {
			panic(err.Error())
		}

		Response = append(Response, category)
	}
	res, err := json.Marshal(Response)
	if err != nil {
		fmt.Println(err.Error())
	}
	response <- res
}

func (c Category) Update(db *sql.DB, response chan []byte, wg *sync.WaitGroup, name string) {
	defer wg.Done()
	Response := make(map[string]interface{})
	_, err := db.Exec("UPDATE `Categories` SET `name`=? WHERE name=?", c.Name, name)
	if err != nil {
		Response["Error"] = true

		middleware.SendResponse(response, Response)
		return
	}

	Response["Error"] = false

	middleware.SendResponse(response, Response)
}

func (c Category) Delete(db *sql.DB, response chan []byte, wg *sync.WaitGroup) {
	defer wg.Done()
	Response := make(map[string]interface{})
	_, err := db.Exec("DELETE FROM `Categories` WHERE name=?", c.Name)
	if err != nil {
		Response["Error"] = true

		middleware.SendResponse(response, Response)
		return
	}

	Response["Error"] = false

	middleware.SendResponse(response, Response)
}

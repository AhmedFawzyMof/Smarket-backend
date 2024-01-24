package models

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"sync"

	"alwadimarkets/middleware"
)

type SubCategory struct {
	Name     string
	Category string
	Image    string
}

func (sc SubCategory) Add(db *sql.DB, response chan []byte, wg *sync.WaitGroup) {
	defer wg.Done()
	Response := make(map[string]interface{})
	_, err := db.Exec("INSERT INTO `SubCategories`(`name`, `category`, `image`) VALUES (?, ?, ?)", sc.Name, sc.Category, sc.Image)
	if err != nil {
		Response["Error"] = true
		middleware.SendResponse(response, Response)
		return
	}

	Response["Error"] = false

	middleware.SendResponse(response, Response)
}

func (sc SubCategory) GetByCategory(db *sql.DB, response chan []byte, wg *sync.WaitGroup) {
	defer wg.Done()
	var SubCategories []SubCategory

	subcategory, err := db.Query("SELECT * FROM `SubCategories` WHERE category =?", sc.Category)

	if err != nil {
		panic(err.Error())
	}

	defer subcategory.Close()

	for subcategory.Next() {
		var sub SubCategory

		if err := subcategory.Scan(&sub.Name, &sub.Category, &sub.Image); err != nil {
			panic(err.Error())
		}

		SubCategories = append(SubCategories, sub)
	}

	res, err := json.Marshal(SubCategories)

	if err != nil {
		fmt.Println(err)
	}

	response <- res
}

func (sc SubCategory) Get(db *sql.DB, response chan []byte, wg *sync.WaitGroup) {
	defer wg.Done()
	var Response []SubCategory
	categories, err := db.Query("SELECT * FROM SubCategories")
	if err != nil {
		fmt.Println(err.Error())
	}

	defer categories.Close()

	for categories.Next() {
		var category SubCategory

		if err := categories.Scan(&category.Name, &category.Category, &category.Image); err != nil {
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

func (sc SubCategory) Update(db *sql.DB, response chan []byte, wg *sync.WaitGroup, name string) {
	defer wg.Done()
	Response := make(map[string]interface{})
	_, err := db.Exec("UPDATE `SubCategories` SET `name`=?, category=?, WHERE name=?", sc.Name, sc.Category, name)
	if err != nil {
		Response["Error"] = true

		middleware.SendResponse(response, Response)
		return
	}

	Response["Error"] = false

	middleware.SendResponse(response, Response)
}

func (sc SubCategory) Delete(db *sql.DB, response chan []byte, wg *sync.WaitGroup) {
	defer wg.Done()
	Response := make(map[string]interface{})
	_, err := db.Exec("DELETE FROM `SubCategories` WHERE name=?", sc.Name)
	if err != nil {
		Response["Error"] = true

		middleware.SendResponse(response, Response)
		return
	}

	Response["Error"] = false

	middleware.SendResponse(response, Response)
}

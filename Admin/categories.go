package admin

import (
	DB "alwadi_markets/db"
	"alwadi_markets/tables"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
)

func GetCategories(res http.ResponseWriter, req *http.Request, params map[string]string) {
	db := DB.Connect()
	res.WriteHeader(http.StatusOK)

	defer db.Close()

	Categories := make(chan []byte, 1)

	wg := &sync.WaitGroup{}

	wg.Add(1)

	go tables.Category.Get(tables.Category{}, db, Categories, wg)

	wg.Wait()

	close(Categories)

	var category []tables.Category

	err := json.Unmarshal(<-Categories, &category)

	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}

	Response := make(map[string]interface{}, 1)

	Response["Categories"] = category

	if err := json.NewEncoder(res).Encode(Response); err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}
}

func AddCategories(res http.ResponseWriter, req *http.Request, params map[string]string) {
	db := DB.Connect()
	res.WriteHeader(http.StatusOK)

	defer db.Close()

	body, err := io.ReadAll(req.Body)

	if err != nil {
		panic(err.Error())
	}

	var categoriesMap map[string]interface{}

	var Category tables.Category

	Error := json.Unmarshal(body, &categoriesMap)

	if Error != nil {
		http.Error(res, Error.Error(), http.StatusInternalServerError)
	}

	Category.Name = fmt.Sprintf("%s", categoriesMap["name"])
	Category.Image = fmt.Sprintf("%s", categoriesMap["image"])

	Categories := make(chan []byte, 1)

	wg := &sync.WaitGroup{}

	wg.Add(1)

	go tables.Category.Add(Category, db, Categories, wg)

	wg.Wait()

	close(Categories)

	var category []tables.Category

	Err := json.Unmarshal(<-Categories, &category)

	if Err != nil {
		http.Error(res, Err.Error(), http.StatusInternalServerError)
	}

	Response := make(map[string]interface{}, 1)

	Response["Categories"] = category

	if err := json.NewEncoder(res).Encode(Response); err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}
}

func EditCategories(res http.ResponseWriter, req *http.Request, params map[string]string) {
	if req.Method == "PUT" {
		db := DB.Connect()
		res.WriteHeader(http.StatusOK)

		defer db.Close()

		body, err := io.ReadAll(req.Body)

		if err != nil {
			panic(err.Error())
		}
		var categoriesMap map[string]interface{}
		var Category tables.Category

		Error := json.Unmarshal(body, &categoriesMap)

		if Error != nil {
			http.Error(res, Error.Error(), http.StatusInternalServerError)
		}

		Category.Name = fmt.Sprintf("%s", categoriesMap["name"])
		var name string = fmt.Sprintf("%s", categoriesMap["oldname"])

		Categories := make(chan []byte, 1)

		wg := &sync.WaitGroup{}

		wg.Add(1)

		go tables.Category.Update(Category, db, Categories, wg, name)

		wg.Wait()

		close(Categories)

		var category []tables.Category

		Err := json.Unmarshal(<-Categories, &category)

		if Err != nil {
			http.Error(res, Err.Error(), http.StatusInternalServerError)
		}

		Response := make(map[string]interface{}, 1)

		Response["Categories"] = category

		if err := json.NewEncoder(res).Encode(Response); err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
		}
	}
}

func DeleteCategories(res http.ResponseWriter, req *http.Request, params map[string]string) {
	if req.Method == "DELETE" {
		db := DB.Connect()
		res.WriteHeader(http.StatusOK)

		defer db.Close()

		body, err := io.ReadAll(req.Body)

		if err != nil {
			panic(err.Error())
		}

		var Category tables.Category

		Error := json.Unmarshal(body, &Category)

		if Error != nil {
			http.Error(res, Error.Error(), http.StatusInternalServerError)
		}

		Categories := make(chan []byte, 1)

		wg := &sync.WaitGroup{}

		wg.Add(1)

		go tables.Category.Delete(Category, db, Categories, wg)

		wg.Wait()

		close(Categories)

		var category []tables.Category

		Err := json.Unmarshal(<-Categories, &category)

		if Err != nil {
			http.Error(res, Err.Error(), http.StatusInternalServerError)
		}

		Response := make(map[string]interface{}, 1)

		Response["Categories"] = category

		if err := json.NewEncoder(res).Encode(Response); err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
		}
	}
}

package admin

import (
	DB "alwadi_markets/db"
	"alwadi_markets/middleware"
	"alwadi_markets/models"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
)

func GetCategories(res http.ResponseWriter, req *http.Request, params map[string]string) {
	res.WriteHeader(http.StatusOK)

	db := DB.Connect()
	defer db.Close()

	body, err := io.ReadAll(req.Body)

	if err != nil {
		middleware.SendError(err, res)
	}

	var categoriesMap map[string]interface{}

	if err := json.Unmarshal(body, &categoriesMap); err != nil {
		middleware.SendError(err, res)
	}

	var token string = fmt.Sprintf("%v", categoriesMap["auth-token"])

	admin := middleware.CheckIsAdmin(token, db)

	if !admin {
		err := fmt.Errorf("user is not admin")
		middleware.SendError(err, res)
	}

	Categories := make(chan []byte, 1)

	wg := &sync.WaitGroup{}

	wg.Add(1)

	go models.Category.Get(models.Category{}, db, Categories, wg)

	wg.Wait()

	close(Categories)

	var category []models.Category

	if err := json.Unmarshal(<-Categories, &category); err != nil {
		middleware.SendError(err, res)
	}

	Response := make(map[string]interface{}, 1)

	Response["Categories"] = category

	if err := json.NewEncoder(res).Encode(Response); err != nil {
		middleware.SendError(err, res)
	}
}

func AddCategories(res http.ResponseWriter, req *http.Request, params map[string]string) {
	db := DB.Connect()
	res.WriteHeader(http.StatusOK)

	defer db.Close()

	body, err := io.ReadAll(req.Body)

	if err != nil {
		middleware.SendError(err, res)
	}

	var categoriesMap map[string]interface{}

	var Category models.Category

	if err := json.Unmarshal(body, &categoriesMap); err != nil {
		middleware.SendError(err, res)
	}

	var token string = fmt.Sprintf("%v", categoriesMap["auth-token"])

	admin := middleware.CheckIsAdmin(token, db)

	if !admin {
		err := fmt.Errorf("user is not admin")
		middleware.SendError(err, res)
	}

	Category.Name = fmt.Sprintf("%s", categoriesMap["name"])
	Category.Image = fmt.Sprintf("%s", categoriesMap["image"])

	Categories := make(chan []byte, 1)

	wg := &sync.WaitGroup{}

	wg.Add(1)

	go models.Category.Add(Category, db, Categories, wg)

	wg.Wait()

	close(Categories)

	var category []models.Category

	if err := json.Unmarshal(<-Categories, &category); err != nil {
		middleware.SendError(err, res)
	}

	Response := make(map[string]interface{}, 1)

	Response["Categories"] = category

	if err := json.NewEncoder(res).Encode(Response); err != nil {
		middleware.SendError(err, res)
	}
}

func EditCategories(res http.ResponseWriter, req *http.Request, params map[string]string) {
	if req.Method == "PUT" {
		res.WriteHeader(http.StatusOK)

		db := DB.Connect()
		defer db.Close()

		body, err := io.ReadAll(req.Body)

		if err != nil {
			middleware.SendError(err, res)
		}
		var categoriesMap map[string]interface{}
		var Category models.Category

		if err := json.Unmarshal(body, &categoriesMap); err != nil {
			middleware.SendError(err, res)
		}

		var token string = fmt.Sprintf("%v", categoriesMap["auth-token"])

		admin := middleware.CheckIsAdmin(token, db)

		if !admin {
			err := fmt.Errorf("user is not admin")
			middleware.SendError(err, res)
		}

		Category.Name = fmt.Sprintf("%s", categoriesMap["name"])
		var name string = fmt.Sprintf("%s", categoriesMap["oldname"])

		Categories := make(chan []byte, 1)

		wg := &sync.WaitGroup{}

		wg.Add(1)

		go models.Category.Update(Category, db, Categories, wg, name)

		wg.Wait()

		close(Categories)

		var category []models.Category

		if err := json.Unmarshal(<-Categories, &category); err != nil {
			middleware.SendError(err, res)
		}

		Response := make(map[string]interface{}, 1)

		Response["Categories"] = category

		if err := json.NewEncoder(res).Encode(Response); err != nil {
			middleware.SendError(err, res)
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
			middleware.SendError(err, res)
		}

		var categoriesMap map[string]interface{}
		var Category models.Category

		if err := json.Unmarshal(body, &categoriesMap); err != nil {
			middleware.SendError(err, res)
		}

		var token string = fmt.Sprintf("%v", categoriesMap["auth-token"])

		admin := middleware.CheckIsAdmin(token, db)

		if !admin {
			err := fmt.Errorf("user is not admin")
			middleware.SendError(err, res)
		}

		Category.Name = fmt.Sprintf("%s", categoriesMap["name"])

		Categories := make(chan []byte, 1)

		wg := &sync.WaitGroup{}

		wg.Add(1)

		go models.Category.Delete(Category, db, Categories, wg)

		wg.Wait()

		close(Categories)

		var category []models.Category

		if err := json.Unmarshal(<-Categories, &category); err != nil {
			middleware.SendError(err, res)
		}

		Response := make(map[string]interface{}, 1)

		Response["Categories"] = category

		if err := json.NewEncoder(res).Encode(Response); err != nil {
			middleware.SendError(err, res)
		}
	}
}

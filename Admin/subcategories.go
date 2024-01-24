package admin

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"

	DB "alwadimarkets/db"
	"alwadimarkets/models"
)

func GetSubCategories(res http.ResponseWriter, req *http.Request, params map[string]string) {
	db := DB.Connect()
	res.WriteHeader(http.StatusOK)

	defer db.Close()

	SubCategories := make(chan []byte, 1)

	wg := &sync.WaitGroup{}

	wg.Add(1)

	go models.SubCategory.Get(models.SubCategory{}, db, SubCategories, wg)

	wg.Wait()

	close(SubCategories)

	var subcategory []models.SubCategory

	err := json.Unmarshal(<-SubCategories, &subcategory)

	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}

	Response := make(map[string]interface{}, 1)

	Response["SubCategories"] = subcategory

	if err := json.NewEncoder(res).Encode(Response); err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}
}

func AddSubCategories(res http.ResponseWriter, req *http.Request, params map[string]string) {
	db := DB.Connect()
	res.WriteHeader(http.StatusOK)

	defer db.Close()

	body, err := io.ReadAll(req.Body)

	if err != nil {
		panic(err.Error())
	}
	var subcategoryMap map[string]interface{}

	var Subcategory models.SubCategory

	Error := json.Unmarshal(body, &subcategoryMap)

	if Error != nil {
		http.Error(res, Error.Error(), http.StatusInternalServerError)
	}

	// body data
	Subcategory.Name = fmt.Sprintf("%s", subcategoryMap["name"])
	Subcategory.Category = fmt.Sprintf("%s", subcategoryMap["category"])
	Subcategory.Image = fmt.Sprintf("%s", subcategoryMap["image"])
	fmt.Printf(Subcategory.Category)
	SubCategory := make(chan []byte, 1)

	wg := &sync.WaitGroup{}

	wg.Add(1)
	go models.SubCategory.Add(Subcategory, db, SubCategory, wg)
	wg.Wait()

	close(SubCategory)
	var subcategory map[string]interface{}

	errors := json.Unmarshal(<-SubCategory, &subcategory)

	if errors != nil {
		http.Error(res, errors.Error(), http.StatusInternalServerError)
	}

	if err := json.NewEncoder(res).Encode(subcategory); err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}
}

func DeleteSubCategories(res http.ResponseWriter, req *http.Request, params map[string]string) {
	if req.Method == "DELETE" {
		db := DB.Connect()
		res.WriteHeader(http.StatusOK)

		defer db.Close()

		body, err := io.ReadAll(req.Body)

		if err != nil {
			panic(err.Error())
		}
		var subcategoryMap map[string]interface{}

		var Subcategory models.SubCategory

		Error := json.Unmarshal(body, &subcategoryMap)

		if Error != nil {
			http.Error(res, Error.Error(), http.StatusInternalServerError)
		}

		// body data
		Subcategory.Name = fmt.Sprintf("%s", subcategoryMap["name"])

		SubCategory := make(chan []byte, 1)

		wg := &sync.WaitGroup{}

		wg.Add(1)
		go models.SubCategory.Delete(Subcategory, db, SubCategory, wg)
		wg.Wait()

		close(SubCategory)
		var subcategory map[string]interface{}

		errors := json.Unmarshal(<-SubCategory, &subcategory)

		if errors != nil {
			http.Error(res, errors.Error(), http.StatusInternalServerError)
		}

		if err := json.NewEncoder(res).Encode(subcategory); err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
		}
	}
}

func UpdateSubCategories(res http.ResponseWriter, req *http.Request, params map[string]string) {
	if req.Method == "PUT" {
		db := DB.Connect()
		res.WriteHeader(http.StatusOK)

		defer db.Close()

		body, err := io.ReadAll(req.Body)

		if err != nil {
			panic(err.Error())
		}
		var subcategoryMap map[string]interface{}

		var Subcategory models.SubCategory

		Error := json.Unmarshal(body, &subcategoryMap)

		if Error != nil {
			http.Error(res, Error.Error(), http.StatusInternalServerError)
		}

		// body data
		Subcategory.Name = fmt.Sprintf("%s", subcategoryMap["name"])
		Subcategory.Category = fmt.Sprintf("%s", subcategoryMap["category"])
		var name string = fmt.Sprintf("%s", subcategoryMap["oldname"])
		SubCategory := make(chan []byte, 1)

		wg := &sync.WaitGroup{}

		wg.Add(1)
		go models.SubCategory.Update(Subcategory, db, SubCategory, wg, name)
		wg.Wait()

		close(SubCategory)
		var subcategory map[string]interface{}

		errors := json.Unmarshal(<-SubCategory, &subcategory)

		if errors != nil {
			http.Error(res, errors.Error(), http.StatusInternalServerError)
		}

		if err := json.NewEncoder(res).Encode(subcategory); err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
		}
	}
}

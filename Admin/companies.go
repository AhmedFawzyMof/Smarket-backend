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

func GetCompanies(res http.ResponseWriter, req *http.Request, params map[string]string) {
	db := DB.Connect()
	res.WriteHeader(http.StatusOK)

	defer db.Close()

	Companies := make(chan []byte, 1)

	wg := &sync.WaitGroup{}

	wg.Add(1)

	go tables.Company.Get(tables.Company{}, db, Companies, wg)

	wg.Wait()

	close(Companies)

	var company []tables.Company

	err := json.Unmarshal(<-Companies, &company)

	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}

	Response := make(map[string]interface{}, 1)

	Response["Companies"] = company

	if err := json.NewEncoder(res).Encode(Response); err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}
}

func AddCompanies(res http.ResponseWriter, req *http.Request, params map[string]string) {
	db := DB.Connect()
	res.WriteHeader(http.StatusOK)

	defer db.Close()

	body, err := io.ReadAll(req.Body)

	if err != nil {
		panic(err.Error())
	}
	var companyMap map[string]interface{}

	var Companies tables.Company

	Error := json.Unmarshal(body, &companyMap)

	if Error != nil {
		http.Error(res, Error.Error(), http.StatusInternalServerError)
	}

	// body data
	Companies.Name = fmt.Sprintf("%s", companyMap["name"])

	Company := make(chan []byte, 1)

	wg := &sync.WaitGroup{}

	wg.Add(1)
	go tables.Company.Add(Companies, db, Company, wg)
	wg.Wait()

	close(Company)
	var company map[string]interface{}

	errors := json.Unmarshal(<-Company, &company)

	if errors != nil {
		http.Error(res, errors.Error(), http.StatusInternalServerError)
	}

	if err := json.NewEncoder(res).Encode(company); err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}
}

func EditCompanies(res http.ResponseWriter, req *http.Request, params map[string]string) {
	if req.Method == "PUT" {
		db := DB.Connect()
		res.WriteHeader(http.StatusOK)

		defer db.Close()

		body, err := io.ReadAll(req.Body)

		if err != nil {
			panic(err.Error())
		}
		var companyMap map[string]interface{}

		var Companies tables.Company

		Error := json.Unmarshal(body, &companyMap)

		if Error != nil {
			http.Error(res, Error.Error(), http.StatusInternalServerError)
		}

		// body data
		Companies.Name = fmt.Sprintf("%s", companyMap["name"])
		var name string = fmt.Sprintf("%s", companyMap["oldname"])

		Company := make(chan []byte, 1)

		wg := &sync.WaitGroup{}

		wg.Add(1)
		go tables.Company.Update(Companies, db, Company, wg, name)
		wg.Wait()

		close(Company)
		var company map[string]interface{}

		errors := json.Unmarshal(<-Company, &company)

		if errors != nil {
			http.Error(res, errors.Error(), http.StatusInternalServerError)
		}

		if err := json.NewEncoder(res).Encode(company); err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
		}
	}
}

func DeleteCompanies(res http.ResponseWriter, req *http.Request, params map[string]string) {
	if req.Method == "DELETE" {
		db := DB.Connect()
		res.WriteHeader(http.StatusOK)

		defer db.Close()

		body, err := io.ReadAll(req.Body)

		if err != nil {
			panic(err.Error())
		}
		var companyMap map[string]interface{}

		var Companies tables.Company

		Error := json.Unmarshal(body, &companyMap)

		if Error != nil {
			http.Error(res, Error.Error(), http.StatusInternalServerError)
		}

		// body data
		Companies.Name = fmt.Sprintf("%s", companyMap["name"])

		Company := make(chan []byte, 1)

		wg := &sync.WaitGroup{}

		wg.Add(1)
		go tables.Company.Delete(Companies, db, Company, wg)
		wg.Wait()

		close(Company)
		var company map[string]interface{}

		errors := json.Unmarshal(<-Company, &company)

		if errors != nil {
			http.Error(res, errors.Error(), http.StatusInternalServerError)
		}

		if err := json.NewEncoder(res).Encode(company); err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
		}
	}
}

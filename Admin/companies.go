package administrator

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"

	DB "alwadimarkets/db"
	"alwadimarkets/middleware"
	"alwadimarkets/models"
)

func GetCompanies(res http.ResponseWriter, req *http.Request, params map[string]string) {
	db := DB.Connect()
	res.WriteHeader(http.StatusOK)

	defer db.Close()

	body, err := io.ReadAll(req.Body)

	if err != nil {
		middleware.SendError(err, res)
	}
	var companyMap map[string]interface{}

	if err := json.Unmarshal(body, &companyMap); err != nil {
		middleware.SendError(err, res)
	}

	var token string = fmt.Sprintf("%s", companyMap["auth-token"])

	admin := middleware.CheckIsAdmin(token, db)

	if !admin {
		err := fmt.Errorf("user is not admin")
		middleware.SendError(err, res)
	}

	Companies := make(chan []byte, 1)

	wg := &sync.WaitGroup{}

	wg.Add(1)

	go models.Company.Get(models.Company{}, db, Companies, wg)

	wg.Wait()

	close(Companies)

	var company []models.Company

	if err := json.Unmarshal(<-Companies, &company); err != nil {
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
		middleware.SendError(err, res)
	}
	var companyMap map[string]interface{}

	var Companies models.Company

	if err := json.Unmarshal(body, &companyMap); err != nil {
		middleware.SendError(err, res)
	}

	var token string = fmt.Sprintf("%s", companyMap["auth-token"])

	admin := middleware.CheckIsAdmin(token, db)

	if !admin {
		err := fmt.Errorf("user is not admin")
		middleware.SendError(err, res)
	}

	// body data
	Companies.Name = fmt.Sprintf("%s", companyMap["name"])

	Company := make(chan []byte, 1)

	wg := &sync.WaitGroup{}

	wg.Add(1)
	go models.Company.Add(Companies, db, Company, wg)
	wg.Wait()

	close(Company)
	var company map[string]interface{}

	if err := json.Unmarshal(<-Company, &company); err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
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
			middleware.SendError(err, res)
		}
		var companyMap map[string]interface{}

		var Companies models.Company

		if err := json.Unmarshal(body, &companyMap); err != nil {
			middleware.SendError(err, res)
		}

		var token string = fmt.Sprintf("%s", companyMap["auth-token"])

		admin := middleware.CheckIsAdmin(token, db)

		if !admin {
			err := fmt.Errorf("user is not admin")
			middleware.SendError(err, res)
		}

		// body data
		Companies.Name = fmt.Sprintf("%s", companyMap["name"])
		var name string = fmt.Sprintf("%s", companyMap["oldname"])

		Company := make(chan []byte, 1)

		wg := &sync.WaitGroup{}

		wg.Add(1)
		go models.Company.Update(Companies, db, Company, wg, name)
		wg.Wait()

		close(Company)
		var company map[string]interface{}

		if err := json.Unmarshal(<-Company, &company); err != nil {
			middleware.SendError(err, res)

		}

		if err := json.NewEncoder(res).Encode(company); err != nil {
			middleware.SendError(err, res)
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
			middleware.SendError(err, res)
		}

		var companyMap map[string]interface{}

		var Companies models.Company

		if err := json.Unmarshal(body, &companyMap); err != nil {
			middleware.SendError(err, res)
		}

		var token string = fmt.Sprintf("%s", companyMap["auth-token"])

		admin := middleware.CheckIsAdmin(token, db)

		if !admin {
			err := fmt.Errorf("user is not admin")
			middleware.SendError(err, res)
		}
		// body data
		Companies.Name = fmt.Sprintf("%s", companyMap["name"])

		Company := make(chan []byte, 1)

		wg := &sync.WaitGroup{}

		wg.Add(1)
		go models.Company.Delete(Companies, db, Company, wg)
		wg.Wait()

		close(Company)
		var company map[string]interface{}

		errors := json.Unmarshal(<-Company, &company)

		if errors != nil {
			middleware.SendError(err, res)

		}

		if err := json.NewEncoder(res).Encode(company); err != nil {
			middleware.SendError(err, res)
		}
	}
}

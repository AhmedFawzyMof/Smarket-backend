package admin

import (
	"Smarket/controller"
	DB "Smarket/db"
	"encoding/json"
	"io"
	"net/http"
)

type Table struct {
	Name string
}

func (t Table) EditTable(res http.ResponseWriter, req *http.Request) {
	db := DB.Connect()

	defer db.Close()
	res.Header().Set("Access-Control-Allow-Origin", "*")
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	body, err := io.ReadAll(req.Body)

	if err != nil {
		panic(err.Error())
	}
	dataForm := make(map[string]string)
	TheToken := json.Unmarshal(body, &dataForm)
	if TheToken != nil {
		panic(TheToken.Error())
	}

	controller.Edit(db, t.Name, dataForm)
}

func (t Table) DeleteTable(res http.ResponseWriter, req *http.Request) {
	db := DB.Connect()

	defer db.Close()
	res.Header().Set("Access-Control-Allow-Origin", "*")
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	body, err := io.ReadAll(req.Body)

	if err != nil {
		panic(err.Error())
	}
	dataForm := make(map[string]interface{})
	TheToken := json.Unmarshal(body, &dataForm)
	if TheToken != nil {
		panic(TheToken.Error())
	}

	var id any
	data, ok := dataForm["id"]
	if !ok {
		id = dataForm["name"]
	} else {
		id = data
	}

	controller.Delete(db, t.Name, id)
}

func (t Table) CreateTable(res http.ResponseWriter, req *http.Request) {
	db := DB.Connect()

	defer db.Close()
	res.Header().Set("Access-Control-Allow-Origin", "*")
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	body, err := io.ReadAll(req.Body)

	if err != nil {
		panic(err.Error())
	}
	dataForm := make(map[string]interface{})
	TheToken := json.Unmarshal(body, &dataForm)
	if TheToken != nil {
		panic(TheToken.Error())
	}

	controller.Create(db, t.Name, dataForm)
}

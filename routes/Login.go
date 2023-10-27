package routes

import (
	"Smarket/controller"
	DB "Smarket/db"
	"encoding/json"
	"io"
	"net/http"
)

func Login(res http.ResponseWriter, req *http.Request) {
	db := DB.Connect()

	defer db.Close()
	res.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)

	body, err := io.ReadAll(req.Body)

	if err != nil {
		panic(err)
	}

	dataForm := make(map[string]interface{})
	mapData := json.Unmarshal(body, &dataForm)
	if mapData != nil {
		panic(mapData.Error())
	}

	UserRes := controller.GetUser(db, dataForm)

	json.NewEncoder(res).Encode(UserRes)
}

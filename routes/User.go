package routes

import (
	"Smarket/controller"
	DB "Smarket/db"
	"encoding/json"
	"io"
	"net/http"
)

func GetUserData(res http.ResponseWriter, req *http.Request) {
	db := DB.Connect()

	defer db.Close()
	res.Header().Set("Access-Control-Allow-Origin", "*")
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)

	body, err := io.ReadAll(req.Body)

	if err != nil {
		panic(err)
	}

	dataForm := make(map[string]string)
	mapData := json.Unmarshal(body, &dataForm)
	if mapData != nil {
		panic(mapData.Error())
	}

	token := dataForm["authToken"]

	UserInfo := controller.GetUserInfo(db, token)

	json.NewEncoder(res).Encode(UserInfo)
}

func EditProfile(res http.ResponseWriter, req *http.Request) {
	db := DB.Connect()

	defer db.Close()
	res.Header().Set("Access-Control-Allow-Origin", "*")
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
	UserRes := controller.EditUserInfo(db, dataForm)

	json.NewEncoder(res).Encode(UserRes)
}

func Register(res http.ResponseWriter, req *http.Request) {
	db := DB.Connect()

	defer db.Close()
	res.Header().Set("Access-Control-Allow-Origin", "*")

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

	UserRes := controller.AddUser(db, dataForm)

	json.NewEncoder(res).Encode(UserRes)
}

func Login(res http.ResponseWriter, req *http.Request) {
	db := DB.Connect()

	defer db.Close()
	res.Header().Set("Access-Control-Allow-Origin", "*")
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

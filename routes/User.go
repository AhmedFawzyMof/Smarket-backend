package routes

import (
	"Smarket/controller"
	DB "Smarket/db"
	"encoding/json"
	"io"
	"net/http"
)

type UserEmail struct {
	Token string
}

func (u UserEmail) GetUserData(res http.ResponseWriter, req *http.Request) {
	db := DB.Connect()

	defer db.Close()
	res.Header().Set("Access-Control-Allow-Origin", "https://beautiful-ruby-sheath-dress.cyclic.app/")
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)

	UserInfo := controller.GetUserInfo(db, u.Token)

	json.NewEncoder(res).Encode(UserInfo)
}

func EditProfile(res http.ResponseWriter, req *http.Request) {
	db := DB.Connect()

	defer db.Close()
	res.Header().Set("Access-Control-Allow-Origin", "https://beautiful-ruby-sheath-dress.cyclic.app/")
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
	res.Header().Set("Access-Control-Allow-Origin", "https://beautiful-ruby-sheath-dress.cyclic.app/")

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
	res.Header().Set("Access-Control-Allow-Origin", "https://wild-pear-millipede.cyclic.app")
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

package routes

import (
	"alwadi/controller"
	DB "alwadi/db"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func Fav(res http.ResponseWriter, req *http.Request) {
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

	FavRes := controller.AddToFav(db, dataForm)

	json.NewEncoder(res).Encode(FavRes)

}

func GetFav(res http.ResponseWriter, req *http.Request) {
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
	var token string = fmt.Sprintf("%v", dataForm["authToken"])
	FavRes := controller.GetUserFav(db, token)

	Res := map[string]interface{}{
		"products": FavRes,
	}
	
	fmt.Println(dataForm)

	json.NewEncoder(res).Encode(Res)
}

func DelFav(res http.ResponseWriter, req *http.Request) {
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

	FavRes := controller.DelUserFav(db, dataForm)

	json.NewEncoder(res).Encode(FavRes)
}

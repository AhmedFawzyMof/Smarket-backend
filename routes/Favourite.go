package routes

import (
	"Smarket/controller"
	DB "Smarket/db"
	"encoding/json"
	"io"
	"net/http"
)

type Favourite struct {
	Token string
}

func Fav(res http.ResponseWriter, req *http.Request) {
	db := DB.Connect()

	defer db.Close()
	res.Header().Set("Access-Control-Allow-Origin", "https://wild-pear-millipede.cyclic.app/")
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

func (f Favourite) GetFav(res http.ResponseWriter, req *http.Request) {
	db := DB.Connect()

	defer db.Close()
	res.Header().Set("Access-Control-Allow-Origin", "https://wild-pear-millipede.cyclic.app/")
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)

	FavRes := controller.GetUserFav(db, f.Token)

	Res := map[string]interface{}{
		"products": FavRes,
	}

	json.NewEncoder(res).Encode(Res)
}

func DelFav(res http.ResponseWriter, req *http.Request) {
	db := DB.Connect()

	defer db.Close()
	res.Header().Set("Access-Control-Allow-Origin", "https://wild-pear-millipede.cyclic.app/")
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

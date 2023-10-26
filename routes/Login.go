package routes

import (
	"Smarket/controller"
	DB "Smarket/db"
	"encoding/json"
	"io"
	"net/http"
	"strings"
)

func Login(res http.ResponseWriter, req *http.Request) {
	db := DB.Connect()

	defer db.Close()
	res.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)

	body, err := io.ReadAll(req.Body)
	splitedBody := strings.Split(string(body), "&")
	if err != nil {
		panic(err)
	}

	theUser := make(map[string]string)

	for _, data := range splitedBody {
		key := strings.Split(data, "=")[0]
		value := strings.Split(data, "=")[1]
		theUser[key] = value
	}

	UserRes := controller.GetUser(db, theUser)

	json.NewEncoder(res).Encode(UserRes)
}

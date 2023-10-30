package routes

import (
	"Smarket/controller"
	DB "Smarket/db"
	"encoding/json"
	"fmt"
	"net/http"
)

type ProductId struct {
	Id int
}

func (p ProductId) GetById(res http.ResponseWriter, req *http.Request) {
	db := DB.Connect()

	defer db.Close()
	res.Header().Set("Access-Control-Allow-Origin", "https://wild-pear-millipede.cyclic.app")
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)

	
	Product := controller.ProductGetId(db, p.Id)

	fmt.Println(Product, p.Id)

	var data = map[string]interface{}{
		"product": Product,
	}

	json.NewEncoder(res).Encode(data)
}

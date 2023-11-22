package admin

import (
	"alwadi/controller"
	DB "alwadi/db"
	"encoding/json"
	"net/http"
	"sync"
)

func GetUsers(res http.ResponseWriter, req *http.Request) {
	db := DB.Connect()

	defer db.Close()
	res.Header().Set("Access-Control-Allow-Origin", "*")
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)

	users := make(chan []controller.Users, 1)
	wg := &sync.WaitGroup{}

	wg.Add(1)

	go controller.GetallUsers(db, users, wg)

	wg.Wait()

	close(users)

	var Users = map[string]interface{}{
		"Users": <-users,
	}

	json.NewEncoder(res).Encode(Users)
}

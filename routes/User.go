package routes

import (
	DB "alwadi_markets/db"
	"alwadi_markets/middleware"
	"alwadi_markets/tables"
	"encoding/json"
	"io"
	"net/http"
	"sync"
)

func Register(res http.ResponseWriter, req *http.Request, params map[string]string) {
	if req.Method == "POST" {
		db := DB.Connect()
		res.WriteHeader(http.StatusOK)

		defer db.Close()

		body, err := io.ReadAll(req.Body)

		if err != nil {
			panic(err.Error())
		}
		var userMap map[string]string

		var User tables.Users

		Error := json.Unmarshal(body, &userMap)

		if Error != nil {
			http.Error(res, Error.Error(), http.StatusInternalServerError)
		}
		User.Email = userMap["email"]
		User.Password = userMap["password"]
		User.Password2 = userMap["password2"]
		User.Phone = userMap["phone"]
		User.Spare_phone = userMap["spare_phone"]
		User.Username = userMap["username"]
		User.Terms = userMap["terms"]

		UserChan := make(chan []byte, 1)

		wg := &sync.WaitGroup{}

		wg.Add(1)
		go tables.Users.Create(User, db, UserChan, wg)
		wg.Wait()

		close(UserChan)

		var UserResponse map[string]interface{}

		errors := json.Unmarshal(<-UserChan, &UserResponse)

		if errors != nil {
			http.Error(res, errors.Error(), http.StatusInternalServerError)
		}

		if err := json.NewEncoder(res).Encode(UserResponse); err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
		}
	}
}

func Login(res http.ResponseWriter, req *http.Request, params map[string]string) {
	if req.Method == "POST" {
		db := DB.Connect()
		res.WriteHeader(http.StatusOK)

		defer db.Close()

		body, err := io.ReadAll(req.Body)

		if err != nil {
			panic(err.Error())
		}
		var userMap map[string]string

		var User tables.Users

		Error := json.Unmarshal(body, &userMap)

		if Error != nil {
			http.Error(res, Error.Error(), http.StatusInternalServerError)
		}
		User.Email = userMap["email"]
		User.Password = userMap["password"]

		UserChan := make(chan []byte, 1)

		wg := &sync.WaitGroup{}

		wg.Add(1)
		go tables.Users.GetByEmail(User, db, UserChan, wg)
		wg.Wait()

		close(UserChan)

		var UserResponse map[string]interface{}

		errors := json.Unmarshal(<-UserChan, &UserResponse)

		if errors != nil {
			http.Error(res, errors.Error(), http.StatusInternalServerError)
		}

		if err := json.NewEncoder(res).Encode(UserResponse); err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
		}
	}
}

func Profile(res http.ResponseWriter, req *http.Request, params map[string]string) {
	if req.Method == "POST" {
		db := DB.Connect()
		res.WriteHeader(http.StatusOK)

		defer db.Close()

		body, err := io.ReadAll(req.Body)

		if err != nil {
			panic(err.Error())
		}
		var userMap map[string]string

		var User tables.Users

		Error := json.Unmarshal(body, &userMap)

		if Error != nil {
			http.Error(res, Error.Error(), http.StatusInternalServerError)
		}

		id, e := middleware.VerifyToken(userMap["token"])
		if e != nil {
			http.Error(res, e.Error(), http.StatusInternalServerError)
		}

		User.Id = id
		UserChan := make(chan []byte, 1)

		wg := &sync.WaitGroup{}

		wg.Add(1)
		go tables.Users.GetById(User, db, UserChan, wg)
		wg.Wait()

		close(UserChan)

		UserResponse := make(map[string]interface{})

		var Users tables.Users

		errors := json.Unmarshal(<-UserChan, &Users)

		if errors != nil {
			http.Error(res, errors.Error(), http.StatusInternalServerError)
		}

		UserResponse["User"] = Users

		if err := json.NewEncoder(res).Encode(UserResponse); err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
		}
	}
}

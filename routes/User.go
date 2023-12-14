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
			middleware.SendError(err, res)
		}
		var userMap map[string]string

		var User tables.Users

		if err := json.Unmarshal(body, &userMap); err != nil {
			middleware.SendError(err, res)
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

		if err := json.Unmarshal(<-UserChan, &UserResponse); err != nil {
			middleware.SendError(err, res)
		}

		if err := json.NewEncoder(res).Encode(UserResponse); err != nil {
			middleware.SendError(err, res)
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
			middleware.SendError(err, res)
		}
		var userMap map[string]string

		var User tables.Users

		if err := json.Unmarshal(body, &userMap); err != nil {
			middleware.SendError(err, res)
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

		if err := json.Unmarshal(<-UserChan, &UserResponse); err != nil {
			middleware.SendError(err, res)
		}

		if err := json.NewEncoder(res).Encode(UserResponse); err != nil {
			middleware.SendError(err, res)
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
			middleware.SendError(err, res)
		}
		var userMap map[string]string

		var User tables.Users

		if err := json.Unmarshal(body, &userMap); err != nil {
			middleware.SendError(err, res)
		}

		id, err := middleware.VerifyToken(userMap["token"])
		if err != nil {
			middleware.SendError(err, res)
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

		if err := json.Unmarshal(<-UserChan, &Users); err != nil {
			middleware.SendError(err, res)
		}

		UserResponse["User"] = Users

		if err := json.NewEncoder(res).Encode(UserResponse); err != nil {
			middleware.SendError(err, res)
		}
	}
}

package admin

import (
	DB "alwadi_markets/db"
	"alwadi_markets/middleware"
	"alwadi_markets/models"
	"crypto/sha256"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
)

func GetUsers(res http.ResponseWriter, req *http.Request, params map[string]string) {
	db := DB.Connect()
	res.WriteHeader(http.StatusOK)

	defer db.Close()

	Users := make(chan []byte, 1)

	wg := &sync.WaitGroup{}

	wg.Add(1)

	go models.Users.Get(models.Users{}, db, Users, wg)

	wg.Wait()

	close(Users)

	var users []models.Users

	err := json.Unmarshal(<-Users, &users)

	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}

	Response := make(map[string]interface{}, 1)

	Response["Users"] = users

	if err := json.NewEncoder(res).Encode(Response); err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}
}

func AdminLogin(res http.ResponseWriter, req *http.Request, params map[string]string) {
	if req.Method == "POST" {
		res.WriteHeader(http.StatusOK)

		db := DB.Connect()

		defer db.Close()

		body, err := io.ReadAll(req.Body)

		if err != nil {
			panic(err.Error())
		}
		var user map[string]string
		var User models.Users

		Error := json.Unmarshal(body, &user)

		if Error != nil {
			http.Error(res, Error.Error(), http.StatusInternalServerError)
		}

		User.Email = user["email"]
		User.Password = user["password"]

		notOk, token := checkIsAdmin(User, db)

		if !notOk {
			Response := map[string]interface{}{
				"Error": true,
			}

			if err := json.NewEncoder(res).Encode(Response); err != nil {
				http.Error(res, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		Response := map[string]interface{}{
			"Error": false,
			"Token": token,
		}

		if err := json.NewEncoder(res).Encode(Response); err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
		}
	}
}

func checkIsAdmin(u models.Users, db *sql.DB) (bool, string) {
	var sampleSecretKey = []byte("Ahmedfawzi made this website")

	sha := sha256.New()
	sha.Write([]byte(u.Password))
	pass := sha.Sum(nil)
	Password := fmt.Sprintf("%x", pass)

	FindEmail := db.QueryRow("SELECT * FROM Users WHERE (email, password) = (?, ?)", u.Email, Password)

	var User models.Users

	err := FindEmail.Scan(&User.Id, &User.Username, &User.Email, &User.Password, &User.Phone, &User.Spare_phone, &User.Role)

	token, Err := middleware.GenerateJWT(User.Id, sampleSecretKey)

	if Err != nil {
		return false, ""
	}

	if err != nil {
		return false, ""
	}

	if User.Role == "Admin" {
		return true, token
	}
	return false, ""
}

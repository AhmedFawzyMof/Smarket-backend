package models

import (
	"crypto/sha256"
	"database/sql"
	"encoding/json"
	"fmt"
	"sync"

	"alwadi_markets/middleware"

	"github.com/google/uuid"
)

var sampleSecretKey = []byte("Ahmedfawzi made this website")

type Users struct {
	Id          string
	Username    string
	Email       string
	Password    string
	Password2   string
	Phone       string
	Spare_phone string
	Role        string
	Terms       string
}

func (u Users) Create(db *sql.DB, response chan []byte, wg *sync.WaitGroup) {
	defer wg.Done()

	Response := make(map[string]interface{})
	id := uuid.New()
	u.Id = fmt.Sprintf("%v", id)

	if u.Password != u.Password2 {
		Response["Error"] = true
		Response["Message"] = "كلمة المرور خطأ"

		middleware.SendResponse(response, Response)
		return
	}

	if u.Terms != "yes" {
		Response["Error"] = true
		Response["Message"] = "لم يتم الموافقة علي شروط والأحكام"

		middleware.SendResponse(response, Response)
		return
	}

	FindEmail := db.QueryRow("SELECT * FROM Users WHERE email=?", u.Email)

	var User Users

	err := FindEmail.Scan(&User.Id, &User.Username, &User.Email, &User.Password, &User.Phone, &User.Spare_phone, &User.Role)
	sha := sha256.New()
	sha.Write([]byte(u.Password))
	pass := sha.Sum(nil)
	Password := fmt.Sprintf("%x", pass)

	if err != nil {
		_, err := db.Exec("INSERT INTO Users (id, username, email, password, phone, spare_phone, role) VALUES (?, ?, ?, ?, ?, ?, ?)", u.Id, u.Username, u.Email, Password, u.Phone, u.Spare_phone, "Normal")
		if err != nil {
			Response["Error"] = true
			Response["Message"] = "حدث خطأ ما يرجى إعادة المحاولة"

			middleware.SendResponse(response, Response)
			return
		}

		token, err := middleware.GenerateJWT(u.Id, sampleSecretKey)
		if err != nil {
			Response["Error"] = true
			Response["Message"] = "حدث خطأ ما يرجى إعادة المحاولة"

			middleware.SendResponse(response, Response)
			return
		}
		Response["Error"] = false
		Response["token"] = token

		middleware.SendResponse(response, Response)
		return
	}

	Response["Error"] = true
	Response["Message"] = "البريد الإلكتروني موجود بالفعل"

	middleware.SendResponse(response, Response)
}

func (u Users) GetByEmail(db *sql.DB, response chan []byte, wg *sync.WaitGroup) {
	defer wg.Done()

	Response := make(map[string]interface{})
	sha := sha256.New()
	sha.Write([]byte(u.Password))
	pass := sha.Sum(nil)
	Password := fmt.Sprintf("%x", pass)

	FindEmail := db.QueryRow("SELECT * FROM Users WHERE (email, password) = (?, ?)", u.Email, Password)
	var User Users

	err := FindEmail.Scan(&User.Id, &User.Username, &User.Email, &User.Password, &User.Phone, &User.Spare_phone, &User.Role)

	token, Err := middleware.GenerateJWT(User.Id, sampleSecretKey)
	if Err != nil {
		Response["Error"] = true
		Response["Message"] = "حدث خطأ ما يرجى إعادة المحاولة"

		middleware.SendResponse(response, Response)
		return
	}

	if err != nil {
		Response["Error"] = true
		Response["Message"] = "لا يمكن تسجيل الدخول ببيانات المقدمة"

		middleware.SendResponse(response, Response)
		return
	}

	Response["Error"] = false
	Response["token"] = token

	middleware.SendResponse(response, Response)
}

func (u Users) GetById(db *sql.DB, response chan []byte, wg *sync.WaitGroup) {
	defer wg.Done()

	Response := make(map[string]interface{})

	FindEmail := db.QueryRow("SELECT * FROM Users WHERE id = ?", u.Id)
	var User Users

	err := FindEmail.Scan(&User.Id, &User.Username, &User.Email, &User.Password, &User.Phone, &User.Spare_phone, &User.Role)

	if err != nil {
		Response["Error"] = true
		Response["Message"] = "لا يمكن تسجيل الدخول ببيانات المقدمة"

		middleware.SendResponse(response, Response)
		return
	}

	res, err := json.Marshal(User)

	if err != nil {
		fmt.Println(err)
	}

	response <- res
}

func (u Users) Get(db *sql.DB, response chan []byte, wg *sync.WaitGroup) {
	defer wg.Done()
	var Response []Users
	users, err := db.Query("SELECT * FROM Users")
	if err != nil {
		fmt.Println(err.Error())
	}

	defer users.Close()

	for users.Next() {
		var user Users

		if err := users.Scan(&user.Id, &user.Username, &user.Email, &user.Password, &user.Phone, &user.Spare_phone, &user.Role); err != nil {
			panic(err.Error())
		}

		Response = append(Response, user)
	}
	res, err := json.Marshal(Response)
	if err != nil {
		fmt.Println(err.Error())
	}
	response <- res
}

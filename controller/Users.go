package controller

import (
	"crypto/sha256"
	"database/sql"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

var sampleSecretKey = []byte("Ahmedfawzi made this website")

type user struct {
	id          string
	username    string
	email       string
	password    string
	phone       string
	spare_phone string
	role        string
}

func generateJWT(email string, key []byte) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": email,
		"exp":   time.Now().Add(time.Hour * 8760).Unix(),
	})
	tokenString, err := token.SignedString(key)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func AddUser(db *sql.DB, userData map[string]string) map[string]interface{} {
	var username string = userData["username"]
	var email string = userData["email"]
	newEmail := strings.Replace(email, "%40", "@", 1)
	var password string = userData["password"]
	var password2 string = userData["password2"]
	var phone string = userData["phone"]
	var spare_phone string = userData["spare_phone"]
	var terms string = userData["terms"]
	id := uuid.New()
	if password != password2 {
		ErrorRes := map[string]interface{}{
			"Error":   true,
			"Message": "كلمة المرور خطأ",
		}

		return ErrorRes
	}

	if terms != "yes" {
		ErrorRes := map[string]interface{}{
			"Error":   true,
			"Message": "لم يتم الموافقة علي شروط والأحكام",
		}

		return ErrorRes
	}

	FindEmail := db.QueryRow("SELECT * FROM Users WHERE email=?", email)

	var User user
	err := FindEmail.Scan(&User.id, &User.username, &User.email, &User.password, &User.password, &User.phone, &User.spare_phone, &User.role)
	sha := sha256.New()
	sha.Write([]byte(password))
	pass := sha.Sum(nil)
	if err != nil {
		_, err := db.Exec("INSERT INTO Users (id, username, email, password, phone, spare_phone, role) VALUES (?, ?, ?, ?, ?, ?, ?)", id, username, newEmail, string(pass), phone, spare_phone, "Normal")
		if err != nil {
			userRes := map[string]interface{}{
				"Error":   true,
				"Message": "البريد الإلكتروني مستخدم بالفعل",
			}
			return userRes
		}
	}
	token, err := generateJWT(email, sampleSecretKey)
	if err != nil {
		userRes := map[string]interface{}{
			"Error":   true,
			"Message": "حدث خطأ ما يرجى إعادة المحاولة",
		}
		return userRes
	}
	userRes := map[string]interface{}{
		"Error": false,
		"token": token,
	}

	return userRes
}

func GetUser(db *sql.DB, userData map[string]string) map[string]interface{} {
	var email string = userData["email"]
	var password string = userData["password"]

	sha := sha256.New()
	sha.Write([]byte(password))
	pass := sha.Sum(nil)

	FindEmail := db.QueryRow("SELECT * FROM Users WHERE (email, password) = (?, ?)", email, pass)
	var User user
	err := FindEmail.Scan(&User.id, &User.username, &User.email, &User.password, &User.password, &User.phone, &User.spare_phone, &User.role)

	token, Err := generateJWT(email, sampleSecretKey)
	if Err != nil {
		userRes := map[string]interface{}{
			"Error":   true,
			"Message": "حدث خطأ ما يرجى إعادة المحاولة",
		}
		return userRes
	}

	if err == nil {
		theUser := map[string]interface{}{
			"Error": false,
			"token": token,
		}
		return theUser
	}

	theUser := map[string]interface{}{
		"Error":   true,
		"Message": "لا يمكن تسجيل الدخول ببيانات المقدمة",
	}

	return theUser
}

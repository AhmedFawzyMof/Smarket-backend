package controller

import (
	"crypto/sha256"
	"database/sql"
	"encoding/json"
	"fmt"
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

func verifyToken(tokenString string) (jwt.MapClaims, error) {
	var sampleSecretKey = []byte("Ahmedfawzi made this website")

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return sampleSecretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, fmt.Errorf("invalid jwt token")
	}
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

func AddUser(db *sql.DB, userData map[string]interface{}) map[string]interface{} {
	var username string = fmt.Sprintf("%v", userData["username"])
	var email string = fmt.Sprintf("%v", userData["email"])
	var password string = fmt.Sprintf("%v", userData["password"])
	var password2 string = fmt.Sprintf("%v", userData["password2"])
	var phone string = fmt.Sprintf("%v", userData["phone"])
	var spare_phone string = fmt.Sprintf("%v", userData["spare_phone"])
	var terms string = fmt.Sprintf("%v", userData["terms"])
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
	err := FindEmail.Scan(&User.id, &User.username, &User.email, &User.password, &User.phone, &User.spare_phone, &User.role)
	sha := sha256.New()
	sha.Write([]byte(password))
	pass := sha.Sum(nil)
	Password := fmt.Sprintf("%x", pass)
	if err != nil {
		_, err := db.Exec("INSERT INTO Users (id, username, email, password, phone, spare_phone, role) VALUES (?, ?, ?, ?, ?, ?, ?)", id, username, email, Password, phone, spare_phone, "Normal")
		if err != nil {
			panic(err.Error())
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

func GetUser(db *sql.DB, userData map[string]interface{}) map[string]interface{} {
	var email string = fmt.Sprintf("%v", userData["email"])
	var password string = fmt.Sprintf("%v", userData["password"])
	sha := sha256.New()
	sha.Write([]byte(password))
	pass := sha.Sum(nil)
	Password := fmt.Sprintf("%x", pass)

	FindEmail := db.QueryRow("SELECT * FROM Users WHERE (email, password) = (?, ?)", email, Password)
	var User user

	err := FindEmail.Scan(&User.id, &User.username, &User.email, &User.password, &User.phone, &User.spare_phone, &User.role)

	token, Err := generateJWT(email, sampleSecretKey)
	if Err != nil {
		userRes := map[string]interface{}{
			"Error":   true,
			"Message": "حدث خطأ ما يرجى إعادة المحاولة",
		}
		return userRes
	}

	if err != nil {
		theUser := map[string]interface{}{
			"Error":   true,
			"Message": "لا يمكن تسجيل الدخول ببيانات المقدمة",
		}
		return theUser
	}

	theUser := map[string]interface{}{
		"Error": false,
		"token": token,
	}
	return theUser
}

func GetUserInfo(db *sql.DB, token string) map[string]interface{} {
	claims, err := verifyToken(token)
	if err != nil {
		UserRes := map[string]interface{}{
			"Error": true,
		}
		return UserRes
	}
	var tm time.Time
	switch iat := claims["exp"].(type) {
	case float64:
		tm = time.Unix(int64(iat), 0)
	case json.Number:
		v, _ := iat.Int64()
		tm = time.Unix(v, 0)
	}

	if tm == time.Now() {
		UserRes := map[string]interface{}{
			"Error": true,
		}
		return UserRes
	}

	FindEmail := db.QueryRow("SELECT * FROM Users WHERE email = ?", claims["email"])

	var User user

	Err := FindEmail.Scan(&User.id, &User.username, &User.email, &User.password, &User.phone, &User.spare_phone, &User.role)

	if Err != nil {
		UserRes := map[string]interface{}{
			"Error": true,
		}
		return UserRes
	}
	UserRes := map[string]interface{}{
		"username": User.username,
		"email":    User.email,
	}
	return UserRes
}
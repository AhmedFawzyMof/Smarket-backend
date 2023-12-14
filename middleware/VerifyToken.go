package middleware

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

func GenerateJWT(id string, key []byte) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  id,
		"exp": time.Now().Add(time.Hour * 8760).Unix(),
	})
	tokenString, err := token.SignedString(key)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func VerifyToken(tokenString string) (string, error) {
	var sampleSecretKey = []byte("Ahmedfawzi made this website")

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return sampleSecretKey, nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		id := getIDFromClaims(claims)
		return id, nil
	} else {
		return "", fmt.Errorf("invalid jwt token")
	}
}

func getIDFromClaims(claims jwt.MapClaims) string {
	idValue := claims["id"]

	id := idValue.(string)

	return id
}

type user struct {
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

func CheckIsAdmin(token string, db *sql.DB) bool {
	id, errors := VerifyToken(token)

	if errors != nil {
		return true
	}

	FindEmail := db.QueryRow("SELECT * FROM Users WHERE id = ?", id)

	var User user

	err := FindEmail.Scan(&User.Id, &User.Username, &User.Email, &User.Password, &User.Phone, &User.Spare_phone, &User.Role)

	if err != nil {
		return false
	}

	if User.Role == "Admin" {
		return true
	}
	return false
}

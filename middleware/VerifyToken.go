package middleware

import (
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
		fmt.Println(claims, ok)
		return "", fmt.Errorf("invalid jwt token")
	}
}

func getIDFromClaims(claims jwt.MapClaims) string {
	idValue := claims["id"]

	id := idValue.(string)

	return id
}

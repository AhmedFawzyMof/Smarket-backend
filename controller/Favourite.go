package controller

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

type fav struct {
	id      int
	user    string
	product int
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

func AddToFav(db *sql.DB, favData map[string]interface{}) map[string]interface{} {
	product := favData["product"]
	token := fmt.Sprintf("%v", favData["token"])
	claims, err := verifyToken(token)
	if err != nil {
		FavRes := map[string]interface{}{
			"Error": true,
		}
		return FavRes
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
		FavRes := map[string]interface{}{
			"Error": true,
		}
		return FavRes
	}

	FindEmail := db.QueryRow("SELECT id FROM Users WHERE email=?", claims["email"])

	var User user

	isErr := FindEmail.Scan(&User.id)

	if isErr == nil {

		FindProduct := db.QueryRow("SELECT * FROM Favourite WHERE (user, product) = (?, ?)", User.id, product)
		var Fav fav

		err := FindProduct.Scan(&Fav.id, &Fav.user, &Fav.product)
		if err != nil {
			fmt.Println(err)
			_, err := db.Exec("INSERT INTO Favourite (user, product) VALUES (?, ?)", User.id, product)
			if err != nil {
				userRes := map[string]interface{}{
					"Error": true,
				}
				return userRes
			}
		}
		if err == nil {
			userRes := map[string]interface{}{
				"Error":   false,
				"Message": "المنتج موجود بالفعل في المفضلة",
			}
			return userRes
		}

	}
	FavRes := map[string]interface{}{
		"Error": false,
	}
	return FavRes
}

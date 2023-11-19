package controller

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
)

type fav struct {
	id      int
	user    string
	product int
}

type favProducts map[string]interface{}

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

func GetUserFav(db *sql.DB, token string) []favProducts {
	var Products []int
	var FavProducts []favProducts

	claims, err := verifyToken(token)
	if err != nil {
		FavRes := map[string]interface{}{
			"Error": true,
		}
		FavProducts = append(FavProducts, FavRes)
		return FavProducts
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
		FavProducts = append(FavProducts, FavRes)
		return FavProducts
	}

	FindEmail := db.QueryRow("SELECT id FROM Users WHERE email=?", claims["email"])

	var User user

	isErr := FindEmail.Scan(&User.id)

	if isErr == nil {
		getProducts, err := db.Query("SELECT * FROM Favourite WHERE user = ?", User.id)

		if err != nil {
			panic(err.Error())
		}

		defer getProducts.Close()

		for getProducts.Next() {
			var Fav fav

			if err = getProducts.Scan(&Fav.id, &Fav.user, &Fav.product); err != nil {
				panic(err.Error())
			}

			Products = append(Products, Fav.product)
		}

		ids := ""

		for i, id := range Products {
			if len(Products) > 0 {
				if i >= 0 {
					ids += fmt.Sprintf("%d,", id)
				}
				if i == len(Products)-1 {
					ids += fmt.Sprintf("%d", id)
				}
			}
		}
		var stmt string = fmt.Sprintf("SELECT * FROM Products WHERE id IN (%s)", ids)
		if ids != "" {
			products, err := db.Query(stmt)

			if err != nil {
				fmt.Println(stmt, err)
			}

			defer products.Close()

			for products.Next() {
				var Product product

				if err := products.Scan(&Product.id, &Product.name, &Product.description, &Product.price, &Product.company, &Product.subcategories, &Product.category, &Product.image, &Product.unit, &Product.available, &Product.offer, &Product.inStock, &Product.pricePerUint, &Product.unitNumber); err != nil {
					panic(err.Error())
				}

				TheProduct := map[string]interface{}{
					"id":            Product.id,
					"name":          Product.name,
					"description":   Product.description,
					"price":         Product.price,
					"company":       Product.company,
					"subcategories": Product.subcategories,
					"category":      Product.category,
					"image":         Product.image,
					"unit":          Product.unit,
					"available":     Product.available,
					"offer":         Product.offer,
					"inStock":       Product.inStock,
					"pricePerUint":  Product.pricePerUint,
					"unitNumber":    Product.unitNumber,
				}

				FavProducts = append(FavProducts, TheProduct)
			}

		}

	}
	return FavProducts
}

func DelUserFav(db *sql.DB, formData map[string]interface{}) map[string]interface{} {
	product := formData["product"]
	token := fmt.Sprintf("%v", formData["token"])

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
		_, err := db.Exec("DELETE FROM Favourite WHERE (user, product) = (?, ?)", User.id, product)

		if err != nil {
			userRes := map[string]interface{}{
				"Error": true,
			}
			return userRes
		}
	}

	FavRes := map[string]interface{}{
		"Error": false,
	}
	return FavRes
}

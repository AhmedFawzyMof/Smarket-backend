package controller

import (
	"database/sql"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

func CreateBlob(image interface{}) []byte {
	iamgeBase64 := fmt.Sprintf("%s", image)
	imgBase64 := strings.Split(iamgeBase64, ";base64,")[1]
	imageBlob, err := base64.StdEncoding.DecodeString(imgBase64)

	if err != nil {
		panic(err.Error())
	}

	return imageBlob
}

func Create(db *sql.DB, table string, data map[string]interface{}) {
	switch table {
	case "category":
		name, image := data["name"], data["image"]

		imageBlob := CreateBlob(image)

		_, Err := db.Exec("INSERT INTO `Categories`(`name`, `image`) VALUES (?,?)", name, imageBlob)

		if Err != nil {
			panic(Err.Error())
		}
	case "company":
		name := data["name"]
		_, err := db.Exec("INSERT INTO `Companies`(`name`) VALUES (?)", name)

		if err != nil {
			panic(err.Error())
		}
	case "offer":
		product, image := data["product"], data["image"]
		imageBlob := CreateBlob(image)

		_, err := db.Exec("INSERT INTO `Offer`(product, image) VALUES(?, ?)", product, imageBlob)

		if err != nil {
			panic(err.Error())
		}
	case "product":
		name, description, price, company, subcategories, category, image, unit, available, offer, inStock, pricePerUint, unitNumber := data["name"], data["description"], data["price"], data["company"], data["subcategories"], data["category"], data["image"], data["unit"], data["available"], data["offer"], data["inStock"], data["pricePerUint"], data["unitNumber"]
		imageBlob := CreateBlob(image)

		_, err := db.Exec("INSERT INTO `Products`(name, description, price, company, subcategories, category, image, unit, available, offer, inStock, pricePerUint, unitNumber) VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?)", name, description, price, company, subcategories, category, imageBlob, unit, available, offer, inStock, pricePerUint, unitNumber)

		if err != nil {
			panic(err.Error())
		}
	case "subcategory":
		name, category, image := data["name"], data["category"], data["image"]

		imageBlob := CreateBlob(image)
		_, err := db.Exec("INSERT INTO `SubCategories`(name, category, image) VALUES(?,?,?)", name, category, imageBlob)

		if err != nil {
			panic(err.Error())
		}
	case "users":
		id := uuid.New()
		username, email, password, phone, spare_phone, role := data["username"], data["email"], data["password"], data["phone"], data["spare_phone"], data["role"]
		_, err := db.Exec("INSERT INTO `Users`(id, username, email, password, phone, spare_phone, role) VALUES(?,?,?,?,?,?,?)", id, username, email, password, phone, spare_phone, role)

		if err != nil {
			panic(err.Error())
		}
	}
}

package controller

import (
	"database/sql"
)

func Delete(db *sql.DB, table string, id any) {
	switch table {
	case "product":
		_, err := db.Exec("DELETE FROM `Products` WHERE id=?", id)

		if err != nil {
			panic(err.Error())
		}
	case "company":
		_, err := db.Exec("DELETE FROM `Companies` WHERE name=?", id)

		if err != nil {
			panic(err.Error())
		}
	case "offer":
		_, err := db.Exec("DELETE FROM `Offer` WHERE id=?", id)

		if err != nil {
			panic(err.Error())
		}
	case "category":
		_, err := db.Exec("DELETE FROM `Categories` WHERE name=?", id)

		if err != nil {
			panic(err.Error())
		}
	case "order":
		_, err := db.Exec("DELETE FROM `Orders` WHERE id=?", id)

		if err != nil {
			panic(err.Error())
		}
	case "user":
		_, err := db.Exec("DELETE FROM `Users` WHERE id=?", id)

		if err != nil {
			panic(err.Error())
		}
	case "subcategory":
		_, err := db.Exec("DELETE FROM SubCategories WHERE name = ?", id)

		if err != nil {
			panic(err.Error())
		}
	}
}

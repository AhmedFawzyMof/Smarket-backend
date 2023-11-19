package controller

import (
	"database/sql"
	"fmt"
)

func Edit(db *sql.DB, table string, values map[string]string) {

	switch table {
	case "product":
		id, name, des, price, company, category, subcategories, units, available, offer, inStock, pricePerUint, unitNumber := values["id"], values["name"], values["description"], values["price"], values["company"], values["category"], values["supcategorires"], values["unit"], values["available"], values["offer"], values["inStock"], values["pricePerUint"], values["unitNumber"]
		_, err := db.Exec("UPDATE `Products` SET `name`= ?, `description`=?, `price`=?, `company`=?, `category`=?, `subcategories`=?, `unit`=?, `available`=?, `offer`=?, `inStock`=?, `pricePerUint`=?, `unitNumber`=? WHERE id =?", name, des, price, company, category, subcategories, units, available, offer, inStock, pricePerUint, unitNumber, id)

		if err != nil {
			panic(err.Error())
		}
	case "user":
		id, username, email, phone, spare_phone, role := values["id"], values["username"], values["email"], values["phone"], values["spare_phone"], values["role"]
		_, err := db.Exec("UPDATE `Users` SET `username`= ?, `email`=?, `phone`=?, `spare_phone`=?, `role`=? WHERE id =?", username, email, phone, spare_phone, role, id)

		if err != nil {
			panic(err.Error())
		}

	case "company":
		name, oldname := values["name"], values["oldname"]
		_, err := db.Exec("UPDATE `Companies` SET `name`= ? WHERE name =?", name, oldname)

		if err != nil {
			panic(err.Error())
		}
	case "category":
		name, oldname := values["name"], values["oldname"]
		_, err := db.Exec("UPDATE `Categories` SET `name`= ? WHERE name =?", name, oldname)

		if err != nil {
			panic(err.Error())
		}

	case "subcategory":
		name, oldname, category := values["name"], values["oldname"], values["category"]
		_, err := db.Exec("UPDATE `SubCategories` SET `name`= ?, `category`=? WHERE name =?", name, category, oldname)

		if err != nil {
			panic(err.Error())
		}
	case "order":
		id, confirmed, delivered, paid := values["id"], values["confirmed"], values["delivered"], values["paid"]
		stmt := fmt.Sprintf("UPDATE `Orders` SET confirmed=%s, delivered=%s, paid=%s WHERE id='%s'", confirmed, delivered, paid, id)
		_, err := db.Exec(stmt)

		if err != nil {
			panic(err.Error())
		}
	}

}

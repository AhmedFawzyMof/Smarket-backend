package controller

import (
	"database/sql"
	"fmt"
)

func CheckAddress(db *sql.DB, address map[string]string) map[string]interface{} {
	building := address["building"]
	city := address["city"]
	floor := address["floor"]
	governorate := address["governorate"]
	street := address["street"]
	user := address["UserId"]

	FindEmail := db.QueryRow("SELECT user FROM Address WHERE user = ?", user)

	var User string

	Err := FindEmail.Scan(&User)
	fmt.Println(Err)
	if Err != nil {
		_, err := db.Exec("INSERT INTO `Address`(`user`, `governorate`, `city`, `street`, `building`, `floor`) VALUES (?, ?, ?, ?, ?, ?)", user, governorate, city, street, building, floor)
		if err != nil {
			panic(err.Error())
		}
		Response := map[string]interface{}{
			"Error": false,
		}
		return Response
	}

	Response := map[string]interface{}{
		"Error": true,
	}
	return Response
}

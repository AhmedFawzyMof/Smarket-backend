package DB

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func Connect() *sql.DB {

	db, err := sql.Open("sqlite3", "./smarket.db")

	if err != nil {
		panic(err.Error())
	}

	return db

}

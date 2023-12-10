package DB

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func Connect() *sql.DB {

	db, Err := sql.Open("sqlite3", "AlwadiMarkts.db")

	if Err != nil {
		panic(Err.Error())
	}

	return db

}

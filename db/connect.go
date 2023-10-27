package DB

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

type Config struct {
	databaseString string
}

func Connect() *sql.DB {
	jsonFile, err := os.Open("config.json")

	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()

	byteValue, _ := io.ReadAll(jsonFile)

	var config map[string]interface{}

	json.Unmarshal(byteValue, &config)

	connectString := fmt.Sprintf("%v", config["databaseString"])

	db, Err := sql.Open("mysql", connectString)

	if Err != nil {
		panic(err.Error())
	}

	return db

}

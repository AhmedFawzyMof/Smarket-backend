package DB

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

// type Config struct {
// 	databaseString string
// }

func Connect() *sql.DB {
	// jsonFile, err := os.Open("config.json")

	// if err != nil {
	// 	fmt.Println(err)
	// }
	// defer jsonFile.Close()

	// byteValue, _ := io.ReadAll(jsonFile)

	// var config map[string]interface{}

	// json.Unmarshal(byteValue, &config)

	// connectString := fmt.Sprintf("%v", config["databaseString"])

	db, Err := sql.Open("sqlite3", "smarket.db")

	if Err != nil {
		panic(Err.Error())
	}

	stats := db.Stats()

	fmt.Println("Open connections:", stats.OpenConnections)

	return db

}

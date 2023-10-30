package DB

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type Config struct {
	databaseString string
}

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

	db, Err := sql.Open("mysql", "ssmarketahmed:ssmarketAhmed20@tcp(db4free.net:3306)/ssmarket?parseTime=true")

	if Err != nil {
		panic(Err.Error())
	}

	stats := db.Stats()

	fmt.Println("Open connections:", stats.OpenConnections)

	return db

}

package controller

import (
	"database/sql"
	"fmt"
)

type OrderProducts map[string]int

func CreateOrderProduct(db *sql.DB, order string, products []OrderProducts) map[string]interface{} {
	orderid := order

	values := ""
	for i, key := range products {
		if i == 0 {
			values += fmt.Sprintf("(%d, '%s', %d)", key["id"], orderid, key["quantity"])
		}
		if i > 0 {
			values += fmt.Sprintf(",(%d, '%s', %d)", key["id"], orderid, key["quantity"])
		}
	}

	var sql string = fmt.Sprintf("INSERT INTO OrderProducts (product, `order`, quantity) VALUES %s", values)
	_, err := db.Exec(sql)
	if err != nil {
		panic(err.Error())
	}

	Response := map[string]interface{}{
		"Error": false,
	}
	return Response
}

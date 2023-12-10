package middleware

import (
	"encoding/json"
	"fmt"
)

func SendResponse(response chan<- []byte, Response map[string]interface{}) {
	res, err := json.Marshal(Response)
	if err != nil {
		fmt.Println(err)
	}
	response <- res
}

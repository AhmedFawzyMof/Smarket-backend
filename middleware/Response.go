package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func SendResponse(response chan<- []byte, Response map[string]interface{}) {
	res, err := json.Marshal(Response)
	if err != nil {
		fmt.Println(err)
	}
	response <- res
}

func SendError(err error, res http.ResponseWriter) {
	http.Error(res, err.Error(), http.StatusInternalServerError)
}

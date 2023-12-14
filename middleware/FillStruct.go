package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

func FillStruct(data map[string]interface{}, result interface{}) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	err = json.Unmarshal(jsonData, result)
	if err != nil {
		return err
	}
	return nil
}

func FillStructInterface(data interface{}, result interface{}) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	err = json.Unmarshal(jsonData, result)
	if err != nil {
		return err
	}
	return nil
}

func ConvertToInt(numberInterface interface{}, res http.ResponseWriter) int {
	Int, err := strconv.Atoi(fmt.Sprintf("%s", numberInterface))

	if err != nil {
		SendError(err, res)
	}

	return Int
}

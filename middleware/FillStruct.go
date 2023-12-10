package middleware

import (
	"encoding/json"
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

package cache

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)


var mutex sync.Mutex

func CacheSet(key string, value map[string]interface{}, ttl time.Time) error {
	theFile := fmt.Sprintf("./cache/%s.json", key)
	CacheFile, err := os.Create(theFile)

	if err != nil {
		return fmt.Errorf("%s", err)
	}
	defer CacheFile.Close()

	cache := make(map[string]interface{})

	cache[key] = value
	cache["exp"] = ttl

	encoder := json.NewEncoder(CacheFile)

	err = encoder.Encode(cache)
	if err != nil {
		return fmt.Errorf("%s", err)
	}

	return nil
}

func CacheGet(key string) (map[string]interface{}, error) {
	file := fmt.Sprintf("./cache/%s.json", key)

	jsonFile, err := os.Open(file)

	if err != nil {
		return nil, fmt.Errorf("file not found")
	}
	defer jsonFile.Close()

	byteValue, _ := io.ReadAll(jsonFile)

	var data map[string]interface{}

	json.Unmarshal(byteValue, &data)

	dateStr := fmt.Sprintf("%s", data["exp"])

	dateFromJson, err := time.Parse(time.RFC3339, dateStr)
	if err != nil {
		return nil, fmt.Errorf("%s", err.Error())
	}

	fiveMinutesAgo := time.Now().Add(-5 * time.Minute)

	if dateFromJson.Before(fiveMinutesAgo) {
		jsonFile.Close()
		err := os.Remove(file)
		if err != nil {
			fmt.Println(err)
		}
		return nil, fmt.Errorf("error time out")
	}

	return data, nil
}

package utils

import (
	"encoding/json"
)

func StringToMap(s string) (stringMap map[string]interface{}) {
	if s == "" {
		return
	}

	if err := json.Unmarshal([]byte(s), &stringMap); err != nil {
		panic(err)
	}

	return stringMap
}

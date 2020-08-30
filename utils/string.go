package utils

import (
	"encoding/json"
	"strings"
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

func IsEmpty(s string) bool {
	return len(strings.TrimSpace(s)) == 0
}

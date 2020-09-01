package utils

import (
	"encoding/json"
	"regexp"
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

func IsValidEmail(email string) bool {
	pattern := "^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$"
	return regexp.MustCompile(pattern).MatchString(email)
}

func IsValidUsername(username string) bool {
	pattern := "\\A[^\\\\\\/<>?\\s]+\\z"
	return regexp.MustCompile(pattern).MatchString(username)
}

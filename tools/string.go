package tools

import (
	"encoding/json"
	"golang.org/x/crypto/bcrypt"
	"regexp"
	"strings"
)

func (t *Tool) StringToMap(s string) (stringMap map[string]interface{}) {
	if s == "" {
		return
	}

	if err := json.Unmarshal([]byte(s), &stringMap); err != nil {
		panic(err)
	}

	return stringMap
}

func (t *Tool) IsEmpty(s string) bool {
	return len(strings.TrimSpace(s)) == 0
}

func (t *Tool) IsValidEmail(email string) bool {
	pattern := "^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$"
	return regexp.MustCompile(pattern).MatchString(email)
}

func (t *Tool) IsValidUsername(username string) bool {
	pattern := "\\A[^\\\\\\/<>?\\s]+\\z"
	return regexp.MustCompile(pattern).MatchString(username)
}

// ToHashedPassword returns the bcrypt hash of the password at the given cost
func (t *Tool) ToHashedPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedBytes), nil
}

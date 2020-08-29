package acct

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

var TokenKey = "6cf6813ba69b1e7cf4bedf4fe6d61221"

func GenerateToken(account *Account) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": account.Username,
		"exp":      time.Now().Add(time.Hour * 72).Unix(),
	})

	return token.SignedString([]byte(TokenKey))
}

package acct

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/xandercheung/ogs-go"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"os"
)

var TokenKey = "6cf6813ba69b1e7cf4bedf4fe6d61221"

func TokenAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !isAuthorized(c) {
			response := ogs.RspBase(ogs.StatusInvalidToken, ogs.ErrorMessage("Invalid Token"))
			c.AbortWithStatusJSON(http.StatusOK, response)
			return
		}

		c.Next()
	}
}

// returns the bcrypt hash of the password at the given cost
func passwordToBcryptHash(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedBytes), nil
}

func isAuthorized(c *gin.Context) bool {
	headerToken := c.GetHeader("Authorization")
	token, err := getParseToken(headerToken)

	if err != nil {
		return false
	}

	if !token.Valid {
		return false
	}

	return true
}

func getTokenKey() string {
	tokenKey := os.Getenv("TOKEN_KEY")
	if tokenKey == "" {
		tokenKey = TokenKey
	}

	return tokenKey
}

func tokenSecretKeyFunc() jwt.Keyfunc {
	return func(token *jwt.Token) (interface{}, error) {
		return []byte(getTokenKey()), nil
	}
}

func getParseToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, tokenSecretKeyFunc())
	return token, err
}

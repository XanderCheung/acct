package utils

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/xandercheung/acct"
	"github.com/xandercheung/ogs-go"
	"net/http"
)

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

func tokenSecretKeyFunc() jwt.Keyfunc {
	return func(token *jwt.Token) (interface{}, error) {
		return []byte(acct.TokenKey), nil
	}
}

func getParseToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, tokenSecretKeyFunc())
	return token, err
}

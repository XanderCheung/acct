package acct

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/xandercheung/ogs-go"
	"net/http"
)

// TokenAuthMiddleware middleware of token authentication
func TokenAuthMiddleware() gin.HandlerFunc {
	return func(g *gin.Context) {
		if !isTokenAuthorized(g) {
			response := ogs.RspBase(ogs.StatusInvalidToken, ogs.ErrorMessage("Invalid Token"))
			g.AbortWithStatusJSON(http.StatusOK, response)
			return
		}

		g.Next()
	}
}

func isTokenAuthorized(g *gin.Context) bool {
	headerToken := g.GetHeader("Authorization")
	token, err := getParseToken(headerToken)

	if err != nil {
		return false
	}

	if !token.Valid {
		return false
	}

	return true
}

func getJwtTokenKey() string {
	return jwtTokenKey
}

func tokenSecretKeyFunc() jwt.Keyfunc {
	return func(token *jwt.Token) (interface{}, error) {
		return []byte(getJwtTokenKey()), nil
	}
}

func getParseToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, tokenSecretKeyFunc())
	return token, err
}

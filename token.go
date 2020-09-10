package acct

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/xandercheung/ogs-go"
	"net/http"
	"time"
)

// TokenClaims JWT claims struct
type TokenClaims struct {
	AccountId uint
	jwt.MapClaims
}

// TokenAuthMiddleware middleware of token authentication
func TokenAuthMiddleware() gin.HandlerFunc {
	return func(g *gin.Context) {
		tokenStr := Utils.HeaderToken(g)

		if !isTokenAuthorized(tokenStr) {
			response := ogs.RspError(ogs.StatusInvalidToken, "Invalid Token")
			g.AbortWithStatusJSON(http.StatusOK, response)
			return
		}

		// check account status
		account := Finder.FindAccountByToken(tokenStr)
		if !account.IsPersisted() || account.Status.IsLocked() {
			response := ogs.RspError(ogs.StatusUnauthorized, "Account Status Is Abnormal")
			g.AbortWithStatusJSON(http.StatusOK, response)
			return
		}

		g.Next()
	}
}

func isTokenAuthorized(tokenStr string) bool {
	token, err := jwt.Parse(tokenStr, tokenSecretKeyFunc())

	if err != nil {
		return false
	}

	if !token.Valid {
		return false
	}

	return true
}

func getJwtTokenKey() string {
	return JwtTokenKey
}

func tokenSecretKeyFunc() jwt.Keyfunc {
	return func(token *jwt.Token) (interface{}, error) {
		return []byte(getJwtTokenKey()), nil
	}
}

func findAccountIdByToken(tokenString string) uint {
	tokenClaims := &TokenClaims{}
	_, _ = jwt.ParseWithClaims(tokenString, tokenClaims, tokenSecretKeyFunc())
	return tokenClaims.AccountId
}

func generateTokenByAccountId(accountId uint) (signedStr string, err error) {
	claims := TokenClaims{
		AccountId: accountId,
		MapClaims: jwt.MapClaims{
			"exp": time.Now(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(getJwtTokenKey()))
}

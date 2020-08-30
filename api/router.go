package api

import (
	"github.com/gin-gonic/gin"
	"github.com/xandercheung/acct"
)

func SetAcctRouter(router *gin.Engine) {
	r := router.Group("/api/v1")
	{
		r.POST("/sign_in", signIn)
		r.POST("/sign_up", signUp)

		// Token Authentication
		r.Use(acct.TokenAuthMiddleware())

		accounts := r.Group("/accounts")
		{
			accounts.GET("/", fetchAccounts)
			accounts.POST("/", createAccount)
			accounts.GET("/:id", fetchAccount)
			accounts.PUT("/:id", updateAccount)
			accounts.POST("/:id", updateAccount)
			accounts.DELETE("/:id", destroyAccount)
		}

	}
}

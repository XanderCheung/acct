package acct

import (
	"github.com/gin-gonic/gin"
)

// SetRouter create routers
func SetRouter(router *gin.Engine) {
	r := router.Group("/api/v1")
	{
		r.POST("/sign_in", Handler.SignIn)
		r.POST("/sign_up", Handler.SignUp)

		// Token Authentication
		r.Use(TokenAuthMiddleware())

		accounts := r.Group("/accounts")
		{
			accounts.GET("", Handler.FetchAccounts)
			accounts.POST("", Handler.CreateAccount)
			accounts.GET("/:id", Handler.FetchAccount)
			accounts.POST("/:id", Handler.UpdateAccount)
			accounts.DELETE("/:id", Handler.DestroyAccount)
			accounts.POST("/:id/password", Handler.UpdateAccountPassword)
		}

		r.GET("/account/info", Handler.FetchCurrentAccountInfo)
	}
}

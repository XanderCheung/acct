package service

import (
	"github.com/gin-gonic/gin"
)

func SetAcctV1Router(router *gin.Engine) {
	r := router.Group("/api/v1")
	{
		r.POST("/user/login", Login)
		//r.POST("/user/logout", service.Logout)
		//r.GET("/user/info", service.GetUserInfo)
	}
}

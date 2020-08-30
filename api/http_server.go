package api

import (
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"os"
)

func RunHttpServer() {
	f, _ := os.OpenFile("log/acct-http.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	gin.DefaultWriter = io.MultiWriter(f)

	gin.SetMode(gin.ReleaseMode)

	router := gin.Default()
	SetAcctRouter(router)

	err := router.Run(":2337")
	if err != nil {
		panic(err)
	}
}

func JSON(c *gin.Context, obj interface{}) {
	c.JSON(http.StatusOK, obj)
}

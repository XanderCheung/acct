package acct

import (
	"errors"
	"github.com/gin-gonic/gin"
)

// RunHttpServer starts listening and serving HTTP requests
func RunHttpServer() error {
	if DB == nil {
		return errors.New("connect database failed")
	}

	gin.SetMode(gin.ReleaseMode)

	router := gin.Default()
	SetRouter(router)

	return router.Run(":" + httpServerPort)
}

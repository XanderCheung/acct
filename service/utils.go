package service

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func JSON(c *gin.Context, obj interface{}) {
	c.JSON(http.StatusOK, obj)
}

func IsEmpty(s string) bool {
	return len(strings.TrimSpace(s)) == 0
}

package utils

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
)

func RequestBodyParams(c *gin.Context) (params map[string]interface{}, err error) {
	err = json.NewDecoder(c.Request.Body).Decode(&params)
	return params, err
}

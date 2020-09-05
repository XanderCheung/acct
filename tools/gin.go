package tools

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
)

// RequestBodyParams get body params from request
func (t *Tool) RequestBodyParams(g *gin.Context) (params map[string]interface{}, err error) {
	err = json.NewDecoder(g.Request.Body).Decode(&params)
	return params, err
}

func (t *Tool) JSON(g *gin.Context, obj interface{}) {
	g.JSON(http.StatusOK, obj)
}

func (t *Tool) HeaderToken(g *gin.Context) string {
	return g.GetHeader("Authorization")
}

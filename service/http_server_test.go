package service

import (
	"github.com/gin-gonic/gin"
	"github.com/xandercheung/acct"
	"io"
	"os"
	"testing"
)

func TestRunHttpServer(t *testing.T) {
	mysqlConnectArgs := "root:@tcp(127.0.0.1:3306)/acct_test?&charset=utf8mb4&parseTime=True&loc=UTC"
	mysqlConnection, err := acct.InitDBConnection("mysql", mysqlConnectArgs)
	if err != nil || mysqlConnection == nil {
		t.Error("connect mysql error: ", err)
	}

	acct.DBSeed()

	// 设置日志文件
	f, _ := os.OpenFile("log/acct-http.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	gin.DefaultWriter = io.MultiWriter(f)

	router := gin.Default()
	//router.Use(TokenAuthMiddleware())
	SetAcctV1Router(router)

	_ = router.Run(":2337")
}

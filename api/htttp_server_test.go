package api

import (
	"github.com/gin-gonic/gin"
	"github.com/xandercheung/acct"
	"testing"
)

func TestHttpServer(t *testing.T) {
	mysqlConnectArgs := "root:@tcp(127.0.0.1:3306)/acct_test?&charset=utf8mb4&parseTime=True&loc=UTC"
	mysqlConnection, err := acct.InitDBConnection("mysql", mysqlConnectArgs)
	if err != nil || mysqlConnection == nil {
		t.Error("connect mysql error: ", err)
	}

	acct.DB.DropTableIfExists(&acct.Account{})
	acct.MigrateTables()
	acct.DBSeed()

	router := gin.Default()
	SetAcctRouter(router)

	err = router.Run(":2337")
	if err != nil {
		t.Error("run http server error: ", err)
	}
}

package api

import (
	"github.com/gin-gonic/gin"
	"github.com/xandercheung/acct"
	"testing"
)

func TestHttpServer(t *testing.T) {
	mysqlConnectArgs := "root:@tcp(127.0.0.1:3306)/acct_test?&charset=utf8mb4&parseTime=True&loc=UTC"
	err := acct.InitDBConnection(mysqlConnectArgs)
	if err != nil {
		t.Error("connect mysql error: ", err)
	}

	_ = acct.DB.Migrator().DropTable(&acct.Account{})
	acct.MigrateTables()
	acct.DBSeed()

	router := gin.Default()
	SetAcctRouter(router)

	err = router.Run(":2337")
	if err != nil {
		t.Error("run http server error: ", err)
	}
}

package acct

import (
	"testing"
)

func TestInitDBConnection(t *testing.T) {
	mysqlConnectArgs := "root:@tcp(127.0.0.1:3306)/acct_test?&charset=utf8mb4&parseTime=True&loc=UTC"
	err := InitDBConnection(mysqlConnectArgs)
	if err != nil {
		t.Error("connect mysql error: ", err)
	}
}

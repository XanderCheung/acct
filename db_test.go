package acct

import (
	"testing"
)

func TestInitDBConnection(t *testing.T) {
	mysqlConnectArgs := "root:@tcp(127.0.0.1:3306)/acct_test?&charset=utf8mb4&parseTime=True&loc=UTC"
	mysqlConnection, err := InitDBConnection("mysql", mysqlConnectArgs)
	if err != nil || mysqlConnection == nil {
		t.Error("connect mysql error: ", err)
	}
}

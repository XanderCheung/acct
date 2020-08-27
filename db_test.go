package acct

import (
	"testing"
)

func TestInitDBConnection(t *testing.T) {
	mysqlConnection, err := InitDBConnection("mysql",
		"root:@tcp(127.0.0.1:3306)/acct_test?&charset=utf8mb4&parseTime=True&loc=Local")

	if err != nil || mysqlConnection == nil {
		t.Error("connect mysql error: ", err)
	}

	postgresConnection, err := InitDBConnection("postgres",
		"host=localhost port=5432 user=xander dbname=acct_test password=123456 sslmode=disable")

	if err != nil || postgresConnection == nil {
		t.Error("connect postgres error: ", err)
	}
}

package acct

import (
	"testing"
)

func TestConnectDB(t *testing.T) {
	mysqlConnectArgs := "root:@tcp(127.0.0.1:3306)/acct_test?&charset=utf8mb4&parseTime=True&loc=UTC"
	err := ConnectDB(mysqlConnectArgs)
	if err != nil {
		t.Error("connect mysql error: ", err)
	}
}

func TestInitDBAndSettings(t *testing.T) {
	// default config
	if err := InitDBAndSettings(nil); err != nil {
		t.Error("connect mysql error: ", err)
	}

	if err := InitDBAndSettings(&Config{
		ConnectionDSN:            "root:@tcp(127.0.0.1:3306)/acct_test?&charset=utf8mb4&parseTime=True&loc=UTC",
		IsLoadDSNFromENV:         false,
		IsLoadJwtTokenKeyFromENV: false,
		JwtTokenKey:              "1cf6813ba62k2e7cf4bedf4fe6d69kd7",
		HttpServerPort:           "2337",
	}); err != nil {
		t.Error("connect mysql error: ", err)
	}
}

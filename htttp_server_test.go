package acct

import (
	"testing"
)

func TestHttpServer(t *testing.T) {
	if err := InitDBAndSettings(&Config{
		ConnectionDSN:            "root:@tcp(127.0.0.1:3306)/acct_test?&charset=utf8mb4&parseTime=True&loc=UTC",
		IsLoadDSNFromENV:         false,
		IsLoadJwtTokenKeyFromENV: false,
		JwtTokenKey:              "1cf6813ba62k2e7cf4bedf4fe6d69kd7",
		HttpServerPort:           "2337",
	}); err != nil {
		t.Error("connect mysql error: ", err)
	}

	if err := RunHttpServer(); err != nil {
		t.Error("connect mysql error: ", err)
	}
}

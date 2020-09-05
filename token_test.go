package acct

import (
	"testing"
)

func TestGenerateToken(t *testing.T) {
	var accountId uint = 1
	tokenStr, err := generateTokenByAccountId(accountId)
	if err != nil {
		t.Error("generate token error: ", err)
	}

	accountIdByToken := findAccountIdByToken(tokenStr)
	if accountId != accountIdByToken {
		t.Error("find account id by token error: ", err)
	}
}

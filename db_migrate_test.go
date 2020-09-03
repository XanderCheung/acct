package acct

import (
	"testing"
)

func TestMigrateTables(t *testing.T) {
	TestConnectDB(t)

	_ = DB.Migrator().DropTable(&Account{})
	if err := migrateTables(); err != nil {
		t.Error("migrate mysql tables error")
	}

	if !DB.Migrator().HasTable(&Account{}) {
		t.Error("migrate mysql tables error")
	}
}

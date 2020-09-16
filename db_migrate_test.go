package acct

import (
	"testing"
)

func TestMigrateTables(t *testing.T) {
	TestConnectDB(t)

	_ = DB.Migrator().DropTable("accounts")
	_ = DB.Migrator().DropTable("schema_migrations")
	if err := MigrateTables(); err != nil {
		t.Error("migrate mysql tables error: " + err.Error())
	}

	if !DB.Migrator().HasTable(&Account{}) {
		t.Error("migrate mysql tables error")
	}
}

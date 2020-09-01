package acct

import (
	"testing"
)

func TestMigrateTables(t *testing.T) {
	TestInitDBConnection(t)

	_ = DB.Migrator().DropTable(&Account{})
	MigrateTables()

	if !DB.Migrator().HasTable(&Account{}) {
		t.Error("migrate mysql tables error")
	}
}

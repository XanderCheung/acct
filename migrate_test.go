package acct

import (
	"testing"
)

func TestMigrateTables(t *testing.T) {
	TestInitDBConnection(t)

	DB.DropTableIfExists(&Account{})
	MigrateTables()

	if !DB.HasTable(&Account{}) {
		t.Error("migrate mysql tables error")
	}
}

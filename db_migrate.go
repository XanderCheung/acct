package acct

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/mysql"
	_ "github.com/golang-migrate/migrate/source/file"
)

func MigrateTables() error {
	db, _ := DB.DB()
	driver, err := mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"mysql",
		driver,
	)

	if err != nil {
		return err
	}

	_ = m.Steps(5)
	return nil
}

func MigrateSeeds() error {
	if !Finder.IsAccountWithDeletedExists(map[string]interface{}{"username": "admin"}, nil) {
		account := Account{
			Email:    "admin@qq.com",
			Username: "admin",
			Nickname: "Admin",
			Password: "admin@123456"}
		if err := account.Create(); err != nil {
			return err
		}
	}
	return nil
}

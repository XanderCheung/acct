package acct

func migrateTables() error {
	return DB.Set("gorm:table_options", "CHARSET=utf8mb4").Debug().AutoMigrate(
		&Account{},
	)
}

func migrateSeeds() error {
	if !Finder.IsAccountWithDeletedExists(map[string]interface{}{"username": "admin"}, nil) {
		account := Account{Email: "admin@qq.com", Username: "admin", Password: "admin@123456"}
		if err := account.Create(); err != nil {
			return err
		}
	}
	return nil
}

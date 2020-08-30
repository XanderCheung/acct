package acct

func MigrateTables() {
	DB.Set("gorm:table_options", "CHARSET=utf8mb4").Debug().AutoMigrate(
		&Account{},
	)
}

func DBSeed() {
	if !IsAccountWithDeletedExists(map[string]interface{}{"username": "admin"}, nil) {
		account := Account{Email: "admin@qq.com", Username: "admin", Password: "admin@123456"}
		_ = account.Create()
	}
}

package acct

type finder struct {
}

var Finder = &finder{}

// IsAccountWithDeletedExists check is account exists
func (f *finder) IsAccountWithDeletedExists(query, notQuery map[string]interface{}) bool {
	account := Account{}
	DB.Unscoped().Model(&Account{}).Where(query).Not(notQuery).Limit(1).Find(&account)
	return account.ID > 0
}

// FindAccountById find account by id
func (f *finder) FindAccountById(id uint) Account {
	DB.Migrator()
	account := Account{}
	DB.Model(&Account{}).Limit(1).Find(&account, id)
	return account
}

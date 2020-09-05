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
	account := Account{}
	DB.Model(&Account{}).Limit(1).Find(&account, id)
	return account
}

// FindAccountByToken find account by token
func (f *finder) FindAccountByToken(tokenString string) Account {
	id := f.FindAccountIdByToken(tokenString)
	return f.FindAccountById(id)
}

// FindAccountIdByToken find account id by token
func (f *finder) FindAccountIdByToken(tokenString string) uint {
	return findAccountIdByToken(tokenString)
}

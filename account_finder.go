package acct

// IsAccountWithDeletedExists check is account exists
func IsAccountWithDeletedExists(query, notQuery map[string]interface{}) bool {
	account := Account{}
	DB.Unscoped().Model(&Account{}).Where(query).Not(notQuery).Limit(1).Find(&account)
	return account.ID > 0
}

// FindAccountById find account by id
func FindAccountById(id uint) Account {
	account := Account{}
	DB.Model(&Account{}).Take(&account, id)
	return account
}

func (c *Account) IsPersisted() bool {
	return c.ID > 0
}

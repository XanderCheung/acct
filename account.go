package acct

import (
	"errors"
	"github.com/jinzhu/gorm"
	"strings"
	"sync"
	"time"
)

type Account struct {
	ID        uint       `gorm:"primary_key" json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"-"`
	Email     string     `gorm:"type:varchar(100);unique_index;not null" json:"email"`
	Username  string     `gorm:"type:varchar(100);unique_index;not null" json:"username"`
	Password  string     `gorm:"column:encrypted_password;not null" json:"-"`
	Token     string     `gorm:"-" json:"token"`
}

var accountLock = sync.RWMutex{}

func (c *Account) BeforeSave(scope *gorm.Scope) (err error) {
	if err = c.Validate(); err != nil {
		return err
	}

	if err = c.encryptPassword(scope); err != nil {
		return err
	}
	return
}

func (c *Account) encryptPassword(scope *gorm.Scope) (err error) {
	if c.Password == "" {
		return nil
	}

	hashedPassword, err := passwordToBcryptHash(c.Password)
	if err != nil {
		return err
	}

	return scope.SetColumn("encrypted_password", hashedPassword)
}

func IsAccountExists(query, notQuery map[string]interface{}) bool {
	account := Account{}
	DB.Model(&Account{}).Where(query).Not(notQuery).Limit(1).Find(&account)
	return account.ID > 0
}

func (c *Account) Create() error {
	accountLock.Lock()
	defer accountLock.Unlock()

	if err := DB.Create(&c).Error; err != nil {
		return err
	}

	if c.ID == 0 {
		return errors.New("create account failed")
	}

	return nil
}

func (c *Account) Validate() error {
	if !strings.Contains(c.Email, "@") {
		return errors.New("email address is required")
	}

	if len(c.Password) < 6 {
		return errors.New("password is to short")
	}

	if IsAccountExists(map[string]interface{}{"email": c.Email}, map[string]interface{}{"id": c.ID}) {
		return errors.New("email address already in use")
	}

	if IsAccountExists(map[string]interface{}{"username": c.Username}, map[string]interface{}{"id": c.ID}) {
		return errors.New("username already in use")
	}

	return nil
}

func (c *Account) GenerateToken() error {
	token, err := GenerateToken(c)
	if err != nil {
		return err
	}

	c.Token = token

	return nil
}

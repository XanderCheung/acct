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

func (c *Account) Create() error {
	accountLock.Lock()
	defer accountLock.Unlock()

	if isValid, err := c.Validate(); !isValid {
		return err
	}

	hashPassword, err := passwordToBcryptHash(c.Password)
	if err != nil {
		return err
	}

	c.Password = hashPassword

	err = DB.Create(&c).Error
	if err != nil {
		return err
	}

	if c.ID == 0 {
		return errors.New("create account failed")
	}

	return nil
}

func (c *Account) Validate() (isValid bool, err error) {
	if !strings.Contains(c.Email, "@") {
		return false, errors.New("email address is required")
	}

	if len(c.Password) < 6 {
		return false, errors.New("password is to short")
	}

	temp := Account{}
	err = DB.Model(&Account{}).Where("email = ?", c.Email).First(&temp).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return false, errors.New("connection error. please retry")
	}

	if temp.ID > 0 {
		return false, errors.New("email address already in use")
	}

	temp = Account{}
	err = DB.Model(&Account{}).Where("username = ?", c.Username).First(&temp).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return false, errors.New("connection error. please retry")
	}

	if temp.ID > 0 {
		return false, errors.New("username already in use")
	}

	return true, nil
}

func (c *Account) GenerateToken() error {
	token, err := GenerateToken(c)
	if err != nil {
		return err
	}

	c.Token = token

	return nil
}

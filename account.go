package acct

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"gorm.io/gorm"
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

func (c *Account) BeforeSave(tx *gorm.DB) (err error) {
	if err = c.Validate(); err != nil {
		return err
	}

	if err = c.encryptPassword(tx); err != nil {
		return err
	}
	return
}

func (c Account) encryptPassword(tx *gorm.DB) (err error) {
	hashedPassword, err := passwordToBcryptHash(c.Password)
	if err != nil {
		return err
	}

	tx.Statement.SetColumn("encrypted_password", hashedPassword)
	return nil
}

func IsAccountWithDeletedExists(query, notQuery map[string]interface{}) bool {
	account := Account{}
	DB.Unscoped().Model(&Account{}).Where(query).Not(notQuery).Limit(1).Find(&account)
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

	if IsAccountWithDeletedExists(map[string]interface{}{"email": c.Email}, map[string]interface{}{"id": c.ID}) {
		return errors.New("email address already in use")
	}

	if IsAccountWithDeletedExists(map[string]interface{}{"username": c.Username}, map[string]interface{}{"id": c.ID}) {
		return errors.New("username already in use")
	}

	return nil
}

func (c *Account) GenerateToken() error {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": c.Username,
		"exp":      time.Now().Add(time.Hour * 72).Unix(),
	})

	signedStr, err := token.SignedString([]byte(getTokenKey()))
	if err != nil {
		return err
	}

	c.Token = signedStr

	return nil
}

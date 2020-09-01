package acct

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/xandercheung/acct/utils"
	"gorm.io/gorm"
	"sync"
	"time"
)

type Account struct {
	ID                uint           `gorm:"primaryKey" json:"id"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	DeletedAt         gorm.DeletedAt `json:"-"`
	Email             string         `gorm:"type:varchar(100);uniqueIndex;not null" json:"email"`
	Username          string         `gorm:"type:varchar(100);uniqueIndex;not null" json:"username"`
	EncryptedPassword string         `gorm:"column:encrypted_password;not null" json:"-"`
	Password          string         `gorm:"column:encrypted_password;not null" json:"-"`
	Token             string         `gorm:"-" json:"token"`
}

var AccountLock = sync.RWMutex{}

func (c *Account) BeforeSave(tx *gorm.DB) (err error) {
	if err = c.Validate(); err != nil {
		return err
	}

	if err = c.encryptPassword(tx); err != nil {
		return err
	}
	return
}

func (c *Account) encryptPassword(tx *gorm.DB) (err error) {
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
	AccountLock.Lock()
	defer AccountLock.Unlock()

	if err := DB.Create(&c).Error; err != nil {
		return err
	}

	if c.ID == 0 {
		return errors.New("create account failed")
	}

	return nil
}

func (c *Account) Validate() error {
	if err := c.validateEmail(); err != nil {
		return err
	}

	if err := c.validateUsername(); err != nil {
		return err
	}

	if err := c.validatePassword(); err != nil {
		return err
	}

	return nil
}

func (c *Account) validateEmail() error {
	if utils.IsEmpty(c.Email) {
		return errors.New("email is required")
	}

	if !utils.IsValidEmail(c.Email) {
		return errors.New("email address is invalid")
	}

	if IsAccountWithDeletedExists(map[string]interface{}{"email": c.Email}, map[string]interface{}{"id": c.ID}) {
		return errors.New("email address already in use")
	}

	return nil
}

func (c *Account) validateUsername() error {
	if utils.IsEmpty(c.Username) {
		return errors.New("username is required")
	}

	if !utils.IsValidUsername(c.Username) {
		return errors.New("username is invalid")
	}

	if IsAccountWithDeletedExists(map[string]interface{}{"username": c.Username}, map[string]interface{}{"id": c.ID}) {
		return errors.New("username already in use")
	}

	return nil
}

func (c *Account) validatePassword() error {
	// TODO not set Password when use map to update
	if c.Password == "" {
		return nil
	}

	if len(c.Password) < 6 {
		return errors.New("password is too short, the minimum is 6 characters")
	}

	if len(c.Password) > 50 {
		return errors.New("password is too long. The maximum length is 50 characters")
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

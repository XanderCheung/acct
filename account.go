package acct

import (
	"github.com/dgrijalva/jwt-go"
	"gorm.io/gorm"
	"sync"
	"time"
)

type Account struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-"`
	Email     string         `gorm:"type:varchar(100);uniqueIndex;not null" json:"email"`
	Username  string         `gorm:"type:varchar(100);uniqueIndex;not null" json:"username"`
	Password  string         `gorm:"column:encrypted_password;not null" json:"-"`
	Token     string         `gorm:"-" json:"token"`
}

var accountLock = sync.RWMutex{}

// Create create account
// permitted columns: "Email", "Username", "Password"
func (c *Account) Create() error {
	accountLock.Lock()
	defer accountLock.Unlock()

	if err := c.validateEmail(); err != nil {
		return err
	}

	if err := c.validateUsername(); err != nil {
		return err
	}

	err := c.validatePassword()
	if err != nil {
		return err
	}

	hashedPassword, err := Utils.ToHashedPassword(c.Password)
	if err != nil {
		return err
	}

	c.Password = hashedPassword

	return DB.Create(&c).Error
}

// Update update account
// permitted columns: "Email", "Username"
func (c *Account) Update() error {
	accountLock.Lock()
	defer accountLock.Unlock()

	if err := c.validateEmail(); err != nil {
		return err
	}

	if err := c.validateUsername(); err != nil {
		return err
	}

	updateParams := Account{
		Email:    c.Email,
		Username: c.Username,
	}

	return DB.Model(&c).Updates(updateParams).Error
}

// UpdatePassword update password
// permitted columns: "Password"
func (c *Account) UpdatePassword() error {
	err := c.validatePassword()
	if err != nil {
		return err
	}

	hashedPassword, err := Utils.ToHashedPassword(c.Password)
	if err != nil {
		return err
	}

	return DB.Model(&c).UpdateColumn("encrypted_password", hashedPassword).Error
}

// GenerateToken generate new token to field of "Token"
func (c *Account) GenerateToken() error {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": c.Username,
		"exp":      time.Now().Add(time.Hour * 72).Unix(),
	})

	signedStr, err := token.SignedString([]byte(getJwtTokenKey()))
	if err != nil {
		return err
	}

	c.Token = signedStr

	return nil
}

func (c *Account) IsPersisted() bool {
	return c.ID > 0
}

package acct

import (
	"errors"
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
	Nickname  string         `gorm:"type:varchar(100)" json:"nickname"`
	Avatar    string         `gorm:"type:varchar(100)" json:"avatar"`
	Status    AccountStatus  `gorm:"type:smallint;default:0" json:"status"`
	Token     string         `gorm:"-" json:"token,omitempty"`
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
		Nickname: c.Nickname,
		Avatar:   c.Avatar,
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

// Lock lock account, unable to login
func (c *Account) Lock() error {
	if c.Status.IsLocked() {
		return errors.New("already locked")
	}

	return DB.Model(&c).UpdateColumn("status", AccountStatusLocked).Error
}

// UnLock unlock account
func (c *Account) UnLock() error {
	if c.Status.IsNormal() {
		return errors.New("already unlocked")
	}

	return DB.Model(&c).UpdateColumn("status", AccountStatusNormal).Error
}

// GenerateToken generate new token to field of "Token"
func (c *Account) GenerateToken() error {
	token, err := generateTokenByAccountId(c.ID)
	c.Token = token
	return err
}

// IsPersisted return true if the record is persisted
func (c *Account) IsPersisted() bool {
	return c.ID > 0
}

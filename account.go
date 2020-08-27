package acct

import (
	"github.com/jinzhu/gorm"
)

type Account struct {
	gorm.Model
	Email    string `gorm:"type:varchar(100);unique_index"json:"email"`
	Username string `gorm:"type:varchar(100);unique_index"json:"username"`
	Password string `gorm:"column:encrypted_password" json:"-"`
	Token    string `gorm:"-" json:"token"`
}

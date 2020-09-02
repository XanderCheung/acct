package acct

import (
	"errors"
	"github.com/xandercheung/acct/utils"
)

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
	if utils.IsEmpty(c.Password) {
		return errors.New("password is required")
	}

	if len(c.Password) < 6 {
		return errors.New("password is too short, the minimum is 6 characters")
	}

	if len(c.Password) > 50 {
		return errors.New("password is too long. The maximum length is 50 characters")
	}

	return nil
}

package acct

import "golang.org/x/crypto/bcrypt"

// returns the bcrypt hash of the password at the given cost
func passwordToBcryptHash(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedBytes), nil
}

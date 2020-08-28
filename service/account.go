package service

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/xandercheung/acct"
	"github.com/xandercheung/ogs-go"
	"golang.org/x/crypto/bcrypt"
)

func Login(c *gin.Context) {
	email := c.PostForm("email")
	username := c.PostForm("username")
	password := c.PostForm("password")

	if IsEmpty(password) || IsEmpty(username) && IsEmpty(password) {
		JSON(c, ogs.RspBase(ogs.StatusSystemError, ogs.ErrorMessage("Invalid Request")))
		return
	}

	// Find user
	account := acct.Account{}
	err := acct.DB.Where(&acct.Account{Email: email, Username: username}).First(&account).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			JSON(c, ogs.RspBase(ogs.StatusUserNotFound, ogs.ErrorMessage("User Not Found")))
			return
		}

		JSON(c, ogs.RspBase(ogs.StatusSystemError, ogs.ErrorMessage("Connection Error. Please Retry")))
		return
	}

	// Compare password
	err = bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(password))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword { // Password does not match!
			errMsg := ogs.CodeText(ogs.StatusErrorPassword)
			JSON(c, ogs.RspBase(ogs.StatusErrorPassword, ogs.ErrorMessage(errMsg)))
			return
		}

		JSON(c, ogs.RspBase(ogs.StatusSystemError, ogs.ErrorMessage(err.Error())))
		return
	}

	// Generate token
	err = account.GenerateToken()
	if err != nil {
		JSON(c, ogs.RspBase(ogs.StatusSystemError, ogs.ErrorMessage("Login Failed")))
		return
	}

	JSON(c, ogs.RspOKWithData(ogs.SuccessMessage("Signed In"), account))
}

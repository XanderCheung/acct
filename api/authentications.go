package api

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/xandercheung/acct"
	"github.com/xandercheung/acct/utils"
	"github.com/xandercheung/ogs-go"
	"golang.org/x/crypto/bcrypt"
)

type tempAccount struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func signIn(c *gin.Context) {
	temp := tempAccount{}
	if err := json.NewDecoder(c.Request.Body).Decode(&temp); err != nil || utils.IsEmpty(temp.Email) && utils.IsEmpty(temp.Username) {
		JSON(c, ogs.RspBase(ogs.StatusSystemError, ogs.ErrorMessage("Invalid Request")))
		return
	}

	// Find user
	account := acct.Account{}
	if err := acct.DB.Where(&acct.Account{Email: temp.Email, Username: temp.Username}).
		Limit(1).Find(&account).Error; err != nil {

		if err == gorm.ErrRecordNotFound {
			JSON(c, ogs.RspBase(ogs.StatusUserNotFound, ogs.ErrorMessage("User Not Found")))
			return
		}

		JSON(c, ogs.RspBase(ogs.StatusSystemError, ogs.ErrorMessage("Connection Error. Please Retry")))
		return
	}

	// Compare password
	if err := bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(temp.Password)); err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword { // Password does not match!
			errMsg := ogs.CodeText(ogs.StatusErrorPassword)
			JSON(c, ogs.RspBase(ogs.StatusErrorPassword, ogs.ErrorMessage(errMsg)))
			return
		}

		JSON(c, ogs.RspBase(ogs.StatusSystemError, ogs.ErrorMessage(err.Error())))
		return
	}

	// Generate token
	if err := account.GenerateToken(); err != nil {
		JSON(c, ogs.RspBase(ogs.StatusSystemError, ogs.ErrorMessage("Sign In Failed")))
		return
	}

	JSON(c, ogs.RspOKWithData(ogs.SuccessMessage("Signed In"), account))
}

func signUp(c *gin.Context) {
	temp := tempAccount{}
	if err := json.NewDecoder(c.Request.Body).Decode(&temp); err != nil {
		JSON(c, ogs.RspBase(ogs.StatusSystemError, ogs.ErrorMessage("Invalid Request")))
		return
	}

	account := acct.Account{
		Email:    temp.Email,
		Username: temp.Username,
		Password: temp.Password,
	}

	if err := account.Create(); err != nil {
		JSON(c, ogs.RspBase(ogs.StatusSystemError, ogs.ErrorMessage(err.Error())))
		return
	}

	JSON(c, ogs.RspOKWithData(ogs.SuccessMessage("Registered Successfully"), account))
}

package acct

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/xandercheung/ogs-go"
	"golang.org/x/crypto/bcrypt"
)

type tempAccount struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func (c *handler) SignIn(g *gin.Context) {
	temp := tempAccount{}
	if err := json.NewDecoder(g.Request.Body).Decode(&temp); err != nil || Utils.IsEmpty(temp.Email) && Utils.IsEmpty(temp.Username) {
		Utils.JSON(g, ogs.RspBase(ogs.StatusSystemError, ogs.ErrorMessage("Invalid Request")))
		return
	}

	// Find user
	account := Account{}
	DB.Where(&Account{Email: temp.Email, Username: temp.Username}).Limit(1).Find(&account)

	if !account.IsPersisted() {
		Utils.JSON(g, ogs.RspBase(ogs.StatusUserNotFound, ogs.ErrorMessage("User Not Found")))
		return
	}

	// Compare password
	if err := bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(temp.Password)); err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword { // Password does not match!
			errMsg := ogs.CodeText(ogs.StatusErrorPassword)
			Utils.JSON(g, ogs.RspBase(ogs.StatusErrorPassword, ogs.ErrorMessage(errMsg)))
			return
		}

		Utils.JSON(g, ogs.RspBase(ogs.StatusSystemError, ogs.ErrorMessage(err.Error())))
		return
	}

	// Generate token
	if err := account.GenerateToken(); err != nil {
		Utils.JSON(g, ogs.RspBase(ogs.StatusSystemError, ogs.ErrorMessage("Sign In Failed")))
		return
	}

	Utils.JSON(g, ogs.RspOKWithData(ogs.SuccessMessage("Signed In"), account))
}

func (c *handler) SignUp(g *gin.Context) {
	temp := tempAccount{}
	if err := json.NewDecoder(g.Request.Body).Decode(&temp); err != nil {
		Utils.JSON(g, ogs.RspBase(ogs.StatusSystemError, ogs.ErrorMessage("Invalid Request")))
		return
	}

	account := Account{
		Email:    temp.Email,
		Username: temp.Username,
		Password: temp.Password,
	}

	if err := account.Create(); err != nil {
		Utils.JSON(g, ogs.RspBase(ogs.StatusSystemError, ogs.ErrorMessage(err.Error())))
		return
	}

	Utils.JSON(g, ogs.RspOKWithData(ogs.SuccessMessage("Registered Successfully"), account))
}

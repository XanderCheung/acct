package acct

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/xandercheung/ogs-go"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

type tempAccount struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
}

func (c *handler) SignIn(g *gin.Context) {
	temp := tempAccount{}
	if err := json.NewDecoder(g.Request.Body).Decode(&temp); err != nil ||
		Utils.IsEmpty(temp.Email) && Utils.IsEmpty(temp.Username) {

		Utils.JSON(g, ogs.RspError(ogs.StatusBadParams, "Bad Params"))
		return
	}

	// Find user
	account := Account{}
	DB.Where(&Account{Email: temp.Email, Username: temp.Username}).Limit(1).Find(&account)

	if !account.IsPersisted() {
		Utils.JSON(g, ogs.RspError(ogs.StatusUserNotFound, "User Not Found"))
		return
	}

	// check account status
	if account.Status.IsLocked() {
		response := ogs.RspError(ogs.StatusUnauthorized, "Account Status Is Abnormal")
		g.AbortWithStatusJSON(http.StatusOK, response)
		return
	}

	// Compare password
	if err := bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(temp.Password)); err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword { // Password does not match!
			errMsg := ogs.CodeText(ogs.StatusErrorPassword)
			Utils.JSON(g, ogs.RspError(ogs.StatusErrorPassword, errMsg))
			return
		}

		Utils.JSON(g, ogs.RspError(ogs.StatusSystemError, err.Error()))
		return
	}

	// Generate token
	if err := account.GenerateToken(); err != nil {
		Utils.JSON(g, ogs.RspError(ogs.StatusSystemError, "Sign In Failed"))
		return
	}

	Utils.JSON(g, ogs.RspOKWithData("Signed In", account))
}

func (c *handler) SignUp(g *gin.Context) {
	temp := tempAccount{}
	if err := json.NewDecoder(g.Request.Body).Decode(&temp); err != nil {
		Utils.JSON(g, ogs.RspError(ogs.StatusBadParams, "Bad Params"))
		return
	}

	account := Account{
		Email:    temp.Email,
		Username: temp.Username,
		Password: temp.Password,
	}

	if err := account.Create(); err != nil {
		Utils.JSON(g, ogs.RspError(ogs.StatusSignUpFailed, err.Error()))
		return
	}

	Utils.JSON(g, ogs.RspOKWithData("Registered Successfully", account))
}

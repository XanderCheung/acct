package api

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/xandercheung/acct"
	"github.com/xandercheung/acct/utils"
	"github.com/xandercheung/ogs-go"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"strconv"
)

func fetchAccounts(c *gin.Context) {
	var accounts []acct.Account
	relation := acct.DB.Model(&acct.Account{}).Order("id desc")

	queryConditions := utils.StringToMap(c.Query("q"))
	if len(queryConditions) > 0 {
		if username, ok := queryConditions["username"].(string); ok {
			relation = relation.Where("username LIKE ?", "%"+username+"%")
		}

		if email, ok := queryConditions["email"].(string); ok {
			relation = relation.Where("name LIKE ?", "%"+email+"%")
		}
	}

	relation, paginate := utils.PaginateGin(relation, c)
	relation.Find(&accounts)

	JSON(c, ogs.RspOKWithPaginate(ogs.BlankMessage(), accounts, paginate))
}

func fetchAccount(c *gin.Context) {
	account, err := loadAccount(c)
	if err != nil {
		return
	}

	JSON(c, ogs.RspOKWithData(ogs.BlankMessage(), account))
}

func createAccount(c *gin.Context) {
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

	err := account.Create()
	if err != nil {
		JSON(c, ogs.RspBase(ogs.StatusCreateFailed, ogs.ErrorMessage(err.Error())))
		return
	}

	JSON(c, ogs.RspOKWithData(ogs.SuccessMessage("Create Successfully"), account))
}

func updateAccount(c *gin.Context) {
	account, err := loadAccount(c)
	if err != nil {
		return
	}

	temp := tempAccount{}
	if err = json.NewDecoder(c.Request.Body).Decode(&temp); err != nil {
		JSON(c, ogs.RspBase(ogs.StatusSystemError, ogs.ErrorMessage("Invalid Request")))
		return
	}

	account.Email = temp.Email
	account.Username = temp.Username

	if err = account.Update(); err != nil {
		JSON(c, ogs.RspBase(ogs.StatusUpdateFailed, ogs.ErrorMessage(err.Error())))
	} else {
		JSON(c, ogs.RspOK(ogs.SuccessMessage("Update Successfully")))
	}
}

func destroyAccount(c *gin.Context) {
	account, err := loadAccount(c)
	if err != nil {
		return
	}

	if err = acct.DB.Delete(&account).Error; err != nil {
		JSON(c, ogs.RspBase(ogs.StatusDestroyFailed, ogs.ErrorMessage("Destroy Failed")))
	} else {
		JSON(c, ogs.RspOK(ogs.SuccessMessage("Destroy Successfully")))
	}
}

func updateAccountPassword(c *gin.Context) {
	account, err := loadAccount(c)
	if err != nil {
		return
	}

	params, err := utils.RequestBodyParams(c)
	if err != nil {
		JSON(c, ogs.RspBase(ogs.StatusSystemError, ogs.ErrorMessage("Invalid Request")))
		return
	}

	newPassword, _ := params["new_password"].(string)
	oldPassword, _ := params["old_password"].(string)

	if err = bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(oldPassword)); err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword { // Password does not match!
			errMsg := ogs.CodeText(ogs.StatusErrorPassword)
			JSON(c, ogs.RspBase(ogs.StatusErrorPassword, ogs.ErrorMessage(errMsg)))
			return
		}

		JSON(c, ogs.RspBase(ogs.StatusSystemError, ogs.ErrorMessage(err.Error())))
		return
	}

	account.Password = newPassword
	if err = account.UpdatePassword(); err != nil {
		JSON(c, ogs.RspBase(ogs.StatusUpdateFailed, ogs.ErrorMessage(err.Error())))
	} else {
		JSON(c, ogs.RspOK(ogs.SuccessMessage("Update Password Successfully")))
	}
}

func loadAccount(c *gin.Context) (account acct.Account, err error) {
	id, _ := strconv.Atoi(c.Param("id"))
	account = acct.FindAccountById(uint(id))

	if !account.IsPersisted() {
		JSON(c, ogs.RspBase(ogs.StatusUserNotFound, ogs.ErrorMessage("Account Not Found")))
		return account, gorm.ErrRecordNotFound
	}

	return account, nil
}

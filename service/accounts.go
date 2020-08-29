package service

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/xandercheung/acct"
	"github.com/xandercheung/acct/utils"
	"github.com/xandercheung/ogs-go"
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
	id := c.Param("id")
	account := acct.Account{}
	if IsEmpty(id) || acct.DB.Limit(1).Find(&account, c.Param("id")).RecordNotFound() {
		JSON(c, ogs.RspBase(ogs.StatusUserNotFound, ogs.ErrorMessage("Account Not Found")))
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
	id := c.Param("id")
	account := acct.Account{}
	if IsEmpty(id) || acct.DB.Limit(1).Find(&account, c.Param("id")).RecordNotFound() {
		JSON(c, ogs.RspBase(ogs.StatusUserNotFound, ogs.ErrorMessage("Account Not Found")))
		return
	}

	temp := tempAccount{}
	if err := json.NewDecoder(c.Request.Body).Decode(&temp); err != nil {
		JSON(c, ogs.RspBase(ogs.StatusSystemError, ogs.ErrorMessage("Invalid Request")))
		return
	}

	if err := acct.DB.Debug().Model(&account).Updates(acct.Account{
		Email:    temp.Email,
		Username: temp.Username,
		Password: temp.Password,
	}).Error; err != nil {
		JSON(c, ogs.RspBase(ogs.StatusUpdateFailed, ogs.ErrorMessage(err.Error())))
	} else {
		JSON(c, ogs.RspOK(ogs.SuccessMessage("Update Successfully")))
	}
}

func destroyAccount(c *gin.Context) {
	id := c.Param("id")
	account := acct.Account{}
	if IsEmpty(id) || acct.DB.Limit(1).Find(&account, c.Param("id")).RecordNotFound() {
		JSON(c, ogs.RspBase(ogs.StatusUserNotFound, ogs.ErrorMessage("Account Not Found")))
		return
	}

	if err := acct.DB.Delete(&account).Error; err != nil {
		JSON(c, ogs.RspBase(ogs.StatusDestroyFailed, ogs.ErrorMessage("Destroy Failed")))
	} else {
		JSON(c, ogs.RspOK(ogs.SuccessMessage("Destroy Successfully")))
	}
}

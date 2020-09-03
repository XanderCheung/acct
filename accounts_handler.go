package acct

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/xandercheung/ogs-go"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"strconv"
)

func (c *handler) FetchAccounts(g *gin.Context) {
	var accounts []Account
	relation := DB.Model(&Account{}).Order("id desc")

	queryConditions := Utils.StringToMap(g.Query("q"))
	if len(queryConditions) > 0 {
		if username, ok := queryConditions["username"].(string); ok {
			relation = relation.Where("username LIKE ?", "%"+username+"%")
		}

		if email, ok := queryConditions["email"].(string); ok {
			relation = relation.Where("name LIKE ?", "%"+email+"%")
		}
	}

	relation, paginate := Utils.PaginateGin(relation, g)
	relation.Find(&accounts)

	Utils.JSON(g, ogs.RspOKWithPaginate(ogs.BlankMessage(), accounts, paginate))
}

func (c *handler) FetchAccount(g *gin.Context) {
	account, err := c.loadAccount(g)
	if err != nil {
		return
	}

	Utils.JSON(g, ogs.RspOKWithData(ogs.BlankMessage(), account))
}

func (c *handler) CreateAccount(g *gin.Context) {
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

	err := account.Create()
	if err != nil {
		Utils.JSON(g, ogs.RspBase(ogs.StatusCreateFailed, ogs.ErrorMessage(err.Error())))
		return
	}

	Utils.JSON(g, ogs.RspOKWithData(ogs.SuccessMessage("Create Successfully"), account))
}

func (c *handler) UpdateAccount(g *gin.Context) {
	account, err := c.loadAccount(g)
	if err != nil {
		return
	}

	temp := tempAccount{}
	if err = json.NewDecoder(g.Request.Body).Decode(&temp); err != nil {
		Utils.JSON(g, ogs.RspBase(ogs.StatusSystemError, ogs.ErrorMessage("Invalid Request")))
		return
	}

	account.Email = temp.Email
	account.Username = temp.Username

	if err = account.Update(); err != nil {
		Utils.JSON(g, ogs.RspBase(ogs.StatusUpdateFailed, ogs.ErrorMessage(err.Error())))
	} else {
		Utils.JSON(g, ogs.RspOK(ogs.SuccessMessage("Update Successfully")))
	}
}

func (c *handler) DestroyAccount(g *gin.Context) {
	account, err := c.loadAccount(g)
	if err != nil {
		return
	}

	if err = DB.Delete(&account).Error; err != nil {
		Utils.JSON(g, ogs.RspBase(ogs.StatusDestroyFailed, ogs.ErrorMessage("Destroy Failed")))
	} else {
		Utils.JSON(g, ogs.RspOK(ogs.SuccessMessage("Destroy Successfully")))
	}
}

func (c *handler) UpdateAccountPassword(g *gin.Context) {
	account, err := c.loadAccount(g)
	if err != nil {
		return
	}

	params, err := Utils.RequestBodyParams(g)
	if err != nil {
		Utils.JSON(g, ogs.RspBase(ogs.StatusSystemError, ogs.ErrorMessage("Invalid Request")))
		return
	}

	newPassword, _ := params["new_password"].(string)
	oldPassword, _ := params["old_password"].(string)

	if err = bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(oldPassword)); err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword { // Password does not match!
			errMsg := ogs.CodeText(ogs.StatusErrorPassword)
			Utils.JSON(g, ogs.RspBase(ogs.StatusErrorPassword, ogs.ErrorMessage(errMsg)))
			return
		}

		Utils.JSON(g, ogs.RspBase(ogs.StatusSystemError, ogs.ErrorMessage(err.Error())))
		return
	}

	account.Password = newPassword
	if err = account.UpdatePassword(); err != nil {
		Utils.JSON(g, ogs.RspBase(ogs.StatusUpdateFailed, ogs.ErrorMessage(err.Error())))
	} else {
		Utils.JSON(g, ogs.RspOK(ogs.SuccessMessage("Update Password Successfully")))
	}
}

func (c *handler) loadAccount(g *gin.Context) (account Account, err error) {
	id, _ := strconv.Atoi(g.Param("id"))
	account = Finder.FindAccountById(uint(id))

	if !account.IsPersisted() {
		Utils.JSON(g, ogs.RspBase(ogs.StatusUserNotFound, ogs.ErrorMessage("Account Not Found")))
		return account, gorm.ErrRecordNotFound
	}

	return account, nil
}

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

	Utils.JSON(g, ogs.RspOKWithPaginate("", accounts, paginate))
}

func (c *handler) FetchAccount(g *gin.Context) {
	account, err := c.loadAccount(g)
	if err != nil {
		return
	}

	Utils.JSON(g, ogs.RspOKWithData("", account))
}

func (c *handler) FetchCurrentAccountInfo(g *gin.Context) {
	headerToken := Utils.HeaderToken(g)
	account := Finder.FindAccountByToken(headerToken)
	if !account.IsPersisted() {
		Utils.JSON(g, ogs.RspError(ogs.StatusUserNotFound, "Get Current Account Info Failed"))
		return
	}

	Utils.JSON(g, ogs.RspOKWithData("", account))
}

func (c *handler) CreateAccount(g *gin.Context) {
	temp := tempAccount{}
	if err := json.NewDecoder(g.Request.Body).Decode(&temp); err != nil {
		Utils.JSON(g, ogs.RspError(ogs.StatusSystemError, "Invalid Request"))
		return
	}

	account := Account{
		Email:    temp.Email,
		Username: temp.Username,
		Password: temp.Password,
		Nickname: temp.Nickname,
		Avatar:   temp.Avatar,
	}

	err := account.Create()
	if err != nil {
		Utils.JSON(g, ogs.RspError(ogs.StatusCreateFailed, err.Error()))
		return
	}

	Utils.JSON(g, ogs.RspOKWithData("Create Successfully", account))
}

func (c *handler) UpdateAccount(g *gin.Context) {
	account, err := c.loadAccount(g)
	if err != nil {
		return
	}

	temp := tempAccount{}
	if err = json.NewDecoder(g.Request.Body).Decode(&temp); err != nil {
		Utils.JSON(g, ogs.RspError(ogs.StatusBadParams, "Bad Params"))
		return
	}

	account.Email = temp.Email
	account.Username = temp.Username
	account.Nickname = temp.Nickname
	account.Avatar = temp.Avatar

	if err = account.Update(); err != nil {
		Utils.JSON(g, ogs.RspError(ogs.StatusUpdateFailed, err.Error()))
	} else {
		Utils.JSON(g, ogs.RspOK("Update Successfully"))
	}
}

func (c *handler) DestroyAccount(g *gin.Context) {
	account, err := c.loadAccount(g)
	if err != nil {
		return
	}

	if err = DB.Delete(&account).Error; err != nil {
		Utils.JSON(g, ogs.RspError(ogs.StatusDestroyFailed, "Destroy Failed"))
	} else {
		Utils.JSON(g, ogs.RspOK("Destroy Successfully"))
	}
}

func (c *handler) UpdateAccountPassword(g *gin.Context) {
	account, err := c.loadAccount(g)
	if err != nil {
		return
	}

	params, err := Utils.RequestBodyParams(g)
	if err != nil {
		Utils.JSON(g, ogs.RspError(ogs.StatusBadParams, "Bad Params"))
		return
	}

	newPassword, _ := params["new_password"].(string)
	oldPassword, _ := params["old_password"].(string)

	if err = bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(oldPassword)); err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword { // Password does not match!
			errMsg := ogs.CodeText(ogs.StatusErrorPassword)
			Utils.JSON(g, ogs.RspError(ogs.StatusErrorPassword, errMsg))
			return
		}

		Utils.JSON(g, ogs.RspError(ogs.StatusSystemError, err.Error()))
		return
	}

	account.Password = newPassword
	if err = account.UpdatePassword(); err != nil {
		Utils.JSON(g, ogs.RspError(ogs.StatusUpdateFailed, err.Error()))
	} else {
		Utils.JSON(g, ogs.RspOK("Update Password Successfully"))
	}
}

func (c *handler) loadAccount(g *gin.Context) (account Account, err error) {
	id, _ := strconv.Atoi(g.Param("id"))
	account = Finder.FindAccountById(uint(id))

	if !account.IsPersisted() {
		Utils.JSON(g, ogs.RspError(ogs.StatusUserNotFound, "Account Not Found"))
		return account, gorm.ErrRecordNotFound
	}

	return account, nil
}

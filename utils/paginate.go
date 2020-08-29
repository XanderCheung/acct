package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/xandercheung/ogs-go"
	"strconv"
)

var defaultPerPage = 20

func Paginate(db *gorm.DB, page, perPage int) (*gorm.DB, ogs.BasePaginate) {
	if page <= 0 {
		page = 1
	}

	if perPage == 0 {
		perPage = defaultPerPage
	}

	totalCount := 0
	db.Count(&totalCount)

	offset := perPage * (page - 1)
	db = db.Limit(perPage).Offset(offset)

	return db, ogs.NewPaginate(page, totalCount, perPage)
}

func PaginateGin(db *gorm.DB, c *gin.Context) (*gorm.DB, ogs.BasePaginate) {
	page, _ := strconv.ParseInt(c.Query("page"), 10, 64)
	perPage, _ := strconv.ParseInt(c.Query("per_page"), 10, 64)
	return Paginate(db, int(page), int(perPage))
}

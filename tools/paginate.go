package tools

import (
	"github.com/gin-gonic/gin"
	"github.com/xandercheung/ogs-go"
	"gorm.io/gorm"
	"strconv"
)

var defaultPerPage = 20

// Paginate paginate query
func (t *Tool) Paginate(db *gorm.DB, page, perPage int) (*gorm.DB, ogs.BasePaginate) {
	if page <= 0 {
		page = 1
	}

	if perPage <= 0 {
		perPage = defaultPerPage
	}

	var totalCount int64 = 0
	db.Count(&totalCount)

	offset := perPage * (page - 1)
	db = db.Limit(perPage).Offset(offset)

	return db, ogs.NewPaginate(page, int(totalCount), perPage)
}

// PaginateGin paginate query
func (t *Tool) PaginateGin(db *gorm.DB, g *gin.Context) (*gorm.DB, ogs.BasePaginate) {
	page, _ := strconv.Atoi(g.Query("page"))
	perPage, _ := strconv.Atoi(g.Query("per_page"))
	return t.Paginate(db, page, perPage)
}

package asi

import (
	"github.com/jinzhu/gorm"
	"time"
)

var DB *gorm.DB

// InitDBConnection initialize a new db connection, need to import driver first, e.g:
// InitDBConnection("mysql", "user:password@tcp(host:port)/dbname?&charset=utf8mb4&parseTime=True&loc=Local")
func InitDBConnection(dialect string, args ...interface{}) (db *gorm.DB, err error) {
	db, err = gorm.Open(dialect, args)
	if err != nil {
		return
	}

	gorm.NowFunc = func() time.Time {
		return time.Now().UTC()
	}

	db.BlockGlobalUpdate(true)
	db.DB().SetMaxIdleConns(0)
	db.DB().SetMaxOpenConns(500)

	DB = db

	return
}
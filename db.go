package acct

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"os"
	"time"
)

var DB *gorm.DB

// InitDBConnection initialize a new db connection, need to import driver first, e.g:
// InitDBConnection("mysql", "user:password@tcp(host:port)/dbname?&charset=utf8mb4&parseTime=True&loc=UTC")
func InitDBConnection(dialect string, args ...interface{}) (db *gorm.DB, err error) {
	db, err = gorm.Open(dialect, args...)
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

func GetMysqlConnectArgsFromEnv() string {
	host := os.Getenv("MYSQL_HOST")
	port := os.Getenv("MYSQL_PORT")
	user := os.Getenv("MYSQL_USER")
	database := os.Getenv("MYSQL_DATABASE")
	password := os.Getenv("MYSQL_PASSWORD")
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?&charset=utf8mb4&parseTime=True&loc=UTC",
		user, password, host, port, database)
}

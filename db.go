package acct

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
)

var DB *gorm.DB

// InitDBConnection initialize a new db connection, need to import driver first, e.g:
// InitDBConnection("user:password@tcp(host:port)/dbname?&charset=utf8mb4&parseTime=True&loc=UTC")
// refer https://github.com/go-sql-driver/mysql#dsn-data-source-name for details
func InitDBConnection(dsn string) (err error) {
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	return err
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

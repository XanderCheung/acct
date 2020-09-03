package acct

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
)

var DB *gorm.DB

// InitDBAndSettings initialize a new db connection, load tables and create data of seeds
func InitDBAndSettings(config *Config) error {
	if config == nil {
		config = DefaultConfig
	}

	if config.IsLoadDSNFromENV {
		if err := Utils.LoadEnv(); err != nil {
			return err
		}

		config.ConnectionDSN = getMysqlConnectArgsFromEnv()
	}

	if err := connectDB(config.ConnectionDSN); err != nil {
		return err
	}

	if err := migrateTables(); err != nil {
		return err
	}

	if err := migrateSeeds(); err != nil {
		return err
	}

	config.load()

	return nil
}

// ConnectDB initialize a new db connection, need to import driver first, e.g:
// ConnectDB("user:password@tcp(host:port)/dbname?&charset=utf8mb4&parseTime=True&loc=UTC")
// refer https://github.com/go-sql-driver/mysql#dsn-data-source-name for details
func connectDB(dsn string) (err error) {
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	return err
}

func getMysqlConnectArgsFromEnv() string {
	host := os.Getenv("MYSQL_HOST")
	port := os.Getenv("MYSQL_PORT")
	user := os.Getenv("MYSQL_USER")
	database := os.Getenv("MYSQL_DATABASE")
	password := os.Getenv("MYSQL_PASSWORD")
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?&charset=utf8mb4&parseTime=True&loc=UTC",
		user, password, host, port, database)
}

package main

import (
	"github.com/xandercheung/acct"
	"github.com/xandercheung/acct/api"
	"github.com/xandercheung/acct/utils"
)

func main() {
	utils.InitEnv()

	connectArgs := acct.GetMysqlConnectArgsFromEnv()
	_, err := acct.InitDBConnection("mysql", connectArgs)
	if err != nil {
		panic(err)
	}

	acct.MigrateTables()
	acct.DBSeed()

	api.RunHttpServer()
}

package main

import (
	"github.com/gin-gonic/gin"
	"github.com/xandercheung/acct"
	"io"
	"os"
)

func main() {
	f, _ := os.OpenFile("log/acct-http.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	gin.DefaultWriter = io.MultiWriter(f)

	if err := acct.InitDBAndSettings(nil); err != nil {
		panic(err)
	}

	if err := acct.RunHttpServer(); err != nil {
		panic(err)
	}
}

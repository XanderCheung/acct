# acct
Service interface for user authentication

```shell script
[GIN-debug] POST   /api/v1/sign_in           --> github.com/xandercheung/acct.(*handler).SignIn-fm (3 handlers)
[GIN-debug] POST   /api/v1/sign_up           --> github.com/xandercheung/acct.(*handler).SignUp-fm (3 handlers)
[GIN-debug] GET    /api/v1/accounts/         --> github.com/xandercheung/acct.(*handler).FetchAccounts-fm (4 handlers)
[GIN-debug] POST   /api/v1/accounts/         --> github.com/xandercheung/acct.(*handler).CreateAccount-fm (4 handlers)
[GIN-debug] GET    /api/v1/accounts/:id      --> github.com/xandercheung/acct.(*handler).FetchAccount-fm (4 handlers)
[GIN-debug] POST   /api/v1/accounts/:id      --> github.com/xandercheung/acct.(*handler).UpdateAccount-fm (4 handlers)
[GIN-debug] DELETE /api/v1/accounts/:id      --> github.com/xandercheung/acct.(*handler).DestroyAccount-fm (4 handlers)
[GIN-debug] POST   /api/v1/accounts/:id/password --> github.com/xandercheung/acct.(*handler).UpdateAccountPassword-fm (4 handlers)
[GIN-debug] Listening and serving HTTP on :2337
```

### As a stand-alone API server
```shell script
go run cmd/acct.go
```

### As a package
```go
import "github.com/xandercheung/acct"

if err := acct.InitDBAndSettings(nil); err != nil {
    panic(err)
}

if err := acct.RunHttpServer(); err != nil {
    panic(err)
}
```

### Only use acct router(gin router)
```go
import "github.com/xandercheung/acct"
import "github.com/gin-gonic/gin"

router := gin.Default()
acct.SetRouter(router)
```

### Only use handler
```go
import "github.com/xandercheung/acct"
import "github.com/gin-gonic/gin"

r := gin.Default()
{
    r.POST("/sign_in", acct.Handler.SignIn)
    r.POST("/sign_up", acct.Handler.SignUp)

    // Token Authentication
    r.Use(acct.TokenAuthMiddleware())
    r.GET("/accounts", acct.Handler.FetchAccounts)
}

```
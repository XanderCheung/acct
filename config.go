package acct

import "os"

type Config struct {
	IsLoadDSNFromENV bool
	ConnectionDSN    string

	IsLoadJwtTokenKeyFromENV bool
	JwtTokenKey              string

	HttpServerPort string
}

var jwtTokenKey = ""
var httpServerPort = ""

// DefaultConfig default config of acct
var DefaultConfig = &Config{
	IsLoadDSNFromENV:         true,
	ConnectionDSN:            "",
	IsLoadJwtTokenKeyFromENV: true,
	JwtTokenKey:              "",
	HttpServerPort:           "2337",
}

func (c *Config) load() {
	if c.IsLoadJwtTokenKeyFromENV {
		jwtTokenKey = os.Getenv("JWT_TOKEN_KEY")
	} else {
		jwtTokenKey = c.JwtTokenKey
	}
	httpServerPort = c.HttpServerPort
}

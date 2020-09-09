package acct

import (
	"os"
)

type Config struct {
	IsLoadDSNFromENV bool
	ConnectionDSN    string

	IsLoadJwtTokenKeyFromENV bool
	JwtTokenKey              string

	HttpServerPort string
}

var JwtTokenKey = ""
var HttpServerPort = "2337"

// DefaultConfig default config of acct
var DefaultConfig = &Config{
	IsLoadDSNFromENV:         true,
	ConnectionDSN:            "",
	IsLoadJwtTokenKeyFromENV: true,
	JwtTokenKey:              "",
	HttpServerPort:           "2337",
}

func (c *Config) Load() {
	if c.IsLoadJwtTokenKeyFromENV {
		JwtTokenKey = os.Getenv("JWT_TOKEN_KEY")
	} else {
		JwtTokenKey = c.JwtTokenKey
	}
	HttpServerPort = c.HttpServerPort
}

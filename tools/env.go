package tools

import (
	"github.com/joho/godotenv"
)

// LoadEnv load env from .env
func (t *Tool) LoadEnv() error {
	return godotenv.Load(".env")
}

package utils

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"path/filepath"
)

func init() {
	err := godotenv.Load(filepath.Join(os.Getenv("GOPATH"), "Transactio", ".env"))
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	GwServiceAddr = os.Getenv("GW_ADDR")
	AuthServiceAddr = os.Getenv("AUTH_ADDR")

	GwLog = os.Getenv("GW_LOG")

	AppEnv = os.Getenv("APP_ENV")
}

var (
	AuthServiceAddr string
	GwServiceAddr   string

	GwLog string

	AppEnv string
)

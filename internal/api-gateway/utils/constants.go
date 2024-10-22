package utils

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

func init() {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	AuthServiceAddr = os.Getenv("AUTH_ADDR")

	GwLog = os.Getenv("GW_LOG")
	GwServiceAddr = os.Getenv("GW_ADDR")

	AppEnv = os.Getenv("APP_ENV")
}

var (
	AuthServiceAddr string
	GwServiceAddr   string

	AppEnv string
	GwLog  string
)

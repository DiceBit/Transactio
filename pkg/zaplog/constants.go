package zaplog

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

	AppEnv = os.Getenv("APP_ENV")
}

var (
	AppEnv string
)

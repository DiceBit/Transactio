package zaplog

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

	AppEnv = os.Getenv("APP_ENV")
}

var (
	AppEnv string
)

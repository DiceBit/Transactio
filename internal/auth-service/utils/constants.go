package utils

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

func init() {
	err := godotenv.Load(filepath.Join(os.Getenv("GOPATH"), "Transactio", ".env"))
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	AuthName = os.Getenv("AUTH_SERV")
	AuthServiceAddr = os.Getenv("AUTH_ADDR")

	AuthLog = os.Getenv("AUTH_LOG")

	jwtPrivateKey = []byte(os.Getenv("JWT_SECRET"))
	expTimeH, _ = strconv.Atoi(os.Getenv("EXP_TIME_HOUR"))
}

var (
	AuthName        string
	AuthServiceAddr string
	AuthLog         string

	jwtPrivateKey []byte
	expTimeH      int
)

var (
	UserRole  = "user"
	AdminRole = "admin"
)

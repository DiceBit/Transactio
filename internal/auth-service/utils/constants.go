package utils

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
)

func init() {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	AuthName = os.Getenv("AUTH_SERV")
	AuthServiceAddr = os.Getenv("AUTH_ADDR")

	AppEnv = os.Getenv("APP_ENV")
	AuthLog = os.Getenv("AUTH_LOG")

	jwtPrivateKey = []byte(os.Getenv("JWT_SECRET"))
	expTimeH, _ = strconv.Atoi(os.Getenv("EXP_TIME_HOUR"))
}

var (
	AuthName        string
	AuthServiceAddr string

	AppEnv  string
	AuthLog string

	jwtPrivateKey []byte
	expTimeH      int
)

var (
	UserRole  = "user"
	AdminRole = "admin"
)

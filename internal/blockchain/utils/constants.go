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

	BcAddr = os.Getenv("BC_ADDR")
	BcName = os.Getenv("BC_SERV")
	BCLog = os.Getenv("BC_LOG")

}

var (
	BcAddr string
	BcName string
	BCLog  string
)

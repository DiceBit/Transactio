package dbConn

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

	Host = os.Getenv("HOST_POSTGRES")
	Port, _ = strconv.Atoi(os.Getenv("PORT_POSTGRES"))
	User = os.Getenv("USER_POSTGRES")
	Pass = os.Getenv("PASSWORD_POSTGRES")
	Name = os.Getenv("DBNAME_POSTGRES")

}

var (
	Host string
	Port int
	User string
	Pass string
	Name string
)

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

	BcAddr = os.Getenv("BC_ADDR")
	BcName = os.Getenv("BC_SERV")
	BCLog = os.Getenv("BC_LOG")

	FsAddr = os.Getenv("FS_ADDR")

	MongoDbName = os.Getenv("DBNAME_MONGO")
	MongoCollections = os.Getenv("COLLECTIONS_NAME")
}

var (
	BcAddr string
	BcName string
	BCLog  string

	FsAddr string

	MongoDbName      string
	MongoCollections string
)

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

	FsAddr = os.Getenv("FS_ADDR")
	FsName = os.Getenv("FS_SERV")
	FsLog = os.Getenv("FS_LOG")
	IpfsAddr = os.Getenv("IPFS_ADDR")

	DefaultKey = os.Getenv("FILE_ENCRYPTING_KEY")
}

var (
	FsAddr   string
	FsName   string
	FsLog    string
	IpfsAddr string

	DefaultKey string
)

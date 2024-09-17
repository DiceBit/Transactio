package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
	"user-service/pkg/gRPC/server"
)

func init() {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	authName = os.Getenv("AUTH_SERV")
	authAddr = os.Getenv("AUTH_ADDR")
	authPort = os.Getenv("AUTH_PORT")
	authServiceAddr = fmt.Sprintf("%s:%s", authAddr, authPort)
}

var (
	authName        string
	authAddr        string
	authPort        string
	authServiceAddr string
)

func main() {
	//ip := os.Getenv("AUTH_ADDR")
	//port := os.Getenv("AUTH_PORT")
	//servName := os.Getenv("AUTH_SERV")
	//log.Printf("Server %s started on port %s\n", servName, port)
	//web.New().Run(ip + ":" + port)

	server.RunServe(authName, authServiceAddr)
}

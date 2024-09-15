package main

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"user-service/pkg/web"
)

func init() {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	ip := os.Getenv("IP_ADDR")
	port := os.Getenv("PORT")
	servName := os.Getenv("AUTH_SERV")
	log.Printf("Server %s started on port %s\n", servName, port)
	web.New().Run(ip + ":" + port)
}

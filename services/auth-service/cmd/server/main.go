package main

import (
	"github.com/joho/godotenv"
	"log"
	"user-service/pkg/web"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	web.New().Run("localhost:8080")
}

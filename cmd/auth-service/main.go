package main

import (
	"Transactio/internal/auth-service"
	"Transactio/internal/auth-service/utils"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	srv, err := auth_service.NewAuthService(utils.AuthName, utils.AuthServiceAddr)
	if err != nil {
		log.Fatal(err)
	}
	go srv.RunServe()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop
	srv.StopServe()
	log.Println("Gracefully stopped")
}

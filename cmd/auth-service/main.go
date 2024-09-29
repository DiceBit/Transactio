package main

import (
	"Transactio/internal/auth-service/server"
	"Transactio/internal/auth-service/utils"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	srv := server.NewAuthService(utils.AuthName, utils.AuthServiceAddr)
	go srv.RunServe()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop
	srv.StopServe()
	log.Println("Gracefully stopped")
}

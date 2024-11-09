package main

import (
	"Transactio/internal/fileStorage"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	srv := fileStorage.NewFsSrv()
	go srv.RunServer()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop
	srv.StopServer()
	log.Println("Gracefully stopped")
}

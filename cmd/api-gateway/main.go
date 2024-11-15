package main

import (
	gateway "Transactio/internal/api-gateway"
	"Transactio/internal/api-gateway/utils"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	gw := gateway.NewGW(utils.GwServiceAddr)
	go gw.RunServer()

	<-stop
	gw.Stop()
	log.Println("Gracefully stopped")
}

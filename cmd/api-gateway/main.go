package main

import (
	gateway "Transactio/internal/gateway/handlers"
	"Transactio/internal/gateway/utils"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	gw := gateway.NewGW(utils.GwServiceAddr)
	go gw.Start()

	<-stop
	gw.Stop()
	log.Println("Gracefully stopped")
}

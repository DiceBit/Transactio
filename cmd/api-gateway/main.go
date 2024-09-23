package main

//
//import (
//	"Transactio/services/api-gateway/gateway"
//	"Transactio/services/api-gateway/utils"
//	"os"
//	"os/signal"
//	"syscall"
//)
//
//func main() {
//
//	stop := make(chan os.Signal, 1)
//	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
//
//	gw := gateway.NewGW(utils.GwServiceAddr)
//	go gw.Start()
//
//	<-stop
//	gw.Stop()
//}

package main

import (
	"Transactio/internal/blockchain"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	srv, err := blockchain.NewBcSrv()
	if err != nil {
		log.Fatal(err)
	}
	go srv.RunServer()

	go testHTTP(srv)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop
	srv.StopServer()
	log.Println("Gracefully stopped")
}

func testHTTP(srv *blockchain.BcSrv) {
	router := mux.NewRouter()

	router.HandleFunc("/", zxc)
	router.HandleFunc("/file", srv.TestGetFile).Methods(http.MethodPost)
	http.ListenAndServe(":8089", router)
}
func zxc(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("<html><body><h1>Hello, World!</h1></body></html>"))
}

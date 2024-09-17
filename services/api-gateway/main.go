package main

import (
	authService "Ecommers/services/api-gateway/gRPC/proto"
	"Ecommers/services/api-gateway/middleware"
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net/http"
	"os"
)

func init() {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	authAddr = os.Getenv("AUTH_ADDR")
	authPort = os.Getenv("AUTH_PORT")
	authServiceAddr = fmt.Sprintf("%s:%s", authAddr, authPort)

	gwAddr = os.Getenv("GW_ADDR")
	gwPort = os.Getenv("GW_PORT")
	gwServiceAddr = fmt.Sprintf(":%s", gwPort)
}

var (
	authAddr        string
	authPort        string
	authServiceAddr string

	gwAddr        string
	gwPort        string
	gwServiceAddr string
)

func main() {
	Gateway(gwServiceAddr)
}

func Gateway(gatewayAddr string) {
	//Connect to Auth service
	grpcAuthServiceConn, err := grpc.NewClient(authServiceAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("Error with connecting to the %s. %s\n", os.Getenv("AUTH_SERV"), err)
	}
	defer grpcAuthServiceConn.Close()

	authClient := authService.NewAuthServiceClient(grpcAuthServiceConn)

	//REST-to-GRPC
	grpcGwMux := runtime.NewServeMux()

	err = authService.RegisterAuthServiceHandler(
		context.Background(),
		grpcGwMux,
		grpcAuthServiceConn,
	)
	if err != nil {
		log.Printf("Fail to start Gateway HTTP server. %s", err)
	}

	//REST
	router := mux.NewRouter()
	router.Use(middleware.AuthMiddleware(authClient))
	router.PathPrefix("/auth/").Handler(grpcGwMux).Methods(http.MethodPost)

	router.HandleFunc("/test", testHttp)

	log.Printf("Server Gateway started on %s\n", gwServiceAddr)
	log.Fatal(http.ListenAndServe(gatewayAddr, router))
}

func testHttp(w http.ResponseWriter, req *http.Request) {
	log.Println("test http") //todo delete
	fmt.Fprintln(w, "URL:", req.URL.String())
}

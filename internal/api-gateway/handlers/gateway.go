package handlers

import (
	"Transactio/internal/api-gateway/middleware"
	"Transactio/internal/api-gateway/services"
	"Transactio/internal/api-gateway/utils"
	"Transactio/pkg/zaplog"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/zap"
	"net/http"
)

type GW struct {
	addr    string
	authSrv *services.AuthServ

	router     *mux.Router
	grpcRouter *runtime.ServeMux

	logger *zap.Logger
}

func gatewayConf(gateway *GW) {
	authSrv := gateway.authSrv
	router := gateway.router
	logger := gateway.logger

	//AUTH-SERVICE
	err := services.RegisterAuthSrvHandler(authSrv.Conn, gateway.grpcRouter)
	if err != nil {
		logger.Error("Fail to start Gateway HTTP server", zap.Error(err))
	}

	router.Use(middleware.LoggerMiddleware(logger))
	router.Use(middleware.AuthMiddleware(authSrv.Client))
	router.PathPrefix("/auth/").Handler(gateway.grpcRouter).Methods(http.MethodPost)

	protectedAdmin := router.PathPrefix("/protected").Subrouter()
	protectedAdmin.Use(middleware.CheckRole([]string{"admin"}))
	protectedAdmin.HandleFunc("/testA", testA)

	protectedUsr := router.PathPrefix("/protected").Subrouter()
	protectedUsr.Use(middleware.CheckRole([]string{"user"}))
	protectedUsr.HandleFunc("/testU", testU)
}

func NewGW(gwAddr string) *GW {
	//REST
	router := mux.NewRouter()

	//REST-to-GRPC
	grpcMux := runtime.NewServeMux(runtime.WithMetadata(middleware.MetaDataForGW))

	logger := zaplog.NewLogger(utils.GwLog)

	authSrv, err := services.NewAuthServ(utils.AuthServiceAddr)
	if err != nil {
		logger.Error("Error creating auth service", zap.Error(err))
	}

	gwSrv := GW{
		addr:    gwAddr,
		authSrv: authSrv,

		router:     router,
		grpcRouter: grpcMux,

		logger: logger,
	}
	return &gwSrv
}

func (gw *GW) Start() {
	logger := gw.logger
	gatewayConf(gw)

	logger.Info(fmt.Sprintf("Server Gateway started on %s", utils.GwServiceAddr))
	err := http.ListenAndServe(gw.addr, gw.router)
	if err != nil {
		logger.Fatal("Error with listening and serve Gateway server", zap.Error(err))
	}
}
func (gw *GW) Stop() {
	gw.logger.Info("Gateway is stopped")

	gw.authSrv.Conn.Close()
	_ = gw.logger.Sync()
}

// ----
// Html test
// ----
func testA(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(w, "Hello, Admin. ", req.URL.String())
}
func testU(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(w, "Hello, User. ", req.URL.String())
}

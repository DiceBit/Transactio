package server

import (
	"auth-service/pkg/db/pgx"
	pb "auth-service/pkg/gRPC/proto"
	"auth-service/pkg/utils"
	"fmt"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"log"
	"net"
)

func (authServ *AuthServiceServer) RunServe() {
	authName := authServ.authName
	authServiceAddr := authServ.authAddr
	logger := authServ.logger

	listen, err := net.Listen("tcp", authServiceAddr)
	if err != nil {
		logger.Fatal("Error with starting", zap.String("servName", authName), zap.Error(err))
	}

	logger.Info(fmt.Sprintf("%s is running on %s", authName, authServiceAddr))

	pb.RegisterAuthServiceServer(authServ.grpcServer, authServ)

	err = authServ.grpcServer.Serve(listen)
	if err != nil {
		logger.Fatal("Error with serve", zap.Error(err))
	}
}

func (authServ *AuthServiceServer) StopServe() {
	msg := fmt.Sprintf("%s is stopped", authServ.authName)
	authServ.logger.Info(msg)

	authServ.db.Close()
	authServ.logger.Sync()
	authServ.grpcServer.GracefulStop()
}

func NewAuthService(authName, authAddr string) *AuthServiceServer {
	logger := utils.NewLogger(utils.AuthLog)

	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(LogInterceptor(logger)))

	db, err := pgx.New(logger)
	if err != nil {
		log.Println("Error with DB", err)
	}

	authServ := AuthServiceServer{
		db:         db,
		logger:     logger,
		authName:   authName,
		authAddr:   authAddr,
		grpcServer: grpcServer,
	}
	return &authServ
}

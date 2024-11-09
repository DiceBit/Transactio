package services

import (
	authService "Transactio/internal/api-gateway/gRPC/authProto"
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type AuthServ struct {
	Addr   string
	Conn   *grpc.ClientConn
	Client authService.AuthServiceClient
}

func NewAuthServ(addr string) (*AuthServ, error) {
	conn, err := authConn(addr)
	if err != nil {
		return nil, err
	}
	client := authClient(addr)

	srv := AuthServ{
		Addr:   addr,
		Conn:   conn,
		Client: client,
	}

	return &srv, nil
}

func RegisterAuthSrvHandler(authConn *grpc.ClientConn, mux *runtime.ServeMux) error {
	err := authService.RegisterAuthServiceHandler(
		context.Background(),
		mux,
		authConn,
	)
	return err
}

func authConn(authAddr string) (*grpc.ClientConn, error) {
	grpcAuthServiceConn, err := grpc.NewClient(
		authAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}
	return grpcAuthServiceConn, nil
}

func authClient(authAddr string) authService.AuthServiceClient {
	conn, _ := authConn(authAddr)
	client := authService.NewAuthServiceClient(conn)
	return client
}

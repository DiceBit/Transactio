package services

import (
	authService "Transactio/internal/api-gateway/gRPC/proto"
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
	conn, err := NewAuthConn(addr)
	if err != nil {
		return nil, err
	}
	client := NewAuthClient(addr)

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

func NewAuthConn(authAddr string) (*grpc.ClientConn, error) {
	grpcAuthServiceConn, err := grpc.NewClient(
		authAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}
	return grpcAuthServiceConn, nil
}

func NewAuthClient(authAddr string) authService.AuthServiceClient {
	conn, _ := NewAuthConn(authAddr)
	client := authService.NewAuthServiceClient(conn)
	return client
}

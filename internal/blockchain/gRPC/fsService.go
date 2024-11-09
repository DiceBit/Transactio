package gRPC

import (
	"Transactio/internal/blockchain/gRPC/fsProto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type FileStorageServ struct {
	Addr   string
	Conn   *grpc.ClientConn
	Client fsProto.FileStorageClient
}

func NewFSServ(addr string) (*FileStorageServ, error) {

	conn, err := fsConn(addr)
	if err != nil {
		return nil, err
	}
	client := fsClient(addr)

	srv := &FileStorageServ{
		Addr:   addr,
		Conn:   conn,
		Client: client,
	}
	return srv, nil
}

func fsConn(addr string) (*grpc.ClientConn, error) {
	client, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return client, nil
}

func fsClient(addr string) fsProto.FileStorageClient {
	conn, _ := fsConn(addr)
	client := fsProto.NewFileStorageClient(conn)
	return client
}

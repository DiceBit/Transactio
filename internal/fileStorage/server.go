package fileStorage

import (
	"Transactio/internal/fileStorage/gRPC/fsProto"
	"Transactio/internal/fileStorage/utils"
	"Transactio/pkg/zaplog"
	"bytes"
	"context"
	"fmt"
	shell "github.com/ipfs/go-ipfs-api"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"net"
)

type FsServer struct {
	fsProto.UnimplementedFileStorageServer
	grpcServer *grpc.Server

	shell  *shell.Shell
	logger *zap.Logger

	FsName string
	FsAddr string
}

func NewFsSrv() *FsServer {
	logger := zaplog.NewLogger(utils.FsLog)
	sh := shell.NewShell(utils.IpfsAddr)

	grpcSrv := grpc.NewServer()

	srv := FsServer{
		grpcServer: grpcSrv,
		shell:      sh,
		logger:     logger,
		FsName:     utils.FsName,
		FsAddr:     utils.FsAddr,
	}

	return &srv
}
func (srv *FsServer) RunServer() {
	fsName := srv.FsName
	fsAddr := srv.FsAddr
	log := srv.logger

	listen, err := net.Listen("tcp", fsAddr)
	if err != nil {
		log.Fatal("Error with starting", zap.String("servName", fsName), zap.Error(err))
	}

	log.Info(fmt.Sprintf("%s is running on %s", srv.FsName, srv.FsAddr))

	fsProto.RegisterFileStorageServer(srv.grpcServer, srv)

	err = srv.grpcServer.Serve(listen)
	if err != nil {
		log.Fatal("Error with serve", zap.Error(err))
	}

}
func (srv *FsServer) StopServer() {
	msg := fmt.Sprintf("%s is stopped", srv.FsName)
	srv.logger.Info(msg)

	_ = srv.logger.Sync()
	srv.grpcServer.GracefulStop()
}

func (srv *FsServer) AddFile(ctx context.Context, req *fsProto.AddFileRequest) (response *fsProto.AddFileResponse, err error) {
	sh := srv.shell
	log := srv.logger

	file, err := encrypted(req.File, req.Password, req.IsSecured)
	if err != nil {
		log.Error("Error when encrypting file", zap.Error(err))
		return nil, status.Errorf(codes.InvalidArgument, "Error when encrypting file. %v", err)
	}

	cid, err := sh.Add(file)
	if err != nil {
		log.Error("Error when adding file", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "Error when adding file. %v", err)
	}

	fmt.Printf("\n %v \n", cid)

	err = sh.Pin(cid)
	if err != nil {
		log.Error("Error when pinning file", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "Error when pinning file. %v", err)
	}

	log.Info("File added")
	return &fsProto.AddFileResponse{Cid: cid}, nil
}
func (srv *FsServer) GetFile(ctx context.Context, req *fsProto.GetFileRequest) (response *fsProto.GetFileResponse, err error) {
	sh := srv.shell
	log := srv.logger

	r, err := sh.Cat(req.Cid)
	if err != nil {
		log.Error("Error when getting data from IPFS", zap.Error(err))
		return nil, status.Errorf(codes.NotFound, "Error when getting data from IPFS. %v", err)
	}

	reader, err := decrypted(r, req.Password, req.IsSecured)
	if err != nil {
		log.Error("Error when decrypting file", zap.Error(err))
		return nil, status.Errorf(codes.DataLoss, "Error when decrypting file. %v", err)
	}

	return &fsProto.GetFileResponse{FileReader: reader}, nil
}
func (srv *FsServer) RemoveFile(ctx context.Context, req *fsProto.RemoveFileRequest) (response *fsProto.RemoveFileResponse, err error) {
	sh := srv.shell
	log := srv.logger

	err = sh.Unpin(req.Cid)
	if err != nil {
		log.Error("Error when unpinning(deleting) file in IPFS", zap.Error(err))
		return &fsProto.RemoveFileResponse{IsDeleted: false},
			status.Errorf(codes.Internal, "Error when unpinning(deleting) file in IPFS. %v", err)
	}

	log.Info("File completely remove")
	return &fsProto.RemoveFileResponse{IsDeleted: true}, nil
}

func encrypted(f []byte, password string, isSecured bool) (reader *bytes.Reader, err error) {

	var encryptedFile []byte

	if !isSecured {
		encryptedFile, err = utils.EncryptData(f, utils.DefaultKey)
		if err != nil {
			return nil, err
		}
	} else {
		encryptedFile, err = utils.EncryptData(f, password)
		if err != nil {
			return nil, err
		}
	}

	reader = bytes.NewReader(encryptedFile)
	return reader, nil
}
func decrypted(r io.Reader, password string, isSecured bool) ([]byte, error) {
	bytesReader, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	var decryptedFile []byte
	if !isSecured {
		decryptedFile, err = utils.DecryptData(bytesReader, utils.DefaultKey)
		if err != nil {
			return nil, err
		}
	} else {
		decryptedFile, err = utils.DecryptData(bytesReader, password)
		if err != nil {
			return nil, err
		}
	}

	return decryptedFile, nil
}

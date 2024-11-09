package blockchain

import (
	"Transactio/internal/blockchain/db"
	"Transactio/internal/blockchain/gRPC"
	"Transactio/internal/blockchain/gRPC/fsProto"
	"Transactio/internal/blockchain/models"
	"Transactio/internal/blockchain/utils"
	mongodb "Transactio/pkg/dbConn/mongo"
	"Transactio/pkg/dbConn/pgx"
	"Transactio/pkg/zaplog"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"log"
	"mime/multipart"
	"net/http"
	"time"
)

type BcSrv struct {
	fsSrv  *gRPC.FileStorageServ
	db     *pgxpool.Pool
	Mongo  *mongo.Client
	logger *zap.Logger

	BcName string
	BcAddr string
}

func NewBcSrv() (*BcSrv, error) {

	logger := zaplog.NewLogger(utils.BCLog)

	bcDb, err := pgx.New()
	if err != nil {
		logger.Error("Error with DB", zap.Error(err))
		return nil, err
	}
	logger.Info("pgx DB connected")

	monDb, err := mongodb.New()
	if err != nil {
		logger.Error("Error with DB", zap.Error(err))
		return nil, err
	}
	logger.Info("MongoDB connected")
	err = db.CreateIndex(context.Background(), monDb)
	if err != nil {
		return nil, err
	}

	fsServ, err := gRPC.NewFSServ(utils.FsAddr)
	if err != nil {
		return nil, err
	}

	srv := BcSrv{
		fsSrv: fsServ,

		db:     bcDb,
		Mongo:  monDb,
		logger: logger,

		BcName: utils.BcName,
		BcAddr: utils.BcAddr,
	}

	srv.setGenesisBlock(context.Background())
	return &srv, nil
}
func (srv *BcSrv) RunServer() {
	srv.logger.Info(fmt.Sprintf("%s is running on %s", srv.BcName, srv.BcAddr))
}
func (srv *BcSrv) StopServer() {
	srv.logger.Info(fmt.Sprintf("%s is stopped", srv.BcName))

	srv.fsSrv.Conn.Close()
	srv.db.Close()
	_ = srv.logger.Sync()
	_ = srv.Mongo.Disconnect(context.Background())
}
func (srv *BcSrv) setGenesisBlock(ctx context.Context) {
	log := srv.logger

	exist, err := db.CheckGenBlock(ctx, srv.db)
	if err != nil {
		log.Error("Error with check genesis block")
		return
	}
	if !exist {
		genesisBlock := models.SetGenesisBlock()
		if _, err = db.AddBlock(ctx, srv.db, genesisBlock); err != nil {
			log.Error("Error with adding genesis block", zap.Error(err))
			return
		}
		log.Info("Genesis Block is installed")
	} else {
		log.Info("Genesis Block was already installed")
	}
}

func (srv *BcSrv) TestGetFile(w http.ResponseWriter, req *http.Request) {
	file, handler, err := req.FormFile("file")
	if err != nil {
		log.Printf("\n %v \n", err)
		return
	}
	if err := srv.SaveFile(context.Background(), file, handler, "testOwner", "zxc", true); err != nil {
		http.Error(w, "Save file error", http.StatusInternalServerError)
		return
	}
	w.Write([]byte("<html><body><h1>File uploaded successfully</h1></body></html>"))
}

/*FUNC FOR GATEWAY SERVICE*/
func (srv *BcSrv) SaveFile(ctx context.Context, file multipart.File, fileInfo *multipart.FileHeader, owner, password string, isSecured bool) error {
	fsSrv := srv.fsSrv.Client

	fileBytes, err := utils.ConvertMultipartToBytes(file)
	if err != nil {
		srv.logger.Error("Failed to convert multipart to bytes.", zap.Error(err))
		return err
	}

	fileResponse, err := fsSrv.AddFile(ctx, &fsProto.AddFileRequest{
		File:      fileBytes,
		Password:  password,
		IsSecured: isSecured,
	})
	if err != nil {
		srv.logger.Error("Error from file-storage service", zap.Error(err))
		return err
	}

	prevBlockHash, err := db.PrevHash(ctx, srv.db)
	if err != nil {
		srv.logger.Error("Error with getting prevBlockHash", zap.Error(err))
		return err
	}

	cid := fileResponse.GetCid()
	data := models.NewData(cid, owner, fileInfo.Filename, int(fileInfo.Size), false, isSecured)
	newBlock := &models.Blockchain{
		PrevBlockHash: prevBlockHash,
		Timestamp:     time.Now().Unix(),
		Fmd:           data,
	}
	newBlock.SetHash()

	index, err := db.AddBlock(ctx, srv.db, newBlock)
	if err != nil {
		srv.logger.Error("Error with adding block", zap.Error(err))
		return err
	}

	if err = db.InsertInfo(ctx, srv.Mongo, owner, fileInfo.Filename, index); err != nil {
		srv.logger.Error("Error when insert data in mongo", zap.Error(err))
		return err
	}

	srv.logger.Info(fmt.Sprintf("Block(%s) has been added", cid[:25]))
	return nil
}
func (srv *BcSrv) DeleteFile(ctx context.Context, cid string) error {
	return nil
}
func (srv *BcSrv) GetFile() {

}

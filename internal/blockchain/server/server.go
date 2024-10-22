package server

import (
	pgxDb "Transactio/internal/blockchain/db"
	"Transactio/internal/blockchain/models"
	utils2 "Transactio/internal/blockchain/utils"
	"Transactio/pkg/dbConn/pgx"
	"Transactio/pkg/zaplog"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"time"
)

type BcSrv struct {
	bc     *models.Blockchain
	db     *pgxpool.Pool
	logger *zap.Logger

	BcName string
	BcAddr string
}

func NewBcSrv() *BcSrv {
	logger := zaplog.NewLogger(utils2.BCLog)

	bcDb, err := pgx.New(logger)

	if err != nil {
		logger.Error("Error with DB", zap.Error(err))
		return nil
	}

	name := utils2.BcName
	addr := utils2.BcAddr

	srv := BcSrv{
		db:     bcDb,
		logger: logger,

		BcName: name,
		BcAddr: addr,
	}

	srv.setGenesisBlock()

	return &srv
}

func (srv *BcSrv) AddBlock(ctx context.Context,
	cid, owner, fileName string, fileSize, version int, status bool) {
	data := models.NewData(cid, owner, fileName, fileSize, version, status)

	prevBlockHash, err := pgxDb.PrevHash(context.Background(), srv.db)
	if err != nil {
		srv.logger.Error("Error with getting prevBlockHash", zap.Error(err))
		return
	}

	newBlock := &models.Blockchain{
		Fmd:           data,
		Hash:          "",
		PrevBlockHash: prevBlockHash,
		Timestamp:     time.Now().Unix(),
	}
	newBlock.SetHash()

	err = pgxDb.AddBlock(context.Background(), srv.db, newBlock)
	if err != nil {
		srv.logger.Error("Error with adding block", zap.Error(err))
		return
	}

	var shortVal = cid
	if len(cid) > 25 {
		shortVal = cid[:25]
	}
	msg := fmt.Sprintf("Block(%s) has been added", shortVal)
	srv.logger.Info(msg)
}

func (srv *BcSrv) RunServer() {

	//TODO: GRPC connections
	srv.logger.Info(fmt.Sprintf("%s is running on %s", srv.BcName, srv.BcAddr))
}
func (srv *BcSrv) StopServer() {
	msg := fmt.Sprintf("%s is stopped", srv.BcName)
	srv.logger.Info(msg)

	srv.db.Close()
	_ = srv.logger.Sync()

	//TODO: Server graceful stop
	//authServ.grpcServer.GracefulStop()
}

func (srv *BcSrv) setGenesisBlock() {
	exist, err := pgxDb.CheckGenBlock(context.Background(), srv.db)
	if err != nil {
		srv.logger.Error("Error with check genesis block")
		return
	} else if !exist {
		genesisBlock := models.SetGenesisBlock()
		err = pgxDb.AddBlock(context.Background(), srv.db, genesisBlock)
		if err != nil {
			srv.logger.Error("Error with adding genesis block", zap.Error(err))
			return
		}
	} else {
		srv.logger.Info("Genesis Block was already installed")
		return
	}
	srv.logger.Info("Genesis Block is installed")
}

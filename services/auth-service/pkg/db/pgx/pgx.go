package pgx

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"os"
	"strconv"
)

type DB struct {
	pool *pgxpool.Pool
}

func New(logger *zap.Logger) (*pgxpool.Pool, error) {
	port, err := strconv.Atoi(os.Getenv("PORT_POSTGRES"))
	if err != nil {
		logger.Error("Error convert string to integer", zap.Error(err))
		return nil, err
	}

	pgxInfo := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("HOST_POSTGRES"),
		port,
		os.Getenv("USER_POSTGRES"),
		os.Getenv("PASSWORD_POSTGRES"),
		os.Getenv("DBNAME_POSTGRES"))

	pgxConn, err := pgxpool.New(context.Background(), pgxInfo)
	if err != nil {
		logger.Error("Error with pgxpool connect", zap.Error(err))
		return nil, err
	}

	err = pgxConn.Ping(context.Background())
	if err != nil {
		return nil, err
	}
	logger.Info("DB connected")
	return pgxConn, nil
}

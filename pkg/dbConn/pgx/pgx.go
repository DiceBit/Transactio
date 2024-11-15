package pgx

import (
	"Transactio/pkg/dbConn"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	pool *pgxpool.Pool
}

func New() (*pgxpool.Pool, error) {
	pgxInfo := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		dbConn.Host, dbConn.Port, dbConn.User, dbConn.Pass, dbConn.Name)

	pgxConn, err := pgxpool.New(context.Background(), pgxInfo)
	if err != nil {
		return nil, err
	}

	err = pgxConn.Ping(context.Background())
	if err != nil {
		return nil, err
	}

	return pgxConn, nil
}

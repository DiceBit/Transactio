package pgx

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"os"
	"strconv"
)

type DB struct {
	pool *pgxpool.Pool
}

func New() (*pgxpool.Pool, error) {
	port, err := strconv.Atoi(os.Getenv("PORT_POSTGRES"))
	if err != nil {
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
		return nil, err
	}

	err = pgxConn.Ping(context.Background())
	if err != nil {
		return nil, err
	}
	log.Println("DB connected")

	return pgxConn, nil
}

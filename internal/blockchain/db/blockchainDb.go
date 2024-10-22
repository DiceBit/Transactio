package pgxDb

import (
	"Transactio/internal/blockchain/models"
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

func AddBlock(ctx context.Context, db *pgxpool.Pool, block *models.Blockchain) error {
	tx, err := db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	fmd := block.Fmd
	var batch = &pgx.Batch{}
	batch.Queue(`insert into filemd(cid, owneraddr, filename, filesize, createat, status, version) 
			VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		fmd.Cid, fmd.OwnerAddr, fmd.FileName, fmd.FileSize, time.Unix(fmd.CreateAt, 0), fmd.Status, fmd.Version)
	batch.Queue(`insert into blockchain(hash, prevblockhash, timestamp) 
		VALUES ($1, $2, $3)`,
		block.Hash, block.PrevBlockHash, time.Unix(block.Timestamp, 0))

	res := tx.SendBatch(ctx, batch)
	err = res.Close()
	if err != nil {
		return err
	}

	tx.Commit(ctx)
	return nil
}

func PrevHash(ctx context.Context, db *pgxpool.Pool) (string, error) {
	tx, err := db.Begin(ctx)
	if err != nil {
		return "", err
	}
	defer tx.Rollback(ctx)

	rows, err := tx.Query(ctx, `select hash from blockchain where index=
                                  (select max(index) from blockchain)`)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	var bc models.Blockchain
	if !rows.Next() {
		return "", errors.New("prevHash not found")
	}
	err = rows.Scan(
		&bc.Hash,
	)
	if err != nil {
		return "", err
	}

	tx.Commit(ctx)
	return bc.Hash, nil
}

func CheckGenBlock(ctx context.Context, db *pgxpool.Pool) (bool, error) {
	var exist bool
	err := db.QueryRow(ctx, `select exists(select 1 from blockchain where index=1)`).Scan(&exist)
	if err != nil {
		return false, err
	}

	return exist, nil
}

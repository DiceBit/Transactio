package db

import (
	"Transactio/internal/blockchain/models"
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

func AddBlock(ctx context.Context, db *pgxpool.Pool, block *models.Blockchain) (int, error) {
	tx, err := db.Begin(ctx)
	if err != nil {
		return -1, err
	}
	defer tx.Rollback(ctx)

	fmd := block.Fmd
	var batch = &pgx.Batch{}
	batch.Queue(`insert into filemd(cid, owneraddr, filename, filesize, isdelete, issecured)
			VALUES ($1, $2, $3, $4, $5, $6)`,
		fmd.Cid, fmd.OwnerAddr, fmd.FileName, fmd.FileSize, fmd.IsDelete, fmd.IsSecured)
	batch.Queue(`insert into blockchain(hash, prevblockhash, timestamp)
		VALUES ($1, $2, $3) RETURNING index`,
		block.Hash, block.PrevBlockHash, time.Unix(block.Timestamp, 0))

	res := tx.SendBatch(ctx, batch)

	_, err = res.Exec()
	if err != nil {
		return -1, err
	}

	var index int
	err = res.QueryRow().Scan(&index)
	if err != nil {
		return -1, err
	}

	err = res.Close()
	if err != nil {
		return 0, err
	}

	tx.Commit(ctx)
	return index, nil
}

func ReadBlock(ctx context.Context, db *pgxpool.Pool, index int) (models.FileMD, error) {
	var md models.FileMD

	err := db.QueryRow(ctx, `select cid, owneraddr, filename, filesize, isdelete, issecured from filemd where id=$1`, index).Scan(
		&md.Cid,
		&md.OwnerAddr,
		&md.FileName,
		&md.FileSize,
		&md.IsDelete,
		&md.IsSecured,
	)
	if err != nil {
		return models.FileMD{}, err
	}

	return md, nil
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

package pgxDb

import (
	pb "Transactio/internal/auth-service/gRPC/authProto"
	"Transactio/internal/auth-service/models"
	"Transactio/internal/auth-service/utils"
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
	"time"
)

func AddUser(ctx context.Context, db *pgxpool.Pool, usrClaims *pb.SignUpRequest) error {
	tx, err := db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(usrClaims.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	usrClaims.Password = string(hashedPassword)
	usrBalance := 0
	usrRoles := []string{utils.UserRole}

	var batch = &pgx.Batch{}
	batch.Queue(`insert into Users(username, email, password, balance, role, createat) values ($1,$2,$3,$4,$5, $6)`,
		usrClaims.Username,
		usrClaims.Email,
		usrClaims.Password,
		usrBalance,
		usrRoles,
		time.Now())

	res := tx.SendBatch(ctx, batch)

	err = res.Close()
	if err != nil {
		return err
	}

	tx.Commit(ctx)
	return nil
}

func UsrByEmail(ctx context.Context, db *pgxpool.Pool, email string) (*models.User, error) {
	tx, err := db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	rows, err := tx.Query(ctx, `select email, password, role from users where email=$1`, email)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var usr models.User
	if !rows.Next() {
		return nil, errors.New("user not found")
	}
	err = rows.Scan(
		&usr.Email,
		&usr.Password,
		&usr.Role,
	)
	if err != nil {
		return nil, err
	}
	tx.Commit(ctx)
	return &usr, nil
}

func CheckIfExistUsr(ctx context.Context, db *pgxpool.Pool, usrClaims *pb.SignUpRequest) (bool, error) {
	var exist bool
	err := db.QueryRow(ctx,
		`select exists(select 1 from users where email=$1 or username=$2)`,
		usrClaims.Email, usrClaims.Username).Scan(&exist)
	if err != nil {
		return false, err
	}

	return exist, nil
}

package userUtils

import (
	"context"
	pgx2 "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
	"log"
	"time"
	pb "user-service/pkg/gRPC/proto"
	"user-service/pkg/models"
	"user-service/pkg/utils"
)

func AddUser(ctx context.Context, db *pgxpool.Pool, usrClaims *pb.SignUpRequest) error {
	tx, err := db.Begin(ctx)
	if err != nil {
		log.Println("Error with transaction:", err)
		return err
	}
	defer tx.Rollback(ctx)

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(usrClaims.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Error hashing password:", err)
		return err
	}

	usrClaims.Password = string(hashedPassword)
	usrBalance := 0
	usrRoles := []string{utils.UserRole}

	var batch = &pgx2.Batch{}
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
		log.Println("Error with transaction:", err)
		return err
	}

	tx.Commit(ctx)
	return nil
}

func UsrByEmail(ctx context.Context, db *pgxpool.Pool, email string) (*models.User, error) {
	tx, err := db.Begin(ctx)
	if err != nil {
		log.Println("Error with transaction:", err)
		return nil, err
	}
	defer tx.Rollback(ctx)

	rows, err := tx.Query(ctx, `select email, password, role from users where email=$1`, email)
	if err != nil {
		log.Println("Error with query:", err)
		return nil, err
	}
	defer rows.Close()

	var usr models.User
	if !rows.Next() {
		log.Println("User not found:", err)
		return nil, err
	}

	err = rows.Scan(
		&usr.Email,
		&usr.Password,
		&usr.Role,
	)
	if err != nil {
		log.Println("Error with  scanning query:", err)
		return nil, err
	}
	tx.Commit(ctx)
	return &usr, nil
}

func CheckIfExistUsr(ctx context.Context, db *pgxpool.Pool, usrClaims *pb.SignUpRequest) (bool, error) {
	tx, err := db.Begin(ctx)
	if err != nil {
		log.Println("Error with transaction:", err)
		return false, err
	}
	defer tx.Rollback(ctx)

	_, err = db.Query(ctx, `select * from users where email=$1 or username=$2`, usrClaims.Email, usrClaims.Username)
	if err == nil {
		tx.Commit(ctx)
		return false, nil
	} else {
		tx.Commit(ctx)
		return true, nil
	}
}

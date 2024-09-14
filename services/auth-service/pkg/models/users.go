package models

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type User struct {
	Id       int
	Username string
	Email    string
	Password string
	Balance  int
	Role     []string

	CreateAt pgtype.Timestamp
}

const (
	user  = iota
	admin = iota
)

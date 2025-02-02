// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package db

import (
	"database/sql"
)

type Machine struct {
	ID      int32
	Name    string
	BuyerID sql.NullInt32
	OwnerID int32
	Ram     int32
	Cpu     int32
	Memory  int32
	Key     sql.NullString
	Host    string
	SshUser string
}

type User struct {
	ID       int32
	Name     string
	Email    string
	Password []byte
}

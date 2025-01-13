package database

import (
	"database/sql"
	"fmt"
	db "machine-marketplace/internal/DB/generated"

	_ "github.com/lib/pq"
)

var (
	DB      *sql.DB
	Queries *db.Queries
)

func Init() error {
	var err error
	DB, err = sql.Open("postgres", "postgresql://postgres:postgres@localhost:5432/machine_market")
	if err != nil {
		return err
	}

	if err := DB.Ping(); err != nil {
		return err
	}

	fmt.Println("database connect successfully")

	Queries = db.New(DB)
	return nil
}

func Close() {
	if DB != nil {
		DB.Close()
	}
}

package database

import (
	"database/sql"
	"fmt"
	"log/slog"
	db "machine-marketplace/internal/DB/generated"
	"os"
	"sync"

	_ "github.com/lib/pq"
)

type (
	Config struct {
		Host     string `short:"h" long:"db-host" description:"Database host" default:"postgres"`
		Port     string `short:"P" long:"db-port" description:"Database port" default:"5432"`
		User     string `short:"u" long:"db-user" description:"Database user" default:"postgres"`
		Password string `short:"w" long:"db-password" description:"Database password" default:"postgres"`
		Name     string `short:"d" long:"db-name" description:"Database name" default:"machine_market"`
		SSLMode  string `short:"s" long:"db-sslmode" description:"Database SSL mode" default:"disable"`
	}
)

var (
	DB            *sql.DB
	Queries       *db.Queries
	once          sync.Once
	closeOnce     sync.Once
	initErr       error
	isInitialized bool
	config        *Config
	l             = slog.New(slog.NewTextHandler(os.Stdout, nil))
)

func Init() (error, *db.Queries) {
	return InitWithConfig(&Config{
		Host:     "postgres",
		Port:     "5432",
		User:     "postgres",
		Password: "postgres",
		Name:     "machine_market",
		SSLMode:  "disable",
	})
}

func InitWithConfig(cfg *Config) (error, *db.Queries) {
	var queries *db.Queries
	once.Do(func() {
		config = cfg
		initErr, queries = initDatabase()
		if initErr == nil {
			isInitialized = true
			Queries = queries
		}
	})
	return initErr, queries
}

func IsInitialized() bool {
	return isInitialized
}

func initDatabase() (error, *db.Queries) {
	l.Info("Init - initializing database connection",
		"host", config.Host,
		"port", config.Port,
		"user", config.User,
		"database", config.Name)

	connectionString := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=%s",
		config.User, config.Password, config.Host, config.Port, config.Name, config.SSLMode)

	var err error
	DB, err = sql.Open("postgres", connectionString)
	if err != nil {
		l.Error("Init - failed to open database connection", "error", err)
		return err, nil
	}

	if err := DB.Ping(); err != nil {
		l.Error("Init - failed to ping database", "error", err)
		return err, nil
	}

	l.Info("Init - database connection established successfully")

	queries := db.New(DB)

	// Setup database schema and seed data
	if err := SetupDatabase(); err != nil {
		l.Error("Init - failed to setup database schema", "error", err)
		return err, nil
	}

	l.Info("Init - database initialization completed successfully")
	return nil, queries
}

// func Init() error {
// 	var err error
// 	DB, err = sql.Open("postgres", "postgresql://postgres:postgres@localhost:5432/machine_market?sslmode=disable")
// 	if err != nil {
// 		return err
// 	}

// 	if err := DB.Ping(); err != nil {
// 		return err
// 	}

// 	fmt.Println("database connect successfully")

// 	Queries = db.New(DB)
// 	return nil
// }

func Close() error {
	var closeErr error
	closeOnce.Do(func() {
		if DB != nil && isInitialized {
			l.Info("Close - closing database connection")
			closeErr = DB.Close()
			if closeErr != nil {
				l.Error("Close - failed to close database connection", "error", closeErr)
			} else {
				l.Info("Close - database connection closed successfully")
				isInitialized = false
			}
		}
	})
	return closeErr
}

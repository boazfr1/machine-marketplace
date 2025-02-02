package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	database "machine-marketplace/internal/DB/generated"
	"machine-marketplace/internal/data"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func SetupDatabase() error {
	schema := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255) UNIQUE NOT NULL,
		email VARCHAR(255) UNIQUE NOT NULL,
		password BYTEA NOT NULL
	);

	CREATE TABLE IF NOT EXISTS machines (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		buyer_id INTEGER REFERENCES users(id),
		owner_id INTEGER NOT NULL REFERENCES users(id),
		ram INTEGER NOT NULL,
		cpu INTEGER NOT NULL,
		memory INTEGER NOT NULL,
		key TEXT,
		host VARCHAR(255) NOT NULL,
		ssh_user VARCHAR(255) NOT NULL
	);`

	// Execute schema creation
	_, err := DB.Exec(schema)
	if err != nil {
		return fmt.Errorf("error creating schema: %v", err)
	}

	// Check if we need to insert mock data
	var count int
	err = DB.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
	if err != nil {
		return fmt.Errorf("error checking users count: %v", err)
	}

	if count == 0 {
		return insertMockData()
	}

	return nil
}

func insertMockData() error {
	ctx := context.Background()

	userParams := []data.UserParams{
		{
			Name:     "boaz",
			Email:    "boaz@test.com",
			Password: "password",
		},
		{
			Name:     "noam",
			Email:    "noam@test.com",
			Password: "password",
		},
	}

	for _, params := range userParams {
		err := insertUserMockData(params, ctx)
		if err != nil {
			return fmt.Errorf("error inserting user: %v", err)
		}
	}

	machines := []struct {
		name    string
		ram     int32
		cpu     int32
		memory  int32
		ownerId int32
		buyerId sql.NullInt32
	}{
		{
			name:    "DevMachine-1",
			ram:     8,
			cpu:     4,
			memory:  500,
			ownerId: 1,
			buyerId: sql.NullInt32{Int32: 2, Valid: true},
		},
		{
			name:    "GameServer-1",
			ram:     16,
			cpu:     8,
			memory:  1000,
			ownerId: 1,
			buyerId: sql.NullInt32{}, // no buyer
		},
		{
			name:    "WebHost-1",
			ram:     4,
			cpu:     2,
			memory:  250,
			ownerId: 2,
			buyerId: sql.NullInt32{Int32: 1, Valid: true},
		},
		{
			name:    "DataCruncher-1",
			ram:     32,
			cpu:     16,
			memory:  2000,
			ownerId: 2,
			buyerId: sql.NullInt32{}, // no buyer
		},
	}

	for _, m := range machines {
		_, err := Queries.CreateMachine(ctx, database.CreateMachineParams{
			Name:    m.name,
			Ram:     m.ram,
			Cpu:     m.cpu,
			Memory:  m.memory,
			Key:     sql.NullString{String: "mock-key", Valid: true},
			OwnerID: m.ownerId,
			Host:    "localhost",
			SshUser: "mockuser",
		})

		if err != nil {
			return fmt.Errorf("error inserting machine %s: %v", m.name, err)
		}
	}

	log.Println("Mock data inserted successfully")
	return nil
}

func insertUserMockData(params data.UserParams, ctx context.Context) error {
	cryptPassword, err := bcrypt.GenerateFromPassword([]byte(params.Password), 14)
	if err != nil {
		return err
	}

	createParams := database.CreateUserParams{
		Name:    params.Name,
		Email:   params.Email,
		Column3: cryptPassword,
	}

	_, err = Queries.CreateUser(ctx, createParams)
	if err != nil {
		if strings.Contains(err.Error(), "unique constraint") {
			return err
		}
		return err
	}
	return nil
}

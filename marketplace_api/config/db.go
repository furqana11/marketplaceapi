package config

import (
	"os"
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"

	controllers "github.com/FarisTF/marketplace_api/controllers"
)

// Connecting to db
func Connect() *pgxpool.Pool {
	databaseUrl :="postgres://postgres:rahasia!@localhost:5432/marketplace_api"
	dbPool, err := pgxpool.Connect(context.Background(), databaseUrl)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Connected to db")
	controllers.InitiateDB(dbPool)
	return dbPool
}
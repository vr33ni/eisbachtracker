package db

import (
	"context"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv" // for loading local .env
)

var Conn *pgxpool.Pool

func Init() error {
	// Only load .env locally
	if os.Getenv("ENV") != "production" {
		_ = godotenv.Load()
	}

	url := os.Getenv("DATABASE_URL")
	pool, err := pgxpool.New(context.Background(), url)
	if err != nil {
		return err
	}
	Conn = pool
	return nil
}

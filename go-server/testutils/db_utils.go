package testutils

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

func SetupTestDB(t *testing.T) *pgxpool.Pool {
	db, err := pgxpool.New(context.Background(), "postgres://vreeni@localhost:5432/eisbach")
	if err != nil {
		t.Fatalf("Failed to connect to DB: %v", err)
	}

	return db
}

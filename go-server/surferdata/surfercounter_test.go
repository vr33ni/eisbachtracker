package surferdata

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func setupService(t *testing.T) *Service {
	LoadTestConfig(t)

	db, err := pgxpool.New(context.Background(), "postgres://vreeni@localhost:5432/eisbach")
	if err != nil {
		t.Fatalf("DB connect fail: %v", err)
	}
	return NewService(db)
}

func TestAddAndGetEntries(t *testing.T) {
	service := setupService(t)

	err := service.AddEntry(5, time.Now())
	if err != nil {
		t.Fatalf("Failed to add entry: %v", err)
	}

	entries, err := service.GetAllEntries()
	if err != nil {
		t.Fatalf("Failed to fetch entries: %v", err)
	}

	if len(entries) == 0 {
		t.Fatal("Expected at least one entry")
	}

	t.Logf("Fetched entries: %+v", entries[0])
}

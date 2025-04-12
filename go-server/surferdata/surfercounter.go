package surferdata

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// SurferEntry represents a user-submitted count of surfers at a specific time.
type SurferEntry struct {
	Timestamp time.Time `json:"timestamp"`
	Count     int       `json:"count"`
}

// Service handles database operations for surfer data.
type Service struct {
	DB *pgxpool.Pool
}

// NewService returns a new instance of Service with the given DB pool connection.
func NewService(db *pgxpool.Pool) *Service {
	return &Service{DB: db}
}

// AddEntry inserts a new surfer entry into the database.
// If no timestamp is provided (zero time), it defaults to the current time.
func (s *Service) AddEntry(count int, when time.Time) error {
	if when.IsZero() {
		when = time.Now()
	}
	_, err := s.DB.Exec(context.Background(),
		`INSERT INTO surfer_entries (timestamp, count) VALUES ($1, $2)`,
		when, count)
	return err
}

// GetAllEntries retrieves all surfer entries, ordered from newest to oldest.
func (s *Service) GetAllEntries() ([]SurferEntry, error) {
	rows, err := s.DB.Query(context.Background(),
		`SELECT timestamp, count FROM surfer_entries ORDER BY timestamp DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []SurferEntry
	for rows.Next() {
		var entry SurferEntry
		if err := rows.Scan(&entry.Timestamp, &entry.Count); err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}
	return entries, nil
}

package surferdata

import (
	"context"

	"github.com/jackc/pgx/v5"
)

// PredictSurferCount returns an average surfer count for the given hour.
// Optionally, filter by temperature range if provided.
func (s *Service) PredictSurferCount(hour int, temperature *float64) (float64, error) {
	var rows pgx.Rows
	var err error

	if temperature != nil {
		tempRangeMin := *temperature - 2
		tempRangeMax := *temperature + 2
		query := `
			SELECT count FROM surfer_entries
			WHERE EXTRACT(HOUR FROM timestamp) = $1
			AND temperature BETWEEN $2 AND $3
		`
		rows, err = s.DB.Query(context.Background(), query, hour, tempRangeMin, tempRangeMax)
	} else {
		query := `
			SELECT count FROM surfer_entries
			WHERE EXTRACT(HOUR FROM timestamp) = $1
		`
		rows, err = s.DB.Query(context.Background(), query, hour)
	}

	if err != nil {
		return 0, err
	}
	defer rows.Close()

	var total int
	var count int
	for rows.Next() {
		var c int
		if err := rows.Scan(&c); err != nil {
			return 0, err
		}
		total += c
		count++
	}

	if count == 0 {
		return 0, nil // instead of returning an error - return nil (no data for this hour)
	}

	return float64(total) / float64(count), nil
}

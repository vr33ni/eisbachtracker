package surferdata

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/vr33ni/eisbachtracker-pwa/go-server/conditions"
)

// SurferEntry represents a user-submitted count of surfers at a specific time.
type SurferEntry struct {
	Timestamp        time.Time `json:"timestamp"`
	Count            int       `json:"count"`
	WaterTemperature *float64  `json:"water_temperature,omitempty"`
	AirTemperature   *float64  `json:"air_temperature,omitempty"`
	WeatherCondition *string   `json:"weather_condition,omitempty"`
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
// If no timestamp is provided (zero time), it defaults to the current time. Temperature is being added based on the latest temperature
func (s *Service) AddEntry(count int, when time.Time) error {
	if when.IsZero() {
		when = time.Now()
	}
	weather, err := conditions.GetCurrentWeather()
	if err != nil {
		log.Println("⚠️ Could not fetch air weather:", err)
		weather = &conditions.WeatherData{
			Temp:      0,
			Condition: "Unknown",
		}
	}

	waterTemp, err := conditions.GetCachedWaterTemperature()
	if err != nil {
		log.Println("⚠️ Could not fetch water temp:", err)
		waterTemp = 0
	}

	_, err = s.DB.Exec(context.Background(),
		`INSERT INTO surfer_entries (timestamp, count, water_temperature, air_temperature, weather_condition) 
	 VALUES ($1, $2, $3, $4, $5)`,
		when, count, waterTemp, weather.Temp, weather.Condition)

	return err
}

// GetAllEntries retrieves all surfer entries, ordered from newest to oldest.
func (s *Service) GetAllEntries() ([]SurferEntry, error) {
	rows, err := s.DB.Query(context.Background(),
		`SELECT timestamp, count, water_temperature, air_temperature, weather_condition FROM surfer_entries ORDER BY timestamp DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []SurferEntry
	for rows.Next() {
		var entry SurferEntry
		if err := rows.Scan(&entry.Timestamp, &entry.Count, &entry.WaterTemperature, &entry.AirTemperature, &entry.WeatherCondition); err != nil {
			return nil, err
		}

		entries = append(entries, entry)
	}
	return entries, nil
}

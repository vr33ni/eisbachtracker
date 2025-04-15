package surferdata

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vr33ni/eisbachtracker-pwa/go-server/conditions"
)

type SurferEntry struct {
	Timestamp        time.Time `json:"timestamp"`
	Count            int       `json:"count"`
	WaterTemperature *float64  `json:"water_temperature,omitempty"`
	AirTemperature   *float64  `json:"air_temperature,omitempty"`
	WeatherCondition *string   `json:"weather_condition,omitempty"`
	WaterLevel       *float64  `json:"water_level,omitempty"`
	WaterFlow        *float64  `json:"water_flow,omitempty"`
}

type SurferEntryResponse struct {
	Timestamp        time.Time `json:"timestamp"`
	Count            int       `json:"count"`
	WaterTemperature float64   `json:"water_temperature"`
	AirTemperature   float64   `json:"air_temperature"`
	WeatherCondition string    `json:"weather_condition"`
	WaterLevel       float64   `json:"water_level"`
	WaterFlow        float64   `json:"water_flow"`
}

type Service struct {
	DB           *pgxpool.Pool
	WaterService conditions.WaterDataProvider // ✅ use the interface here
}

func NewService(db *pgxpool.Pool, ws *conditions.WaterService) *Service {
	return &Service{DB: db, WaterService: ws}
}

func (s *Service) AddEntry(count int, when time.Time, waterLevel *float64, waterFlow *float64) error {
	if when.IsZero() {
		when = time.Now()
	}

	weather, err := conditions.GetCurrentWeather()
	if err != nil {
		log.Println("⚠️ Could not fetch air weather:", err)
		weather = &conditions.WeatherData{Temp: 0, Condition: "Unknown"}
	}

	waterTemp, err := conditions.GetCachedWaterTemperature()
	if err != nil {
		log.Println("⚠️ Could not fetch water temp:", err)
		waterTemp = 0
	}

	if waterLevel == nil || waterFlow == nil {
		wl, wf, err := s.WaterService.GetCurrentWaterConditions()
		if err != nil {
			log.Println("⚠️ Could not fetch water level/flow:", err)
			wl, wf = 0, 0
		}
		waterLevel = &wl
		waterFlow = &wf
	}

	_, err = s.DB.Exec(context.Background(),
		`INSERT INTO surfer_entries (timestamp, count, water_temperature, air_temperature, weather_condition, water_level, water_flow) 
		 VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		when, count, waterTemp, weather.Temp, weather.Condition, waterLevel, waterFlow,
	)

	return err
}

func (s *Service) GetAllEntries() ([]SurferEntryResponse, error) {
	rows, err := s.DB.Query(context.Background(),
		`SELECT timestamp, count, water_temperature, air_temperature, weather_condition, water_level, water_flow 
		FROM surfer_entries ORDER BY timestamp DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []SurferEntryResponse
	for rows.Next() {
		var e SurferEntry
		if err := rows.Scan(&e.Timestamp, &e.Count, &e.WaterTemperature, &e.AirTemperature, &e.WeatherCondition, &e.WaterLevel, &e.WaterFlow); err != nil {
			return nil, err
		}

		entries = append(entries, SurferEntryResponse{
			Timestamp:        e.Timestamp,
			Count:            e.Count,
			WaterTemperature: safeFloat(e.WaterTemperature),
			AirTemperature:   safeFloat(e.AirTemperature),
			WeatherCondition: safeString(e.WeatherCondition),
			WaterLevel:       safeFloat(e.WaterLevel),
			WaterFlow:        safeFloat(e.WaterFlow),
		})
	}
	return entries, nil
}

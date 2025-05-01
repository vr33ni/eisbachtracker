package surferdata

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type MLPredictionParams struct {
	Hour             int     `json:"hour"`
	WaterTemp        float64 `json:"water_temp"`
	AirTemp          float64 `json:"air_temp"`
	WaterLevel       float64 `json:"water_level"`
	WeatherCondition string  `json:"weather_condition"`
}

type MLPredictionResponse struct {
	SurferCount int `json:"surfer_count"`
}

func (s *Service) PredictSurferCountML(params MLPredictionParams) (int, error) {
	// Prepare the request payload
	payload := map[string]interface{}{
		"hour":                     params.Hour,
		"water_temp":               params.WaterTemp,
		"air_temp":                 params.AirTemp,
		"water_level":              params.WaterLevel,
		"weather_condition_cloudy": 0,
		"weather_condition_rainy":  0,
		"weather_condition_snow":   0,
		"weather_condition_stormy": 0,
	}

	// Map the weather condition to the correct encoded field
	switch params.WeatherCondition {
	case "cloudy":
		payload["weather_condition_cloudy"] = 1
	case "rainy":
		payload["weather_condition_rainy"] = 1
	case "snow":
		payload["weather_condition_snow"] = 1
	case "stormy":
		payload["weather_condition_stormy"] = 1
	}

	// Convert payload to JSON
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return 0, fmt.Errorf("failed to marshal payload: %w", err)
	}

	url := os.Getenv("FLASK_API_URL")
	// Make the HTTP POST request to the Flask API
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		return 0, fmt.Errorf("failed to call ML prediction API: %w", err)
	}
	defer resp.Body.Close()

	// Check for non-200 status codes
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return 0, fmt.Errorf("ML prediction API returned status %d: %s", resp.StatusCode, string(body))
	}

	// Parse the response
	var mlResponse MLPredictionResponse
	if err := json.NewDecoder(resp.Body).Decode(&mlResponse); err != nil {
		return 0, fmt.Errorf("failed to decode ML prediction response: %w", err)
	}

	// Return the predicted surfer count
	return mlResponse.SurferCount, nil
}

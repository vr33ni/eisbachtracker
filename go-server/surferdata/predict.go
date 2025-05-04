package surferdata

import (
	"context"
	"fmt"
	"math"

	"github.com/vr33ni/eisbachtracker-pwa/go-server/conditions"
)

type PredictionParams struct {
	Hour             int
	WaterTemp        *float64
	AirTemp          *float64
	WeatherCondition int
	WaterLevel       float64
	WaterFlow        float64
}

// BasePredictionByHour fetches avg surfer count from DB for given hour
func (s *Service) basePredictionByHour(hour int) (float64, error) {
	var avg *float64
	err := s.DB.QueryRow(context.Background(),
		`SELECT AVG(count) FROM surfer_entries WHERE EXTRACT(HOUR FROM timestamp) = $1`,
		hour,
	).Scan(&avg)

	if err != nil {
		return 0, err
	}

	// fallback logic for weird hours (no data or tiny value)
	if avg == nil || *avg < 1 {
		// night hours fallback (basically no one)
		if hour >= 22 || hour <= 5 {
			return 0, nil // super low base
		}
		return 1, nil // minimal base for daytime
	}

	return *avg, nil
}

func (s *Service) PredictSurferCountAdvanced(params PredictionParams) (interface{}, error) {
	// Step 1: Get the base prediction by hour (rule-based fallback)
	base, err := s.basePredictionByHour(params.Hour)
	if err != nil {
		return 0, err
	}

	// Step 2: Calculate the rule-based factor
	weatherData := &conditions.WeatherData{
		Temp:      safeFloat(params.AirTemp),
		Condition: params.WeatherCondition,
	}
	factor := calculateFactor(params.Hour, params.WaterTemp, weatherData, params.WaterLevel, params.WaterFlow)
	ruleBasedPrediction := int(math.Round(base * factor))
	if ruleBasedPrediction < 0 {
		ruleBasedPrediction = 0
	}

	// Step 3: Call the ML-based prediction
	mlParams := MLPredictionParams{
		Hour:             params.Hour,
		WaterTemp:        safeFloat(params.WaterTemp),
		AirTemp:          safeFloat(params.AirTemp),
		WaterLevel:       params.WaterLevel,
		WeatherCondition: params.WeatherCondition,
	}
	mlPrediction, explanation, err := s.PredictSurferCountML(mlParams)
	if err != nil {
		fmt.Printf("Error predicting surfer count: %v\n", err)
		return 0, err
	}

	fmt.Printf("Predicted Surfer Count: %d\n", mlPrediction)
	fmt.Println("Feature Contributions:")
	for feature, contribution := range explanation {
		fmt.Printf("  %s: %.2f\n", feature, contribution)
	}

	// Step 4: Combine the predictions (optional)
	// Combine the response
	response := map[string]interface{}{
		"hour":              params.Hour,
		"water_temperature": safeFloat(params.WaterTemp),
		"air_temperature":   safeFloat(params.AirTemp),
		"weather_condition": params.WeatherCondition,
		"water_level":       params.WaterLevel,
		"prediction":        mlPrediction,
		"explanation":       explanation,
	}

	return response, nil
}

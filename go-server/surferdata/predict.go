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
	WeatherCondition string
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

func (s *Service) PredictSurferCountAdvanced(params PredictionParams) (int, error) {
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
	mlPrediction, err := s.PredictSurferCountML(mlParams)
	if err != nil {
		// Log the error and fall back to rule-based prediction
		fmt.Printf("ML prediction failed: %v\n", err)
		return ruleBasedPrediction, nil
	}

	// Step 4: Combine the predictions (optional)
	// For now, we return the ML prediction
	return mlPrediction, nil
}

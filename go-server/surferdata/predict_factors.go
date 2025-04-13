package surferdata

import (
	"log"

	"github.com/vr33ni/eisbachtracker-pwa/go-server/conditions"
	"github.com/vr33ni/eisbachtracker-pwa/go-server/config"
)

// calculateFactor applies all dynamic factors based on the current context
func calculateFactor(hour int, waterTemp *float64, weatherData *conditions.WeatherData) float64 {
	factor := 1.0

	log.Printf("Hour Factor: %+v", config.Predict.HourFactor)
	log.Printf("Water Temp Factor: %+v", config.Predict.WaterTempFactor)
	log.Printf("Air Temp Factor: %+v", config.Predict.AirTempFactor)
	log.Printf("Weather Condition Factor: %+v", config.Predict.WeatherConditionFactor)

	// --- Hour Factor ---
	if hour >= 18 || hour <= 8 {
		factor *= config.Predict.HourFactor.Offpeak
	} else {
		factor *= config.Predict.HourFactor.Peak
	}

	// --- Water Temp Factor ---
	if waterTemp != nil {
		factor *= waterTempFactor(*waterTemp)
	}

	// --- Air Temp Factor ---
	if weatherData != nil {
		factor *= airTempFactor(weatherData.Temp)

		// --- Weather Condition Factor ---
		if conditionFactor, ok := config.Predict.WeatherConditionFactor[weatherData.Condition]; ok {
			factor *= conditionFactor
		} else {
			factor *= config.Predict.WeatherConditionFactor["Unknown"]
		}
	}

	return factor
}

// --- Helpers ---

func waterTempFactor(temp float64) float64 {
	switch {
	case temp < 10:
		return config.Predict.WaterTempFactor.Cold
	case temp < 15:
		return config.Predict.WaterTempFactor.Medium
	default:
		return config.Predict.WaterTempFactor.Warm
	}
}

func airTempFactor(temp float64) float64 {
	switch {
	case temp < 10:
		return config.Predict.AirTempFactor.Cold
	case temp < 15:
		return config.Predict.AirTempFactor.Medium
	default:
		return config.Predict.AirTempFactor.Hot
	}
}

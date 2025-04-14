package surferdata

import (
	"github.com/vr33ni/eisbachtracker-pwa/go-server/conditions"
	"github.com/vr33ni/eisbachtracker-pwa/go-server/config"
)

// calculateFactor applies all dynamic factors based on the current context
func calculateFactor(hour int, waterTemp *float64, weatherData *conditions.WeatherData, waterLevel float64, waterFlow float64) float64 {
	factor := 1.0

	// Time of day influence
	if hour >= 6 && hour <= 8 {
		factor += 0.3
	} else if hour >= 12 && hour <= 14 {
		factor += 0.2
	} else if hour >= 22 || hour <= 5 {
		factor -= 0.4
	}

	// Water temperature influence
	if waterTemp != nil && *waterTemp < 10 {
		factor -= 0.2
	}

	// Weather influence
	if weatherData.Condition == "Rain" || weatherData.Condition == "Snow" {
		factor -= 0.3
	}

	// Water level influence
	if waterLevel < 140 {
		factor -= 0.3
	} else if waterLevel > 145 {
		factor += 0.2
	}

	// Water flow -> ignored for now, but kept for future
	_ = waterFlow // placeholder to avoid unused warning

	if factor < 0.5 {
		factor = 0.5
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

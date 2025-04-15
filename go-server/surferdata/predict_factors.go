package surferdata

import (
	"github.com/vr33ni/eisbachtracker-pwa/go-server/conditions"
)

// calculateFactor applies all dynamic factors based on the current context
func calculateFactor(
	hour int,
	waterTemp *float64,
	weatherData *conditions.WeatherData,
	waterLevel float64,
	waterFlow float64,
) float64 {
	factor := 1.0

	// 🕒 Time of day influence
	if hour >= 6 && hour <= 8 {
		factor += 0.3 // Early morning surf crowd
	} else if hour >= 12 && hour <= 14 {
		factor += 0.2 // Lunchtime bump
	} else if hour >= 22 || hour <= 5 {
		factor -= 0.4 // Night time drop
	}

	// ❄️ Water temperature influence
	if waterTemp != nil && *waterTemp < 10 {
		factor -= 0.2
	}

	// 🌧️ Weather influence
	if weatherData.Condition == "Rain" || weatherData.Condition == "Snow" {
		factor -= 0.3
	}

	// 🌊 Water level influence
	if waterLevel < 140 {
		factor -= 0.3
	} else if waterLevel > 145 {
		factor += 0.2
	}

	// 🏞️ (optional future): Water flow — currently unused
	_ = waterFlow

	// ✨ Safety cap
	if factor < 0.5 {
		factor = 0.5
	}

	return factor
}

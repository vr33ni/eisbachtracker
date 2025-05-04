package testutils

import "github.com/vr33ni/eisbachtracker-pwa/go-server/conditions"

func MockWeatherData(temp float64, condition int) *conditions.WeatherData {
	return &conditions.WeatherData{
		Temp:      temp,
		Condition: condition,
	}
}

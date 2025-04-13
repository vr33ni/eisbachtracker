package conditions

import (
	"encoding/json"
	"net/http"
)

func GetCurrentWeather() (*WeatherData, error) {
	url := "https://api.open-meteo.com/v1/forecast?latitude=48.137154&longitude=11.576124&current_weather=true"

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var apiResp struct {
		CurrentWeather struct {
			Temp        float64 `json:"temperature"`
			WeatherCode int     `json:"weathercode"`
		} `json:"current_weather"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, err
	}

	return &WeatherData{
		Temp:      apiResp.CurrentWeather.Temp,
		Condition: mapWeatherCode(apiResp.CurrentWeather.WeatherCode),
	}, nil
}

func mapWeatherCode(code int) string {
	switch {
	case code == 0:
		return "Clear"
	case code == 1 || code == 2 || code == 3:
		return "Cloudy"
	case code >= 45 && code < 60:
		return "Fog/Drizzle"
	case code >= 61 && code < 70:
		return "Rain"
	case code >= 71 && code < 80:
		return "Snow"
	case code >= 95:
		return "Thunderstorm"
	default:
		return "Unknown"
	}
}

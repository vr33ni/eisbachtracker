package conditions

import (
	"encoding/json"
	"net/http"
)

type AirService struct{}

func NewAirService() *AirService {
	return &AirService{}
}

type AirDataProvider interface {
	GetCurrentWeather() (*WeatherData, error)
}

func (ws *AirService) GetCurrentWeather() (*WeatherData, error) {
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
		Condition: apiResp.CurrentWeather.WeatherCode, // Use numeric WeatherCode directly
	}, nil
}

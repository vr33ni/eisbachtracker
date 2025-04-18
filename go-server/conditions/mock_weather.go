package conditions

func GetCurrentWeather() (*WeatherData, error) {
	// Mock data
	return &WeatherData{
		Temp:      22,
		Condition: "Clear",
	}, nil
}

func GetLatestWaterTemperature() (float64, error) {
	return 18.5, nil
}

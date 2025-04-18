package conditions

// MockWaterService is a mock implementation of WaterDataProvider for testing.
type MockWaterService struct{}

// GetLatestWaterLevelAndFlow implements WaterDataProvider.
func (m *MockWaterService) GetLatestWaterLevelAndFlow() (*WaterLevelAndFlow, error) {
	panic("unimplemented")
}

func (m *MockWaterService) GetCurrentWeather() (*WeatherData, error) {
	return &WeatherData{
		Temp:      20.0,
		Condition: "Clear",
	}, nil
}

func (m *MockWaterService) GetCachedWaterTemperature() (float64, error) {
	return 16.5, nil
}

func (m *MockWaterService) GetCurrentWaterConditions() (float64, float64, error) {
	return 143.0, 9.5, nil
}

func (m *MockWaterService) GetLatestWaterTemperature() (float64, error) {
	return 16.5, nil
}

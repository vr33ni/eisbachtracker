package conditions

import "time"

// MockWaterService is a mock implementation of WaterDataProvider for testing.
type MockWaterService struct{}

// GetLatestWaterLevelAndFlow implements WaterDataProvider.
func (m *MockWaterService) GetLatestWaterLevelAndFlow() (*WaterLevelAndFlow, error) {
	return &WaterLevelAndFlow{
		Level:       143.0,
		Flow:        9.5,
		RequestDate: time.Now().Format(time.RFC3339),
	}, nil
}

func (m *MockWaterService) GetCachedWaterTemperature() (float64, error) {
	return 16.5, nil
}

func (m *MockWaterService) GetLatestWaterTemperature() (float64, error) {
	return 16.5, nil
}

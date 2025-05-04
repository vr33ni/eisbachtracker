package surferdata

import (
	"testing"
	"time"

	"github.com/vr33ni/eisbachtracker-pwa/go-server/conditions"
)

type MockWaterService struct{}

// GetLatestWaterTemperature implements conditions.WaterDataProvider.
func (m *MockWaterService) GetLatestWaterTemperature() (float64, error) {
	panic("unimplemented")
}

func (m *MockWaterService) GetCachedWaterTemperature() (float64, error) {
	return 15.5, nil
}

func (m *MockWaterService) GetLatestWaterLevelAndFlow() (*conditions.WaterLevelAndFlow, error) {
	return &conditions.WaterLevelAndFlow{
		Level: 120.0,
		Flow:  20.5,
	}, nil
}

type MockAirService struct{}

func (m *MockAirService) GetCurrentWeather() (*conditions.WeatherData, error) {
	return &conditions.WeatherData{
		Temp:      22.3,
		Condition: 0,
	}, nil
}

func TestAddAndGetEntries(t *testing.T) {
	service := setupTestService(t)

	service.WaterService = &MockWaterService{}
	service.AirService = &MockAirService{}

	err := service.AddEntry(5, time.Now(), nil, nil, nil)
	if err != nil {
		t.Fatalf("Failed to add entry: %v", err)
	}

	entries, err := service.GetAllEntries()
	if err != nil {
		t.Fatalf("Failed to fetch entries: %v", err)
	}

	if len(entries) == 0 {
		t.Fatal("Expected at least one entry")
	}

	t.Logf("Fetched entry: %+v", entries[0])
}

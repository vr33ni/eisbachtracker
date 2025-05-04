package surferdata

import (
	"testing"

	"github.com/vr33ni/eisbachtracker-pwa/go-server/conditions"
	"github.com/vr33ni/eisbachtracker-pwa/go-server/testutils"
	"github.com/vr33ni/eisbachtracker-pwa/go-server/utils"
)

func TestCalculateFactorSunnyWarmPeak(t *testing.T) {
	testutils.LoadTestConfig(t)

	f := calculateFactor(
		14,                // hour
		utils.Float64(20), // water temp
		&conditions.WeatherData{Temp: 25, Condition: 0}, // clear
		146, // water level
		15,  // water flow (ignored for now)
	)

	t.Logf("factor: %.2f", f)

	if f <= 1.0 {
		t.Error("Expected factor to increase for sunny warm peak conditions")
	}
}

func TestCalculateFactorColdRainyOffpeak(t *testing.T) {
	testutils.LoadTestConfig(t)

	f := calculateFactor(
		6,                // hour
		utils.Float64(5), // water temp
		&conditions.WeatherData{Temp: 5, Condition: 61}, // rain
		135, // water level
		10,  // water flow (ignored for now)
	)

	t.Logf("factor: %.2f", f)

	if f >= 1.0 {
		t.Error("Expected factor to decrease for cold rainy off-peak conditions")
	}
}

func TestCalculateFactorLowWaterLevel(t *testing.T) {
	testutils.LoadTestConfig(t)

	f := calculateFactor(10, utils.Float64(15), &conditions.WeatherData{Temp: 15, Condition: 0}, 135, 10) // weather = clear

	t.Logf("factor: %.2f", f)

	if f >= 1.0 {
		t.Error("Expected factor to decrease for low water level")
	}
}

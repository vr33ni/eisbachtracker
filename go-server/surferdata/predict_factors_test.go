package surferdata

import (
	"testing"

	"github.com/vr33ni/eisbachtracker-pwa/go-server/testutils"
	"github.com/vr33ni/eisbachtracker-pwa/go-server/utils"

	"github.com/vr33ni/eisbachtracker-pwa/go-server/conditions"
)

func TestCalculateFactorSunnyWarmPeak(t *testing.T) {
	testutils.LoadTestConfig(t)

	f := calculateFactor(14, utils.Float64(20), &conditions.WeatherData{Temp: 25, Condition: "Clear"})
	t.Logf("factor: %.2f", f)

	if f <= 1.0 {
		t.Error("Expected factor to increase for sunny warm peak conditions")
	}
}

func TestCalculateFactorColdRainyOffpeak(t *testing.T) {
	testutils.LoadTestConfig(t)

	f := calculateFactor(6, utils.Float64(5), &conditions.WeatherData{Temp: 5, Condition: "Rain"})
	if f >= 1.0 {
		t.Error("Expected factor to decrease for cold rainy off-peak conditions")
	}
}

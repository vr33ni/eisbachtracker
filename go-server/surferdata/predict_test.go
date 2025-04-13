package surferdata

import (
	"testing"

	"github.com/vr33ni/eisbachtracker-pwa/go-server/utils"
)

func TestPredictSurferCount_HourOnly(t *testing.T) {
	service := setupTestService(t)
	pred, err := service.PredictSurferCountAdvanced(PredictionParams{
		Hour: 14,
	})
	if err != nil {
		t.Fatalf("Prediction failed: %v", err)
	}

	t.Logf("Prediction for hour=14 → %d", pred)
}

func TestPredictSurferCount_WithWaterTemp(t *testing.T) {
	service := setupTestService(t)

	pred, err := service.PredictSurferCountAdvanced(PredictionParams{
		Hour:      18,
		WaterTemp: utils.Float64(18),
	})
	if err != nil {
		t.Fatalf("Prediction failed: %v", err)
	}

	t.Logf("Prediction for hour=18 with 18°C water temp → %d", pred)
}

func TestPredictSurferCount_AllFactorsSunny(t *testing.T) {
	service := setupTestService(t)

	pred, err := service.PredictSurferCountAdvanced(PredictionParams{
		Hour:             14,
		WaterTemp:        utils.Float64(18),
		AirTemp:          utils.Float64(25),
		WeatherCondition: "Clear",
	})
	if err != nil {
		t.Fatalf("Prediction failed: %v", err)
	}

	t.Logf("Prediction for hour=14 sunny warm → %d", pred)
}

func TestPredictSurferCount_AllFactorsBad(t *testing.T) {
	service := setupTestService(t)

	pred, err := service.PredictSurferCountAdvanced(PredictionParams{
		Hour:             5,
		WaterTemp:        utils.Float64(4),
		AirTemp:          utils.Float64(2),
		WeatherCondition: "Rain",
	})
	if err != nil {
		t.Fatalf("Prediction failed: %v", err)
	}

	t.Logf("Prediction for hour=5 cold rainy → %d", pred)
}

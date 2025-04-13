package testutils

import (
	"os"
	"testing"

	"github.com/vr33ni/eisbachtracker-pwa/go-server/config"
)

func LoadTestConfig(t *testing.T) {
	os.Setenv("PREDICT_CONFIG", "../config/predict.toml")

	if err := config.LoadConfig(); err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}
}

package testutils

import (
	"os"
	"testing"

	"github.com/vr33ni/eisbachtracker-pwa/go-server/config"
)

func LoadTestConfig(t *testing.T) {
	pathsToTry := []string{
		"./config/predict.toml",  // normal
		"../config/predict.toml", // when running from subpackage
	}

	var loaded bool
	for _, path := range pathsToTry {
		if _, err := os.Stat(path); err == nil {
			os.Setenv("PREDICT_CONFIG", path)
			if err := config.LoadConfig(); err != nil {
				t.Fatalf("Failed to load config: %v", err)
			}
			loaded = true
			break
		}
	}

	if !loaded {
		t.Fatalf("Failed to find predict.toml in expected paths")
	}
}

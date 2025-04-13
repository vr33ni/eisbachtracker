package surferdata

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/vr33ni/eisbachtracker-pwa/go-server/config"
)

func LoadTestConfig(t *testing.T) {
	rootPath, _ := filepath.Abs(filepath.Join("..", "config", "predict.toml"))
	os.Setenv("PREDICT_CONFIG", rootPath)

	if err := config.LoadConfig(); err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}
}

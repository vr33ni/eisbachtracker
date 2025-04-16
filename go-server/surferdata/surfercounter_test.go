package surferdata

import (
	"testing"
	"time"

	"github.com/vr33ni/eisbachtracker-pwa/go-server/conditions"
)

func TestAddAndGetEntries(t *testing.T) {
	service := setupTestService(t)

	mockWS := &conditions.MockWaterService{}

	service.WaterService = mockWS

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

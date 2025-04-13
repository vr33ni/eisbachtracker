package surferdata

import (
	"testing"
	"time"
)

func TestAddAndGetEntries(t *testing.T) {
	service := setupTestService(t)

	err := service.AddEntry(5, time.Now())
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

	t.Logf("Fetched entries: %+v", entries[0])
}

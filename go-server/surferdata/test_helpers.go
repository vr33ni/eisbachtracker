package surferdata

import (
	"testing"

	"github.com/vr33ni/eisbachtracker-pwa/go-server/testutils"
)

// shared setup for all surferdata tests
func setupTestService(t *testing.T) *Service {
	testutils.LoadTestConfig(t)

	db := testutils.SetupTestDB(t)
	return NewService(db)
}

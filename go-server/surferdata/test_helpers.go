package surferdata

import (
	"testing"

	"github.com/vr33ni/eisbachtracker-pwa/go-server/conditions"
	"github.com/vr33ni/eisbachtracker-pwa/go-server/testutils"
)

func setupTestService(t *testing.T) *Service {
	testutils.LoadTestConfig(t)

	db := testutils.SetupTestDB(t)
	waterService := conditions.NewWaterService()

	return NewService(db, waterService)
}

package conditions

import (
	"sync"
	"time"
)

var (
	waterLevelCacheMu     sync.Mutex
	cachedWaterLevel      float64
	cachedWaterFlow       float64
	lastWaterLevelFetched time.Time
)

func GetCachedWaterConditions(fetchFunc func() (float64, float64, error)) (float64, float64, error) {
	waterLevelCacheMu.Lock()
	defer waterLevelCacheMu.Unlock()

	if time.Since(lastWaterLevelFetched) < time.Minute &&
		cachedWaterLevel != 0 && cachedWaterFlow != 0 {
		return cachedWaterLevel, cachedWaterFlow, nil
	}

	level, flow, err := fetchFunc()
	if err != nil {
		return 0, 0, err
	}

	cachedWaterLevel = level
	cachedWaterFlow = flow
	lastWaterLevelFetched = time.Now()
	return level, flow, nil
}

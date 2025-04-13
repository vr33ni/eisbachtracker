package conditions

import (
	"sync"
	"time"
)

var (
	cacheLock     sync.Mutex
	lastWaterTemp *float64
	lastFetched   time.Time
)

// GetCachedWaterTemperature returns cached temp if recent, else fetches
func GetCachedWaterTemperature() (float64, error) {
	cacheLock.Lock()
	defer cacheLock.Unlock()

	if lastWaterTemp != nil && time.Since(lastFetched) < 10*time.Minute {
		return *lastWaterTemp, nil
	}

	temp, err := GetLatestWaterTemperature()
	if err != nil {
		if lastWaterTemp != nil {
			// fallback to last known value
			return *lastWaterTemp, nil
		}
		return 0, err
	}

	lastWaterTemp = &temp
	lastFetched = time.Now()
	return temp, nil
}

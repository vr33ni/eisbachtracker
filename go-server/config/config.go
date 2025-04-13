package config

import (
	"os"

	"github.com/pelletier/go-toml"
)

type PredictConfig struct {
	HourFactor struct {
		Peak    float64
		Offpeak float64
	} `toml:"hour_factor"`

	WaterTempFactor struct {
		Cold   float64
		Medium float64
		Warm   float64
	} `toml:"water_temp_factor"`

	AirTempFactor struct {
		Cold   float64
		Medium float64
		Hot    float64
	} `toml:"air_temp_factor"`

	WeatherConditionFactor map[string]float64 `toml:"weather_condition_factor"`
}

var Predict PredictConfig

func LoadConfig() error {
	path := os.Getenv("PREDICT_CONFIG")
	if path == "" {
		path = "./config/predict.toml" // fallback default
	}

	config, err := toml.LoadFile(path)
	if err != nil {
		return err
	}
	return config.Unmarshal(&Predict)
}

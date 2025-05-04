package conditions

type WeatherData struct {
	Temp      float64 `json:"temp"`
	Condition int     `json:"condition"` // Use numeric WMO codes
}

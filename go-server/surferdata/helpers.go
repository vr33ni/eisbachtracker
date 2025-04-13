package surferdata

func safeFloat(f *float64) float64 {
	if f == nil {
		return 0
	}
	return *f
}

func safeString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

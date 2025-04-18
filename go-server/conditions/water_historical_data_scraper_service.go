package conditions

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// Represents a single row from the water level history table
type HistoricalWaterLevel struct {
	DateTime string
	Value    float64
}

// Scrapes historical water level values from the HND Bayern site
func ScrapeWaterLevelHistory() ([]HistoricalWaterLevel, error) {
	url := os.Getenv("HND_BAYERN_URL")

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch page: %w", err)
	}
	defer resp.Body.Close()
	// bodyBytes, _ := io.ReadAll(resp.Body)
	// fmt.Println("🧾 RAW HTML:\n", string(bodyBytes))

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	var results []HistoricalWaterLevel
	doc.Find("table.tblsort tbody tr").Each(func(i int, s *goquery.Selection) {
		cols := s.Find("td")
		if cols.Length() >= 2 {
			dateText := strings.TrimSpace(cols.Eq(0).Text())
			valueText := strings.ReplaceAll(strings.TrimSpace(cols.Eq(1).Text()), ",", ".")

			var value float64
			fmt.Sscanf(valueText, "%f", &value)

			results = append(results, HistoricalWaterLevel{
				DateTime: dateText,
				Value:    value,
			})
		}
	})

	return results, nil
}

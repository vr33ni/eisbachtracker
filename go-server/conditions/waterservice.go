package conditions

import (
	"archive/zip"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strings"
	"time"
)

type WaterService struct{}

func NewWaterService() *WaterService {
	return &WaterService{}
}

type WaterConditions struct {
	Level float64
	Flow  float64
}

// --- Interface
type WaterDataProvider interface {
	GetCurrentWeather() (*WeatherData, error)
	GetCachedWaterTemperature() (float64, error)
	GetLatestWaterTemperature() (float64, error) // <-- Add this line
	GetCurrentWaterConditions() (float64, float64, error)
}

// --- API Methods

func (ws *WaterService) GetCurrentWeather() (*WeatherData, error) {
	return GetCurrentWeather()
}

func (ws *WaterService) GetCachedWaterTemperature() (float64, error) {
	return GetCachedWaterTemperature()
}

func (ws *WaterService) GetCurrentWaterConditions() (float64, float64, error) {
	url := os.Getenv("PEGEL_ALARM_API_URL")

	resp, err := http.Get(url)
	if err != nil {
		return 0, 0, err
	}
	defer resp.Body.Close()

	var data struct {
		Payload struct {
			Stations []struct {
				Data []struct {
					Value float64 `json:"value"`
				} `json:"data"`
			} `json:"stations"`
		} `json:"payload"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return 0, 0, err
	}

	if len(data.Payload.Stations) == 0 || len(data.Payload.Stations[0].Data) < 2 {
		return 0, 0, errors.New("invalid response structure")
	}

	level := data.Payload.Stations[0].Data[0].Value
	flow := data.Payload.Stations[0].Data[1].Value

	return level, flow, nil
}

// --- Public Fetching Method ---

func (ws *WaterService) GetLatestWaterTemperature() (float64, error) {
	client, err := createHTTPClient()
	if err != nil {
		return 0, fmt.Errorf("creating HTTP client: %w", err)
	}

	token, err := requestDownloadToken(client)
	if err != nil {
		return 0, fmt.Errorf("getting token: %w", err)
	}

	downloadURL := fmt.Sprintf("https://www.gkd.bayern.de/de/downloadcenter/download?token=%s&dl=1", token)

	zipPath, err := pollAndDownloadZip(client, downloadURL)
	if err != nil {
		return 0, fmt.Errorf("downloading zip: %w", err)
	}
	defer os.Remove(zipPath)

	records, err := extractCSV(zipPath)
	if err != nil {
		return 0, fmt.Errorf("parsing CSV: %w", err)
	}

	return parseLatestWaterTemperature(records)
}

// --- Internal Helpers ---

func createHTTPClient() (*http.Client, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}
	return &http.Client{Jar: jar}, nil
}

func requestDownloadToken(client *http.Client) (string, error) {
	page := "https://www.gkd.bayern.de/de/fluesse/wassertemperatur/kelheim/muenchen-himmelreichbruecke-16515005/download"

	// Load page first (important for cookies/session)
	_, err := client.Get(page)
	if err != nil {
		return "", err
	}

	form := url.Values{
		"zr":       {"monat"},
		"beginn":   {"01.04.2025"},
		"ende":     {"05.04.2025"},
		"email":    {"test@test.de"},
		"geprueft": {"0"},
		"wertart":  {"tmw"},
		"f":        {""},
		"t":        {`{"16515005":["fluesse.wassertemperatur"]}`},
	}

	req, err := http.NewRequest("POST", "https://www.gkd.bayern.de/de/downloadcenter/enqueue_download", strings.NewReader(form.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("User-Agent", "Mozilla/5.0")
	req.Header.Set("Referer", page)
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("Origin", "https://www.gkd.bayern.de")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	bodyStr := string(body)

	tokenStart := strings.Index(bodyStr, "token=")
	if tokenStart == -1 {
		return "", fmt.Errorf("token not found")
	}

	tokenRaw := bodyStr[tokenStart+6:]
	tokenEnd := strings.IndexAny(tokenRaw, `"'><&`)
	if tokenEnd == -1 {
		tokenEnd = len(tokenRaw)
	}

	return strings.TrimSuffix(strings.TrimSpace(tokenRaw[:tokenEnd]), `\`), nil
}

func pollAndDownloadZip(client *http.Client, downloadURL string) (string, error) {
	for i := 0; i < 15; i++ {
		head, _ := client.Head(downloadURL)
		if head != nil &&
			head.StatusCode == 200 &&
			head.ContentLength > 0 &&
			strings.Contains(head.Header.Get("Content-Type"), "zip") {
			goto Ready
		}
		time.Sleep(3 * time.Second)
	}
	return "", fmt.Errorf("download not ready")

Ready:
	resp, err := client.Get(downloadURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	path := os.TempDir() + "/data.zip"

	out, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return path, err
}

func extractCSV(zipPath string) ([][]string, error) {
	zipFile, err := zip.OpenReader(zipPath)
	if err != nil {
		return nil, err
	}
	defer zipFile.Close()

	for _, f := range zipFile.File {
		if strings.HasSuffix(f.Name, ".csv") {
			rc, _ := f.Open()
			defer rc.Close()
			content, _ := io.ReadAll(rc)
			lines := strings.Split(string(content), "\n")

			for i, line := range lines {
				if strings.HasPrefix(strings.TrimSpace(line), "Datum") {
					reader := csv.NewReader(strings.NewReader(strings.Join(lines[i:], "\n")))
					reader.Comma = ';'
					reader.FieldsPerRecord = -1
					return reader.ReadAll()
				}
			}
		}
	}
	return nil, fmt.Errorf("no valid CSV found")
}

func parseLatestWaterTemperature(rows [][]string) (float64, error) {
	if len(rows) < 2 {
		return 0, fmt.Errorf("no data in CSV")
	}
	last := rows[len(rows)-1]
	tempStr := strings.ReplaceAll(last[1], ",", ".")
	var temp float64
	fmt.Sscanf(tempStr, "%f", &temp)
	return temp, nil
}

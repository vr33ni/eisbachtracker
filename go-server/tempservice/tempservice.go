package tempservice

import (
	"archive/zip"
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strings"
	"time"
)

func GetLatestTemperature() (float64, error) {
	jar, _ := cookiejar.New(nil)
	client := &http.Client{Jar: jar}

	page := "https://www.gkd.bayern.de/de/fluesse/wassertemperatur/kelheim/muenchen-himmelreichbruecke-16515005/download"

	// Step 1: Get session
	_, err := client.Get(page)
	if err != nil {
		return 0, err
	}

	// Step 2: POST form
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
		return 0, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("User-Agent", "Mozilla/5.0")
	req.Header.Set("Referer", page)
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("Origin", "https://www.gkd.bayern.de")

	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	bodyStr := string(body)

	tokenStart := strings.Index(bodyStr, "token=")
	if tokenStart == -1 {
		return 0, fmt.Errorf("token not found")
	}
	tokenRaw := bodyStr[tokenStart+6:]
	tokenEnd := strings.IndexAny(tokenRaw, `"'><&`)
	if tokenEnd == -1 {
		tokenEnd = len(tokenRaw)
	}
	token := strings.TrimSuffix(strings.TrimSpace(tokenRaw[:tokenEnd]), `\`)

	downloadUrl := fmt.Sprintf("https://www.gkd.bayern.de/de/downloadcenter/download?token=%s&dl=1", token)

	// Step 3: Poll for readiness
	isReady := false
	for i := 0; i < 15; i++ {
		head, _ := client.Head(downloadUrl)
		if head != nil && head.StatusCode == 200 && head.ContentLength > 0 && strings.Contains(head.Header.Get("Content-Type"), "zip") {
			isReady = true
			break
		}
		time.Sleep(3 * time.Second)
	}
	if !isReady {
		return 0, fmt.Errorf("download not ready")
	}

	// Step 4: Download ZIP
	zipPath := "./data.zip"
	zipRes, err := client.Get(downloadUrl)
	if err != nil {
		return 0, err
	}
	defer zipRes.Body.Close()

	if !strings.Contains(zipRes.Header.Get("Content-Type"), "zip") {
		bodyPreview, _ := io.ReadAll(zipRes.Body)
		return 0, fmt.Errorf("not a zip file. Body: %s", string(bodyPreview[:300]))
	}

	out, _ := os.Create(zipPath)
	defer out.Close()
	io.Copy(out, zipRes.Body)

	// Step 5: Unzip & find CSV
	zipFile, err := zip.OpenReader(zipPath)
	if err != nil {
		return 0, err
	}
	defer zipFile.Close()

	var csvLines []string
	for _, f := range zipFile.File {
		if strings.HasSuffix(f.Name, ".csv") {
			rc, _ := f.Open()
			defer rc.Close()
			bytes, _ := io.ReadAll(rc)
			allLines := strings.Split(string(bytes), "\n")

			for i, line := range allLines {
				if strings.HasPrefix(strings.TrimSpace(line), "Datum") {
					csvLines = allLines[i:]
					break
				}
			}
			break
		}
	}

	if len(csvLines) < 2 {
		return 0, fmt.Errorf("no temperature data found")
	}

	// Step 6: Parse CSV
	reader := csv.NewReader(strings.NewReader(strings.Join(csvLines, "\n")))
	reader.Comma = ';'
	reader.FieldsPerRecord = -1
	rows, _ := reader.ReadAll()

	last := rows[len(rows)-1]
	tempStr := strings.ReplaceAll(last[1], ",", ".")
	var temp float64
	fmt.Sscanf(tempStr, "%f", &temp)

	os.Remove(zipPath)
	return temp, nil
}

package main

import (
	"archive/zip"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strings"
	"time"
)

func main() {
	http.HandleFunc("/api/temperature", withCORS(handleTemperature))

	fmt.Println("🌍 Listening on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// CORS middleware
func withCORS(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		handler(w, r)
	}
}

func handleTemperature(w http.ResponseWriter, r *http.Request) {
	temp, err := fetchLatestTemperature()
	if err != nil {
		log.Println("❌", err)
		http.Error(w, "Internal Server Error", 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"temperature": temp,
	})
}

func fetchLatestTemperature() (float64, error) {
	jar, _ := cookiejar.New(nil)
	client := &http.Client{Jar: jar}

	page := "https://www.gkd.bayern.de/de/fluesse/wassertemperatur/kelheim/muenchen-himmelreichbruecke-16515005/download"

	// Step 1: Visit download page to get session
	_, err := client.Get(page)
	if err != nil {
		return 0, err
	}

	// Step 2: Send POST request with headers
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
	fmt.Println("📦 Response body:", bodyStr)

	tokenStart := strings.Index(bodyStr, "token=")
	if tokenStart == -1 {
		return 0, fmt.Errorf("token not found")
	}

	tokenRaw := bodyStr[tokenStart+6:]
	tokenEnd := strings.IndexAny(tokenRaw, `"'><&`)
	if tokenEnd == -1 {
		tokenEnd = len(tokenRaw)
	}
	token := tokenRaw[:tokenEnd]

	// Optional: clean leftover junk
	token = strings.TrimSpace(token)
	token = strings.TrimSuffix(token, `\`)

	downloadUrl := fmt.Sprintf("https://www.gkd.bayern.de/de/downloadcenter/download?token=%s&dl=1", token)
	fmt.Println("⬇️ Download URL:", downloadUrl)

	// Step 3: Wait for file to be ready
	isReady := false
	for i := 0; i < 15; i++ {
		head, _ := client.Head(downloadUrl)
		length := head.ContentLength
		status := head.StatusCode
		contentType := head.Header.Get("Content-Type")

		log.Printf("🔁 Poll %d - Status: %d, Length: %d, Type: %s", i+1, status, length, contentType)

		if status == 200 && strings.Contains(contentType, "zip") && length > 0 {
			isReady = true
			break
		}

		time.Sleep(3 * time.Second)
	}

	if !isReady {
		return 0, fmt.Errorf("download not ready")
	}

	// Step 4: Download zip
	zipPath := "./data.zip"
	zipRes, err := client.Get(downloadUrl)
	if err != nil {
		return 0, err
	}
	defer zipRes.Body.Close()

	contentType := zipRes.Header.Get("Content-Type")
	if !strings.Contains(contentType, "zip") {
		bodyPreview, _ := io.ReadAll(zipRes.Body)
		return 0, fmt.Errorf("not a zip file: %s\nBody: %s", contentType, string(bodyPreview[:300]))
	}

	out, _ := os.Create(zipPath)
	defer out.Close()
	io.Copy(out, zipRes.Body)

	// Step 5: Unzip and parse CSV
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

	reader := csv.NewReader(strings.NewReader(strings.Join(csvLines, "\n")))
	reader.Comma = ';'
	reader.FieldsPerRecord = -1
	rows, _ := reader.ReadAll()

	last := rows[len(rows)-1]
	tempStr := strings.ReplaceAll(last[1], ",", ".")
	var temp float64
	fmt.Sscanf(tempStr, "%f", &temp)

	// Clean up
	os.Remove(zipPath)

	return temp, nil
}

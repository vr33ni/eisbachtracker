package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/vr33ni/eisbachtracker-pwa/go-server/db"
	"github.com/vr33ni/eisbachtracker-pwa/go-server/surferdata"
	"github.com/vr33ni/eisbachtracker-pwa/go-server/tempservice"
)

var surferService *surferdata.Service

func main() {
	if os.Getenv("FLY_APP_NAME") == "" {
		err := godotenv.Load()
		if err != nil {
			log.Println("⚠️ Could not load .env file")
		} else {
			log.Println("✅ Loaded local .env")
		}
	}

	if err := db.Init(); err != nil {
		log.Fatal(err)
	}
	defer db.Conn.Close()

	fmt.Println("🌍 DATABASE_URL:", os.Getenv("DATABASE_URL"))

	surferService = surferdata.NewService(db.Conn)

	http.HandleFunc("/api/temperature", withCORS(handleTemperature))
	http.HandleFunc("/api/surfers", withCORS(handleSurferEntries))
	http.HandleFunc("/api/surfers/predict", withCORS(handlePrediction(surferService)))

	if os.Getenv("ENV") != "production" {
		runMigrations()
	}
	fmt.Println("🚀 Listening on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func runMigrations() {
	cmd := exec.Command("flyway", "migrate")
	if err := cmd.Run(); err != nil {
		log.Fatalf("Failed to execute Flyway migrations: %v", err)
	}
	log.Println("Database migrations applied successfully.")
}

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
	temp, err := tempservice.GetLatestTemperature()
	if err != nil {
		log.Println("❌", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"temperature": temp,
	})
}

func handlePrediction(service *surferdata.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		hourStr := r.URL.Query().Get("hour")
		tempStr := r.URL.Query().Get("temperature")

		hour, err := strconv.Atoi(hourStr)
		if err != nil {
			http.Error(w, "Invalid or missing hour", http.StatusBadRequest)
			return
		}

		var tempPtr *float64
		if tempStr != "" {
			t, err := strconv.ParseFloat(tempStr, 64)
			if err == nil {
				tempPtr = &t
			}
		}

		prediction, err := service.PredictSurferCount(hour, tempPtr)
		if err != nil {
			http.Error(w, "Could not compute prediction: "+err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(map[string]interface{}{
			"hour":        hour,
			"temperature": tempPtr,
			"prediction":  prediction,
		})
	}
}

// 🏄‍♂️ Handle GET + POST /api/surfers
func handleSurferEntries(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		entries, err := surferService.GetAllEntries()
		if err != nil {
			http.Error(w, "Failed to fetch entries", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(entries)

	case http.MethodPost:
		var input struct {
			Count int       `json:"count"`
			Time  time.Time `json:"timestamp"` // optional
		}
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}
		if input.Count <= 0 {
			http.Error(w, "Surfer count must be positive", http.StatusBadRequest)
			return
		}

		if err := surferService.AddEntry(input.Count, input.Time); err != nil {
			http.Error(w, "Failed to save entry", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"message": "Entry saved"})

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

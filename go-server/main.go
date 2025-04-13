package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"

	"github.com/joho/godotenv"
	"github.com/vr33ni/eisbachtracker-pwa/go-server/config"
	"github.com/vr33ni/eisbachtracker-pwa/go-server/db"
	"github.com/vr33ni/eisbachtracker-pwa/go-server/routes"
	"github.com/vr33ni/eisbachtracker-pwa/go-server/surferdata"
)

func main() {
	if err := config.LoadConfig(); err != nil {
		log.Fatal("Failed to load config:", err)
	}

	if os.Getenv("ENV") != "production" {
		if err := godotenv.Load(); err == nil {
			log.Println("✅ Loaded local .env")
		}
	}

	if err := db.Init(); err != nil {
		log.Fatal(err)
	}
	defer db.Conn.Close()

	fmt.Println("🌍 DATABASE_URL:", os.Getenv("DATABASE_URL"))

	surferService := surferdata.NewService(db.Conn)

	// Register Routes
	routes.RegisterRoutes(surferService)

	// Run Migrations locally only
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

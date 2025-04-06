package main

import (
	"fmt"
	"log"

	"github.com/vr33ni/eisbachtracker-pwa/go-server/tempservice"
)

func main() {
	temp, err := tempservice.GetLatestTemperature()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("🌡️ Temperature: %.2f°C\n", temp)
}

# 🌊 EisbachTracker PWA

## Backend (Go)

The Go server downloads a CSV file from gkd.bayern.de, unzips it, and extracts the latest water temperature.

```cmd
cd ../go-server
go run main.go
```

This serves the /api/temperature endpoint on port 8080.

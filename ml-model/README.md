# 🌊 EisbachTracker PWA

## Prediction model

The prediction model uses a machine learning approach to estimate the number of surfers based on environmental conditions and time of day. A linear regression model is trained on dummy data that simulates realistic tendencies, such as more surfers on sunny days, during warmer temperatures, and at peak surfing hours (e.g., early morning or lunchtime). The model is integrated into a Flask API, allowing real-time predictions based on user-provided inputs like weather, water temperature, and time. This setup serves as a foundation for future enhancements with real-world data.

### Features Used:
- Hour of the day: Captures time-based surfing patterns.
- Water temperature: Warmer water tends to attract more surfers.
- Air temperature: Warmer air increases surfer activity.
- Water level: Determines surfability of the wave.
- Weather conditions: Includes sunny, cloudy, rainy, snowy, and stormy conditions.

### API Endpoint:

POST /predict: Accepts JSON input with the above features and returns the predicted surfer count.

### Deployment:

The Flask API is hosted on Render and communicates with the Go backend for seamless integration.

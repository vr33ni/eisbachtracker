# 🌊 EisbachTracker PWA

## Prediction model

The prediction model uses a machine learning approach to estimate the number of surfers based on environmental conditions and time of day. A linear regression model is trained on dummy data that simulates realistic tendencies, such as more surfers on sunny days, during warmer temperatures, and at peak surfing hours (e.g., early morning or lunchtime). The model is integrated into a Flask API, allowing real-time predictions based on user-provided inputs like weather, water temperature, and time. This setup serves as a foundation for future enhancements with real-world data.

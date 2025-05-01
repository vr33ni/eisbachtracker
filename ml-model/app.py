from flask import Flask, request, jsonify
import joblib
import numpy as np
import pandas as pd  # Import pandas

app = Flask(__name__)
model = joblib.load("surfer_prediction_model.pkl")

@app.route("/predict", methods=["POST"])
def predict():
    data = request.json
    # Create a DataFrame with the same column names as the training data
    feature_dict = {
        "hour": data["hour"],
        "water_temp": data["water_temp"],
        "air_temp": data["air_temp"],
        "water_level": data["water_level"],
        "weather_condition_cloudy": data.get("weather_condition_cloudy", 0),
        "weather_condition_rainy": data.get("weather_condition_rainy", 0),
        "weather_condition_snow": data.get("weather_condition_snow", 0),
        "weather_condition_stormy": data.get("weather_condition_stormy", 0),
    }
    features = pd.DataFrame([feature_dict])  # Create a DataFrame with one row
    prediction = model.predict(features)
    return jsonify({"surfer_count": int(prediction[0])})

if __name__ == "__main__":
    app.run(debug=True)
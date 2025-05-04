from flask import Flask, request, jsonify
import joblib
import numpy as np
import pandas as pd 
import os 

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
        "weather_condition": data.get("weather_condition"),
    }
    features = pd.DataFrame([feature_dict])  # Create a DataFrame with one row

    # Initialize default values for surfer_count and explanation
    surfer_count = 0
    explanation = {
        "hour": 0.0,
        "water_temp": 0.0,
        "air_temp": 0.0,
        "water_level": 0.0,
        "weather_condition": 0.0,
    }

    # Enforce the rule: water_level < 130 means no surfers
    if feature_dict["water_level"] >= 130:
        prediction = model.predict(features)
        surfer_count = max(0, int(prediction[0]))  # Ensure the count is at least 0

        # Calculate feature contributions
        coefficients = model.coef_  # Get model coefficients
        contributions = {
            feature: coef * features.iloc[0][feature]
            for feature, coef in zip(features.columns, coefficients)
        }
        explanation = contributions  # Feature contributions to the prediction

    return jsonify({"surfer_count": surfer_count, "explanation": explanation})

if __name__ == "__main__":
    app.run(host="0.0.0.0", port=int(os.environ.get("PORT", 5001)))
from flask import Flask, request, jsonify
import joblib
import numpy as np
import pandas as pd  # Import pandas
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

    # Enforce the rule: water_level < 130 means no surfers
    if feature_dict["water_level"] < 130:
        surfer_count = 0
    else:
        prediction = model.predict(features)
        surfer_count = max(0, int(prediction[0]))  # Ensure the count is at least 0

    return jsonify({"surfer_count": surfer_count})

if __name__ == "__main__":
    app.run(host="0.0.0.0", port=int(os.environ.get("PORT", 5001)))
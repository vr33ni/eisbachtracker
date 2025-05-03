import pandas as pd
from sklearn.model_selection import train_test_split
from sklearn.linear_model import LinearRegression
from sklearn.metrics import mean_squared_error
import joblib

# Load data
df = pd.read_csv("combined_feature_and_target_data.csv")

# Ensure all one-hot encoded columns are present
expected_columns = [
    "hour", "water_temp", "air_temp", "water_level",
    "weather_condition"
]
for col in expected_columns:
    if col not in df.columns:
        df[col] = 0  # Add missing columns with default value 0

# Features (X) and target (y)
X = df[expected_columns]  # Use only the expected columns
y = df["surfer_count"]  # Target column

# Split data into training and testing sets
X_train, X_test, y_train, y_test = train_test_split(X, y, test_size=0.2, random_state=42)

# Train a linear regression model
model = LinearRegression()
model.fit(X_train, y_train)

# Evaluate the model
y_pred = model.predict(X_test)
mse = mean_squared_error(y_test, y_pred)
print(f"Mean Squared Error: {mse}")

# Save the model
joblib.dump(model, "surfer_prediction_model.pkl")
print("Model saved as surfer_prediction_model.pkl")
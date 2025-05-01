import pandas as pd
import numpy as np

# Generate dummy data
np.random.seed(42)
data = {
    "hour": np.random.randint(0, 24, 1000),  # Hour of the day
    "water_temp": np.random.uniform(2, 20, 1000),  # Water temperature
    "air_temp": np.random.uniform(-10, 35, 1000),  # Air temperature
    "weather_condition": np.random.choice(["sunny", "cloudy", "rainy", "snow", "stormy"], 1000),  # Weather
    "water_level": np.random.uniform(130, 155, 1000),  # Water level
    #"water_flow": np .random.uniform(18, 28, 1000),  # Water flow
}

# Convert to DataFrame
df = pd.DataFrame(data)

# Add surfer_count with tendencies
def generate_surfer_count(row):
    base_count = np.random.randint(0, 10)  # Base random count
    if row["water_level"] < 130:
        return 0  #  Not surfable
    if row["water_level"] > 145:
        base_count += 5  # More surfers with higher water levels
    if row["water_level"] < 140:
        base_count -= 8  # Less surfers with lower water levels
    if row["weather_condition"] == "sunny":
        base_count += 5  # More surfers on sunny days
    if row["weather_condition"] == "snow":
        base_count -= 3  # Less surfers on winter days
    if row["air_temp"] > 20:
        base_count += 5  # More surfers in warmer air temperatures
    if row["air_temp"] < 0:
        base_count -= 10  # Less surfers in colder air temperatures
    if row["water_temp"] > 15:
        base_count += 5  # More surfers in warmer water temperatures
    if row["water_temp"] < 5:
        base_count -= 8  # Less surfers in lower water temperatures
    # if row["water_flow"] > 25:
    #     base_count += 5  # More surfers with higher water flow; water flow kinda equivalent with higher water level
    if 11 <= row["hour"] <= 14 or 6 <= row["hour"] <= 7 or 17 <= row["hour"] <= 19:  # Lunch/afternoon or early morning
        base_count += 5
    return max(0, min(base_count, 30))  # Ensure count is between 0 and 30
df["surfer_count"] = df.apply(generate_surfer_count, axis=1)

# One-hot encode weather_condition
df = pd.get_dummies(df, columns=["weather_condition"], drop_first=True)

# Save to CSV
df.to_csv("dummy_surfer_data.csv", index=False)
print("Dummy data saved to dummy_surfer_data.csv")
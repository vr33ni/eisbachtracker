import os
import pandas as pd
from bs4 import BeautifulSoup
import requests
import numpy as np
from openmeteo_requests import Client
import requests_cache
from retry_requests import retry

# Function to process temperature data
def process_temperature_data(folder_path):
    all_data = []

    # Iterate through all files in the folder
    for filename in os.listdir(folder_path):
        if filename.startswith("temp_") and filename.endswith(".csv"):
            file_path = os.path.join(folder_path, filename)

            # Read the file, skipping the first 9 rows
            temp_data = pd.read_csv(file_path, sep=";", skiprows=9, header=0)

            # Rename columns to standard names
            temp_data = temp_data.rename(columns={"Datum": "date", "Mittelwert": "water_temp"})

            # Convert date and temperature columns
            temp_data["date"] = pd.to_datetime(temp_data["date"], format="%Y-%m-%d")
            temp_data["water_temp"] = temp_data["water_temp"].str.replace(",", ".").astype(float)

            # Append to the list of all data
            all_data.append(temp_data[["date", "water_temp"]])

    # Combine all data into a single DataFrame
    combined_data = pd.concat(all_data, ignore_index=True)

    # Sort by date to ensure chronological order
    combined_data = combined_data.sort_values(by="date").reset_index(drop=True)

    return combined_data

# Function to scrape water level data
def scrape_historical_water_level(url):
    # Fetch the HTML page
    response = requests.get(url)
    if response.status_code != 200:
        raise Exception(f"Failed to fetch page: {response.status_code}")
    
    # Parse the HTML content
    soup = BeautifulSoup(response.content, "html.parser")
    table = soup.find("table", class_="tblsort")
    if not table:
        raise Exception("Failed to find the water level table in the HTML")

    # Extract rows from the table
    rows = table.find("tbody").find_all("tr")
    data = []
    for row in rows:
        cols = row.find_all("td")
        if len(cols) >= 2:
            date_text = cols[0].get_text(strip=True)
            value_text = cols[1].get_text(strip=True).replace(",", ".")
            try:
                value = float(value_text)
                data.append({"date": date_text, "water_level": value})
            except ValueError:
                continue  # Skip rows with invalid data

    # Convert to DataFrame
    df = pd.DataFrame(data)
    df["date"] = pd.to_datetime(df["date"], format="%d.%m.%Y %H:%M", errors="coerce")  # Convert to datetime
    return df

# Function to expand temperature data to include all hours of the day
def expand_temperature_data(temperature_data):
    # Repeat each daily temperature for all 24 hours
    expanded_data = temperature_data.loc[temperature_data.index.repeat(24)].reset_index(drop=True)
    expanded_data["hour"] = list(range(24)) * len(temperature_data)
    return expanded_data



# Function to fetch weather data
def fetch_weather_data():
    """Fetch weather data from Open-Meteo API."""
    cache_session = requests_cache.CachedSession('.cache', expire_after=-1)
    retry_session = retry(cache_session, retries=5, backoff_factor=0.2)
    openmeteo = Client(session=retry_session)

    # Define API parameters
    url = "https://archive-api.open-meteo.com/v1/archive"
    params = {
        "latitude": 48.137154,
        "longitude": 11.576124,
        "start_date": "2024-01-01",  # Adjust as needed
        "end_date": "2024-12-31",    # Adjust as needed
        "hourly": ["temperature_2m", "weather_code"]
    }

    # Fetch weather data
    responses = openmeteo.weather_api(url, params=params)
    response = responses[0]  # Assuming single location

    # Process hourly data
    hourly = response.Hourly()
    weather_data = pd.DataFrame({
        "date": pd.date_range(
            start=pd.to_datetime(hourly.Time(), unit="s", utc=True),
            end=pd.to_datetime(hourly.TimeEnd(), unit="s", utc=True),
            freq=pd.Timedelta(seconds=hourly.Interval()),
            inclusive="left"
        ),
        "hour": pd.to_datetime(hourly.Time(), unit="s", utc=True).hour,
        "air_temp": hourly.Variables(0).ValuesAsNumpy(),
        "weather_code": hourly.Variables(1).ValuesAsNumpy()
    })

    # Ensure date column is in local time (optional)
    weather_data["date"] = weather_data["date"].dt.tz_convert(None)
    return weather_data



# Combine temperature, water level, and weather data
def combine_temperature_and_water_level_with_weather(temp_folder, water_level_url):
    try:
        print("Processing temperature data...")
        temperature_data = process_temperature_data(temp_folder)
        print(f"Temperature data shape: {temperature_data.shape}")
        if temperature_data.empty:
            print("Temperature data is empty! Check your temp_folder path or CSV files.")
            return

        print("Scraping water level data...")
        water_level_data = scrape_historical_water_level(water_level_url)
        print(f"Water level data shape: {water_level_data.shape}")
        if water_level_data.empty:
            print("Water level data is empty! Check your water_level_url.")
            return

        print("Extracting hour from water level data...")
        water_level_data["hour"] = water_level_data["date"].dt.hour

        print("Ensuring date columns are consistent...")
        water_level_data["date"] = pd.to_datetime(water_level_data["date"].dt.date)  # Keep only the date part
        temperature_data["date"] = pd.to_datetime(temperature_data["date"])

        print("Expanding temperature data to include all hours...")
        expanded_temperature_data = expand_temperature_data(temperature_data)
        print(f"Expanded temperature data shape: {expanded_temperature_data.shape}")

        print("Merging temperature and water level data...")
        combined_data = pd.merge(
            expanded_temperature_data,
            water_level_data,
            on=["date", "hour"],
            how="inner"
        )
        print(f"Combined data shape: {combined_data.shape}")

        if combined_data.empty:
            print("Combined data is empty! Check your merge logic.")
            return

        print("Fetching weather data...")
        weather_data = fetch_weather_data()
        print(f"Weather data shape: {weather_data.shape}")

        print("Merging weather data with combined data...")
        combined_data = pd.merge(
            combined_data,
            weather_data,
            on=["date", "hour"],
            how="left"  # Use "left" to keep all rows from combined_data
        )
        print(f"Final combined data shape: {combined_data.shape}")

        print("Replacing dummy columns with actual weather data...")
        import pdb; pdb.set_trace()

        combined_data["air_temp"] = weather_data["air_temp"]
        combined_data["weather_condition"] = weather_data["weather_code"]
        combined_data = combined_data.drop(columns=["weather_code"])        
        print("Saving combined data to CSV...")
        combined_data.to_csv("combined_feature_data.csv", index=False)
        print("Combined feature data saved to combined_feature_data.csv")

        return combined_data

    except Exception as e:
        print(f"An error occurred: {e}")

# Example usage
if __name__ == "__main__":
    folder_path = "temperature-data"  # Replace with the path to your folder containing the CSV files
    water_level_url = "https://www.hnd.bayern.de/pegel/isar/muenchen-himmelreichbruecke-16515005/tabelle?methode=wasserstand&days=365"  # Replace with the correct URL

    try:
        combined_data = combine_temperature_and_water_level_with_weather(folder_path, water_level_url)
        if combined_data is not None:
            print("Script executed successfully!")
    except Exception as e:
        print(f"An error occurred during execution: {e}")
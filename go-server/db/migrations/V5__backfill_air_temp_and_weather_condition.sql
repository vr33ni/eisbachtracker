UPDATE surfer_entries
SET air_temperature = 18,          -- could vary this later
    weather_condition = -1         -- safe default
WHERE air_temperature IS NULL;

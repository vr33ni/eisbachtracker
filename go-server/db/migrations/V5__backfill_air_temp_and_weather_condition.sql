UPDATE surfer_entries
SET air_temperature = 18,               -- could vary this later
    weather_condition = 'Clear'         -- safe default
WHERE air_temperature IS NULL;

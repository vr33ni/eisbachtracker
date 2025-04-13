-- Set existing NULL values to 0
UPDATE surfer_entries SET water_temperature = 0 WHERE water_temperature IS NULL;
UPDATE surfer_entries SET air_temperature = 0 WHERE air_temperature IS NULL;
UPDATE surfer_entries SET weather_condition = 'Unknown' WHERE weather_condition IS NULL;

-- Set DEFAULT 0 and NOT NULL constraint
ALTER TABLE surfer_entries
ALTER COLUMN water_temperature SET DEFAULT 0,
ALTER COLUMN water_temperature SET NOT NULL,
ALTER COLUMN air_temperature SET DEFAULT 0,
ALTER COLUMN air_temperature SET NOT NULL,
ALTER COLUMN weather_condition SET DEFAULT 'Unknown',
ALTER COLUMN weather_condition SET NOT NULL;

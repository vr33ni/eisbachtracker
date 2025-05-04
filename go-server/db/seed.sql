INSERT INTO surfer_entries (timestamp, count, water_temperature, air_temperature, weather_condition) VALUES
(NOW() - interval '10 minutes', 6, 17, 23, 3),
(NOW() - interval '30 minutes', 2, 15, 20, 2),
(NOW() - interval '1 hour', 1, 14, 18, 61); 

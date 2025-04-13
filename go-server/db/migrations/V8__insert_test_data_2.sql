INSERT INTO surfer_entries (timestamp, count, water_temperature, air_temperature, weather_condition) VALUES
    (NOW() - INTERVAL '1 hour', 3, 10.5, 15.0, 'Cloudy'),
    (NOW() - INTERVAL '2 hours', 5, 10.7, 16.0, 'Clear'),
    (NOW() - INTERVAL '3 hours', 2, 10.2, 14.5, 'Rain'),
    (NOW() - INTERVAL '4 hours', 7, 10.9, 17.0, 'Clear'),
    (NOW() - INTERVAL '5 hours', 4, 9.8, 13.0, 'Cloudy'),
    (NOW() - INTERVAL '6 hours', 1, 9.5, 12.0, 'Rain'),
    (NOW() - INTERVAL '7 hours', 6, 11.0, 18.0, 'Clear');

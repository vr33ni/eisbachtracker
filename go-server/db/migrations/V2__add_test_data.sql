INSERT INTO surfer_entries (timestamp, count, temperature) VALUES
  (NOW() - interval '1 day', 3, 15.0),
  (NOW() - interval '2 hours', 5, 18.5),
  (NOW() - interval '1 hour', 6, 20.0),
  (NOW(), 4, 19.0),
  (date_trunc('day', NOW() - interval '5 days') + interval '18 hours', 8, 15.0),
  (date_trunc('day', NOW() - interval '4 days') + interval '18 hours', 10, 16.5),
  (date_trunc('day', NOW() - interval '3 days') + interval '18 hours', 7, 17.0),
  (date_trunc('day', NOW() - interval '2 days') + interval '18 hours', 12, 18.0),
  (date_trunc('day', NOW() - interval '1 days') + interval '18 hours', 9, 19.5),
  (date_trunc('day', NOW()) + interval '18 hours', 11, 20.0);


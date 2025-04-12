CREATE TABLE IF NOT EXISTS surfer_entries (
  id SERIAL PRIMARY KEY,
  timestamp TIMESTAMP NOT NULL,
  count INTEGER NOT NULL,
  temperature REAL
);

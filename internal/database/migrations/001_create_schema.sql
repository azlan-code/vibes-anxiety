-- enable timescaledb (run once per cluster)
CREATE EXTENSION IF NOT EXISTS timescaledb CASCADE;

-- create the base table if it doesn't exist (use timestamptz for time)
CREATE TABLE IF NOT EXISTS location_vibes (
  iso3 TEXT NOT NULL,
  coordinates POINT NOT NULL,
  city TEXT NOT NULL,
  day timestamptz NOT NULL,
  score double precision NOT NULL,
  PRIMARY KEY (coordinates, day)
);

-- convert to hypertable (if_not_exists avoids error if already hypertable)
SELECT create_hypertable('location_vibes', 'day', if_not_exists => TRUE, chunk_time_interval => INTERVAL '7 days');
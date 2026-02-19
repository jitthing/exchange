CREATE TABLE IF NOT EXISTS academic_events (
    id         TEXT PRIMARY KEY,
    type       TEXT NOT NULL,
    title      TEXT NOT NULL,
    start_date TEXT NOT NULL,
    end_date   TEXT NOT NULL,
    priority   INT  NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS travel_windows (
    id         TEXT PRIMARY KEY,
    start_date TEXT NOT NULL,
    end_date   TEXT NOT NULL,
    score      INT  NOT NULL DEFAULT 0,
    conflicts  JSONB NOT NULL DEFAULT '[]'
);

CREATE TABLE IF NOT EXISTS trips (
    id             TEXT PRIMARY KEY,
    owner_id       TEXT NOT NULL,
    destination    TEXT NOT NULL,
    window_id      TEXT NOT NULL,
    members        JSONB NOT NULL DEFAULT '[]',
    itinerary      JSONB NOT NULL DEFAULT '[]',
    estimated_cost DOUBLE PRECISION NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS budget_entries (
    id       TEXT PRIMARY KEY,
    user_id  TEXT NOT NULL,
    category TEXT NOT NULL,
    amount   DOUBLE PRECISION NOT NULL,
    currency TEXT NOT NULL DEFAULT 'EUR',
    date     TEXT NOT NULL,
    trip_id  TEXT NOT NULL DEFAULT '',
    note     TEXT NOT NULL DEFAULT ''
);

CREATE TABLE IF NOT EXISTS monthly_budgets (
    user_id TEXT PRIMARY KEY,
    budget  DOUBLE PRECISION NOT NULL DEFAULT 900
);

CREATE TABLE IF NOT EXISTS destinations (
    city             TEXT PRIMARY KEY,
    base_travel_hrs  DOUBLE PRECISION NOT NULL,
    transport_base   DOUBLE PRECISION NOT NULL,
    hostel_night_eur DOUBLE PRECISION NOT NULL,
    tags             JSONB NOT NULL DEFAULT '[]'
);

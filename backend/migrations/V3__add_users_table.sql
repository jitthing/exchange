-- Users table (maps to Supabase auth.users)
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY,
    email TEXT NOT NULL UNIQUE,
    display_name TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Add user_id to trips
ALTER TABLE trips ADD COLUMN IF NOT EXISTS user_id UUID REFERENCES users(id);

-- Add index on budget_entries(user_id) for fast lookups
CREATE INDEX IF NOT EXISTS idx_budget_entries_user_id ON budget_entries(user_id);

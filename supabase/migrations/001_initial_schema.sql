-- Territory Run Phase 1 schema
-- Run in Supabase Dashboard → SQL Editor → New query → Run

-- Extensions (optional for future heatmaps)
-- CREATE EXTENSION IF NOT EXISTS postgis;

-- Users (synced from Supabase Auth on first API call)
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY,
    username TEXT UNIQUE NOT NULL,
    display_name TEXT NOT NULL,
    email TEXT,
    avatar_url TEXT,
    trust_score SMALLINT DEFAULT 100,
    created_at TIMESTAMPTZ DEFAULT now()
);

-- Activities
CREATE TABLE IF NOT EXISTS activities (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    started_at TIMESTAMPTZ,
    ended_at TIMESTAMPTZ,
    distance_m REAL DEFAULT 0,
    duration_s INT DEFAULT 0,
    status TEXT NOT NULL DEFAULT 'recording',
    trust_score SMALLINT DEFAULT 100,
    capture_score INT DEFAULT 0,
    idempotency_key TEXT UNIQUE,
    device_info JSONB DEFAULT '{}'::jsonb,
    points_json JSONB DEFAULT '[]'::jsonb,
    h3_cells JSONB DEFAULT '[]'::jsonb,
    created_at TIMESTAMPTZ DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_activities_user ON activities(user_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_activities_status ON activities(status);

-- Territory cells (H3 resolution 10)
CREATE TABLE IF NOT EXISTS territory_cells (
    h3_index BIGINT NOT NULL,
    h3_index_hex TEXT PRIMARY KEY,
    parent_h3_hex TEXT NOT NULL,
    owner_id UUID REFERENCES users(id) ON DELETE SET NULL,
    status TEXT NOT NULL DEFAULT 'OWNED',
    captured_at TIMESTAMPTZ DEFAULT now(),
    last_defended_at TIMESTAMPTZ,
    capture_count INT DEFAULT 1
);

CREATE INDEX IF NOT EXISTS idx_territory_parent ON territory_cells(parent_h3_hex);
CREATE INDEX IF NOT EXISTS idx_territory_owner ON territory_cells(owner_id);

-- Territory events log
CREATE TABLE IF NOT EXISTS territory_events (
    id UUID PRIMARY KEY,
    event_type TEXT NOT NULL,
    h3_parent_hex TEXT,
    actor_id UUID REFERENCES users(id),
    victim_id UUID REFERENCES users(id),
    activity_id UUID REFERENCES activities(id),
    metadata JSONB DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_territory_events_actor ON territory_events(actor_id, created_at DESC);

-- User stats (materialized counters)
CREATE TABLE IF NOT EXISTS user_territory_stats (
    user_id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    cells_owned INT DEFAULT 0,
    territories_captured INT DEFAULT 0,
    territories_stolen INT DEFAULT 0,
    territories_lost INT DEFAULT 0,
    total_capture_score BIGINT DEFAULT 0,
    activity_count INT DEFAULT 0,
    updated_at TIMESTAMPTZ DEFAULT now()
);

-- Row Level Security (optional — API uses service role via backend)
ALTER TABLE users ENABLE ROW LEVEL SECURITY;
ALTER TABLE activities ENABLE ROW LEVEL SECURITY;
ALTER TABLE territory_cells ENABLE ROW LEVEL SECURITY;

-- Allow authenticated users to read their own data via Supabase client (optional)
CREATE POLICY "users_read_own" ON users FOR SELECT USING (auth.uid() = id);
CREATE POLICY "activities_read_own" ON activities FOR SELECT USING (auth.uid() = user_id);

-- Realtime (optional Phase 2) — enable in Dashboard → Database → Replication
-- ALTER PUBLICATION supabase_realtime ADD TABLE territory_events;

-- Territory Run Phase 2: steal, decay, leaderboards, battles
-- Run in Supabase Dashboard → SQL Editor after 001_initial_schema.sql

-- Steal / decay counters already exist on user_territory_stats:
--   territories_stolen, territories_lost
-- territory_cells already has last_defended_at for decay.

ALTER TABLE territory_cells
    ALTER COLUMN last_defended_at SET DEFAULT now();

UPDATE territory_cells
SET last_defended_at = COALESCE(last_defended_at, captured_at, now())
WHERE last_defended_at IS NULL;

-- Fast reads for global / local leaderboards
CREATE TABLE IF NOT EXISTS leaderboard_snapshots (
    scope TEXT NOT NULL,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    rank INT NOT NULL,
    score BIGINT NOT NULL DEFAULT 0,
    cells_owned INT NOT NULL DEFAULT 0,
    territories_stolen INT NOT NULL DEFAULT 0,
    refreshed_at TIMESTAMPTZ DEFAULT now(),
    PRIMARY KEY (scope, user_id)
);

CREATE INDEX IF NOT EXISTS idx_leaderboard_scope_rank
    ON leaderboard_snapshots (scope, rank);

-- Head-to-head battles (created when the same pair steals 3+ times in 7 days)
CREATE TABLE IF NOT EXISTS battles (
    id UUID PRIMARY KEY,
    challenger_id UUID NOT NULL REFERENCES users(id),
    defender_id UUID NOT NULL REFERENCES users(id),
    status TEXT NOT NULL DEFAULT 'active',
    challenger_score INT DEFAULT 0,
    defender_score INT DEFAULT 0,
    started_at TIMESTAMPTZ DEFAULT now(),
    ends_at TIMESTAMPTZ,
    CONSTRAINT battles_distinct_users CHECK (challenger_id <> defender_id)
);

CREATE INDEX IF NOT EXISTS idx_battles_challenger ON battles (challenger_id, started_at DESC);
CREATE INDEX IF NOT EXISTS idx_battles_defender ON battles (defender_id, started_at DESC);

CREATE INDEX IF NOT EXISTS idx_territory_events_victim
    ON territory_events (victim_id, created_at DESC);

CREATE INDEX IF NOT EXISTS idx_territory_events_type
    ON territory_events (event_type, created_at DESC);

CREATE INDEX IF NOT EXISTS idx_territory_cells_defended
    ON territory_cells (last_defended_at);

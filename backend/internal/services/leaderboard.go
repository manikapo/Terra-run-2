package services

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/territory-run/api/internal/models"
)

func (s *TerritoryService) Decay(ctx context.Context) (*models.DecayResult, error) {
	days := s.cfg.TerritoryDecayDays
	if days <= 0 {
		days = 14
	}
	tag, err := s.db.Exec(ctx, `
		UPDATE territory_cells
		SET owner_id = NULL, status = 'NEUTRAL'
		WHERE owner_id IS NOT NULL
		  AND COALESCE(last_defended_at, captured_at, now()) < now() - ($1::int || ' days')::interval
	`, days)
	if err != nil {
		return nil, err
	}
	return &models.DecayResult{Neutralized: int(tag.RowsAffected())}, nil
}

func (s *TerritoryService) GlobalLeaderboard(ctx context.Context, limit int) (*models.LeaderboardResponse, error) {
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	rows, err := s.db.Query(ctx, `
		SELECT u.id, u.username, u.display_name,
			COALESCE(s.cells_owned, 0),
			COALESCE(s.territories_captured, 0),
			COALESCE(s.territories_stolen, 0),
			COALESCE(s.total_capture_score, 0)
		FROM user_territory_stats s
		JOIN users u ON u.id = s.user_id
		ORDER BY s.total_capture_score DESC, s.cells_owned DESC
		LIMIT $1
	`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanLeaderboard(rows, "global")
}

func (s *TerritoryService) LocalLeaderboard(ctx context.Context, parentHex string, limit int) (*models.LeaderboardResponse, error) {
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	rows, err := s.db.Query(ctx, `
		SELECT u.id, u.username, u.display_name,
			COUNT(*)::int AS cells_owned,
			COUNT(DISTINCT c.parent_h3_hex)::int AS territories_captured,
			COALESCE(s.territories_stolen, 0),
			COALESCE(s.total_capture_score, 0)
		FROM territory_cells c
		JOIN users u ON u.id = c.owner_id
		LEFT JOIN user_territory_stats s ON s.user_id = u.id
		WHERE c.parent_h3_hex = $1 AND c.owner_id IS NOT NULL
		GROUP BY u.id, u.username, u.display_name, s.territories_stolen, s.total_capture_score
		ORDER BY cells_owned DESC, s.total_capture_score DESC
		LIMIT $2
	`, parentHex, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanLeaderboard(rows, "local:"+parentHex)
}

func scanLeaderboard(rows interface {
	Next() bool
	Scan(dest ...any) error
	Err() error
}, scope string) (*models.LeaderboardResponse, error) {
	out := &models.LeaderboardResponse{Scope: scope, Entries: []models.LeaderboardEntry{}}
	rank := 0
	for rows.Next() {
		rank++
		var e models.LeaderboardEntry
		if err := rows.Scan(&e.UserID, &e.Username, &e.DisplayName,
			&e.CellsOwned, &e.TerritoriesCaptured, &e.TerritoriesStolen, &e.TotalCaptureScore); err != nil {
			return nil, err
		}
		e.Rank = rank
		out.Entries = append(out.Entries, e)
	}
	return out, rows.Err()
}

func (s *TerritoryService) RecentEvents(ctx context.Context, userID uuid.UUID, limit int) ([]models.TerritoryEvent, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	rows, err := s.db.Query(ctx, `
		SELECT id, event_type, COALESCE(h3_parent_hex, ''), actor_id, victim_id, activity_id,
			COALESCE(metadata, '{}'::jsonb), created_at
		FROM territory_events
		WHERE actor_id = $1 OR victim_id = $1
		ORDER BY created_at DESC
		LIMIT $2
	`, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []models.TerritoryEvent
	for rows.Next() {
		var e models.TerritoryEvent
		var created time.Time
		var meta []byte
		if err := rows.Scan(&e.ID, &e.EventType, &e.H3ParentHex, &e.ActorID, &e.VictimID, &e.ActivityID, &meta, &created); err != nil {
			return nil, err
		}
		_ = json.Unmarshal(meta, &e.Metadata)
		e.CreatedAt = created.UTC().Format(time.RFC3339)
		events = append(events, e)
	}
	if events == nil {
		events = []models.TerritoryEvent{}
	}
	return events, rows.Err()
}

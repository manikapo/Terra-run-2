package services

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/territory-run/api/internal/anticheat"
	"github.com/territory-run/api/internal/config"
	"github.com/territory-run/api/internal/h3util"
	"github.com/territory-run/api/internal/models"
)

type ActivityService struct {
	db  *pgxpool.Pool
	cfg config.Config
}

func NewActivityService(db *pgxpool.Pool, cfg config.Config) *ActivityService {
	return &ActivityService{db: db, cfg: cfg}
}

func (s *ActivityService) Create(ctx context.Context, userID uuid.UUID, req models.CreateActivityRequest) (uuid.UUID, error) {
	id := uuid.New()
	_, err := s.db.Exec(ctx, `
		INSERT INTO activities (id, user_id, started_at, status, idempotency_key, device_info, created_at)
		VALUES ($1, $2, $3, 'recording', $4, $5, now())
		ON CONFLICT (idempotency_key) DO NOTHING
	`, id, userID, req.StartedAt, req.IdempotencyKey, req.DeviceInfo)
	if err != nil {
		return uuid.Nil, err
	}

	// If idempotency key existed, fetch existing id
	if req.IdempotencyKey != "" {
		var existing uuid.UUID
		err = s.db.QueryRow(ctx, `SELECT id FROM activities WHERE idempotency_key = $1`, req.IdempotencyKey).Scan(&existing)
		if err == nil {
			return existing, nil
		}
	}
	return id, nil
}

func (s *ActivityService) AppendPoints(ctx context.Context, activityID, userID uuid.UUID, points []models.GPSPoint) error {
	data, err := json.Marshal(points)
	if err != nil {
		return err
	}
	tag, err := s.db.Exec(ctx, `
		UPDATE activities
		SET points_json = COALESCE(points_json, '[]'::jsonb) || $1::jsonb
		WHERE id = $2 AND user_id = $3 AND status = 'recording'
	`, data, activityID, userID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("activity not found or not recording")
	}
	return nil
}

func (s *ActivityService) Complete(ctx context.Context, activityID, userID uuid.UUID, req models.CompleteActivityRequest) (*models.CompleteActivityResponse, error) {
	points, err := s.loadPoints(ctx, activityID, userID)
	if err != nil {
		return nil, err
	}

	filtered := anticheat.FilterPoints(points, s.cfg.MaxSpeedMPS, s.cfg.MaxAccuracyMeters)
	check := anticheat.Analyze(filtered, s.cfg.MaxSpeedMPS, s.cfg.MaxAccuracyMeters, s.cfg.MinActivityDurationSec)

	distance := req.DistanceM
	if distance <= 0 {
		distance = anticheat.DistanceM(filtered)
	}
	duration := req.DurationSec
	if duration <= 0 {
		duration = anticheat.DurationSec(filtered)
	}

	routePoints := make([]struct{ Lat, Lon float64 }, len(filtered))
	for i, p := range filtered {
		routePoints[i] = struct{ Lat, Lon float64 }{p.Lat, p.Lon}
	}

	cells, err := h3util.CellsFromRoute(routePoints, s.cfg.H3CaptureResolution)
	if err != nil {
		return nil, err
	}

	captureResult := &models.CaptureResult{H3Cells: cells, TotalCells: len(cells)}

	if check.Status == models.ActivityVerified {
		newCells, score, err := s.captureTerritories(ctx, userID, cells)
		if err != nil {
			return nil, err
		}
		captureResult.NewCells = newCells
		captureResult.CaptureScore = score
	}

	filteredJson, err := json.Marshal(filtered)
	if err != nil {
		return nil, err
	}

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, `
		UPDATE activities SET
			ended_at = $1,
			distance_m = $2,
			duration_s = $3,
			status = $4,
			trust_score = $5,
			capture_score = $6,
			h3_cells = $7,
			points_json = $8
		WHERE id = $9 AND user_id = $10
	`, req.EndedAt, distance, duration, check.Status, check.TrustScore,
		captureResult.CaptureScore, cells, filteredJson, activityID, userID)
	if err != nil {
		return nil, err
	}

	if err := s.upsertUserStats(ctx, tx, userID); err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	resp := &models.CompleteActivityResponse{
		ActivityID:    activityID,
		Status:        check.Status,
		CaptureResult: captureResult,
		TrustScore:    check.TrustScore,
	}
	if check.Reason != "" && check.Status != models.ActivityVerified {
		resp.Message = check.Reason
	}
	return resp, nil
}

func (s *ActivityService) loadPoints(ctx context.Context, activityID, userID uuid.UUID) ([]models.GPSPoint, error) {
	var raw []byte
	err := s.db.QueryRow(ctx, `
		SELECT COALESCE(points_json, '[]'::jsonb) FROM activities
		WHERE id = $1 AND user_id = $2
	`, activityID, userID).Scan(&raw)
	if err != nil {
		return nil, err
	}
	var points []models.GPSPoint
	if err := json.Unmarshal(raw, &points); err != nil {
		return nil, err
	}
	return points, nil
}

func (s *ActivityService) captureTerritories(ctx context.Context, userID uuid.UUID, cells []string) (int, int, error) {
	newCells := 0
	score := 0

	for _, cellHex := range cells {
		parent, err := h3util.ParentCell(cellHex, s.cfg.H3CaptureResolution-1)
		if err != nil {
			parent = cellHex
		}

		h3Int, err := h3util.CellToInt64(cellHex)
		if err != nil {
			h3Int = 0
		}
		tag, err := s.db.Exec(ctx, `
			INSERT INTO territory_cells (h3_index, h3_index_hex, parent_h3_hex, owner_id, status, captured_at, capture_count)
			VALUES ($1, $2, $3, $4, 'OWNED', now(), 1)
			ON CONFLICT (h3_index_hex) DO UPDATE SET
				owner_id = EXCLUDED.owner_id,
				captured_at = now(),
				capture_count = territory_cells.capture_count + 1,
				status = 'OWNED'
			WHERE territory_cells.owner_id IS NULL OR territory_cells.owner_id = EXCLUDED.owner_id
		`, h3Int, cellHex, parent, userID)
		if err != nil {
			return newCells, score, err
		}
		if tag.RowsAffected() > 0 {
			// Count as new capture if we owned null or same user refresh
			newCells++
			score += 10
			_, _ = s.db.Exec(ctx, `
				INSERT INTO territory_events (id, event_type, h3_parent_hex, actor_id, metadata, created_at)
				VALUES ($1, 'CAPTURED', $2, $3, $4, now())
			`, uuid.New(), parent, userID, map[string]interface{}{"h3_cell": cellHex})
		}
	}
	return newCells, score, nil
}

func (s *ActivityService) upsertUserStats(ctx context.Context, tx pgx.Tx, userID uuid.UUID) error {
	_, err := tx.Exec(ctx, `
		INSERT INTO user_territory_stats (user_id, cells_owned, territories_captured, total_capture_score, activity_count, updated_at)
		SELECT
			$1,
			(SELECT COUNT(*)::int FROM territory_cells WHERE owner_id = $1),
			(SELECT COUNT(DISTINCT parent_h3_hex)::int FROM territory_cells WHERE owner_id = $1),
			(SELECT COALESCE(SUM(capture_score), 0)::int FROM activities WHERE user_id = $1 AND status = 'verified'),
			(SELECT COUNT(*)::int FROM activities WHERE user_id = $1 AND status IN ('verified', 'completed')),
			now()
		ON CONFLICT (user_id) DO UPDATE SET
			cells_owned = EXCLUDED.cells_owned,
			territories_captured = EXCLUDED.territories_captured,
			total_capture_score = EXCLUDED.total_capture_score,
			activity_count = EXCLUDED.activity_count,
			updated_at = now()
	`, userID)
	return err
}

func cellToInt64(hex string) int64 {
	cell, err := h3util.ParseCellHex(hex)
	if err != nil {
		return 0
	}
	return int64(cell)
}

type TerritoryService struct {
	db  *pgxpool.Pool
	cfg config.Config
}

func NewTerritoryService(db *pgxpool.Pool, cfg config.Config) *TerritoryService {
	return &TerritoryService{db: db, cfg: cfg}
}

func (s *TerritoryService) TileGeoJSON(ctx context.Context, tileHex string) (*models.TerritoryGeoJSON, error) {
	children, err := h3util.ChildrenAtResolution(tileHex, s.cfg.H3CaptureResolution)
	if err != nil {
		return nil, err
	}

	rows, err := s.db.Query(ctx, `
		SELECT h3_index_hex, owner_id, status FROM territory_cells
		WHERE h3_index_hex = ANY($1)
	`, children)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ownerMap := make(map[string]struct {
		OwnerID uuid.UUID
		Status  string
	})
	for rows.Next() {
		var hex string
		var ownerID uuid.UUID
		var status string
		if err := rows.Scan(&hex, &ownerID, &status); err != nil {
			return nil, err
		}
		ownerMap[hex] = struct {
			OwnerID uuid.UUID
			Status  string
		}{ownerID, status}
	}

	features := make([]models.TerritoryFeature, 0, len(ownerMap))
	for hex, info := range ownerMap {
		geom, err := h3util.CellBoundaryGeoJSON(hex)
		if err != nil {
			continue
		}
		features = append(features, models.TerritoryFeature{
			Type: "Feature",
			Properties: map[string]interface{}{
				"h3":      hex,
				"owner":   info.OwnerID.String(),
				"status":  info.Status,
			},
			Geometry: geom,
		})
	}

	return &models.TerritoryGeoJSON{Type: "FeatureCollection", Features: features}, nil
}

type UserService struct {
	db *pgxpool.Pool
}

func NewUserService(db *pgxpool.Pool) *UserService {
	return &UserService{db: db}
}

func (s *UserService) EnsureProfile(ctx context.Context, userID uuid.UUID, email string) error {
	username := email
	displayName := email
	if email == "" || strings.Contains(email, "@guest") || strings.HasSuffix(email, ".local") {
		short := strings.ReplaceAll(userID.String(), "-", "")
		if len(short) > 8 {
			short = short[:8]
		}
		username = "guest_" + short
		displayName = "Guest Runner"
		email = "guest_" + userID.String() + "@territory.run"
	} else if at := indexByte(email, '@'); at > 0 {
		username = email[:at]
		displayName = username
	}
	_, err := s.db.Exec(ctx, `
		INSERT INTO users (id, username, display_name, email, created_at)
		VALUES ($1, $2, $3, $4, now())
		ON CONFLICT (id) DO UPDATE SET
			display_name = EXCLUDED.display_name,
			email = EXCLUDED.email
	`, userID, username, displayName, email)
	return err
}

func indexByte(s string, c byte) int {
	for i := 0; i < len(s); i++ {
		if s[i] == c {
			return i
		}
	}
	return -1
}

func (s *UserService) Profile(ctx context.Context, userID uuid.UUID) (*models.UserProfile, error) {
	var p models.UserProfile
	err := s.db.QueryRow(ctx, `
		SELECT u.id, u.username, u.display_name, COALESCE(u.avatar_url, ''),
			COALESCE(s.cells_owned, 0), COALESCE(s.territories_captured, 0),
			COALESCE(s.total_capture_score, 0), COALESCE(s.activity_count, 0)
		FROM users u
		LEFT JOIN user_territory_stats s ON s.user_id = u.id
		WHERE u.id = $1
	`, userID).Scan(&p.UserID, &p.Username, &p.DisplayName, &p.AvatarURL,
		&p.CellsOwned, &p.TerritoriesCaptured, &p.TotalCaptureScore, &p.ActivityCount)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (s *UserService) ListActivities(ctx context.Context, userID uuid.UUID, limit int) ([]map[string]interface{}, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	rows, err := s.db.Query(ctx, `
		SELECT id, started_at, ended_at, distance_m, duration_s, capture_score, status, created_at
		FROM activities WHERE user_id = $1
		ORDER BY created_at DESC LIMIT $2
	`, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []map[string]interface{}
	for rows.Next() {
		var id uuid.UUID
		var started, ended *time.Time
		var distance float64
		var duration, captureScore int
		var status string
		var created time.Time
		if err := rows.Scan(&id, &started, &ended, &distance, &duration, &captureScore, &status, &created); err != nil {
			return nil, err
		}
		out = append(out, map[string]interface{}{
			"id":            id,
			"started_at":    started,
			"ended_at":      ended,
			"distance_m":    distance,
			"duration_s":    duration,
			"capture_score": captureScore,
			"status":        status,
			"created_at":    created,
		})
	}
	return out, nil
}

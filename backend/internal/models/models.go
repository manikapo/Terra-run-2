package models

import (
	"time"

	"github.com/google/uuid"
)

type ActivityStatus string

const (
	ActivityRecording ActivityStatus = "recording"
	ActivityCompleted ActivityStatus = "completed"
	ActivityVerified  ActivityStatus = "verified"
	ActivityRejected  ActivityStatus = "rejected"
	ActivityQuarantined ActivityStatus = "quarantined"
)

type GPSPoint struct {
	Lat       float64 `json:"lat"`
	Lon       float64 `json:"lon"`
	Timestamp int64   `json:"ts"`
	Accuracy  float64 `json:"acc"`
	Speed     float64 `json:"speed,omitempty"`
}

type CreateActivityRequest struct {
	StartedAt   time.Time `json:"started_at"`
	DeviceInfo  map[string]interface{} `json:"device_info,omitempty"`
	IdempotencyKey string `json:"idempotency_key"`
}

type UploadPointsRequest struct {
	Points []GPSPoint `json:"points"`
}

type CompleteActivityRequest struct {
	EndedAt      time.Time `json:"ended_at"`
	DistanceM    float64   `json:"distance_m"`
	DurationSec  int       `json:"duration_s"`
	H3CellsHint  []string  `json:"h3_cells_hint,omitempty"`
}

type CaptureResult struct {
	NewCells      int   `json:"new_cells"`
	TotalCells    int   `json:"total_cells"`
	CaptureScore  int   `json:"capture_score"`
	H3Cells       []string `json:"h3_cells"`
}

type CompleteActivityResponse struct {
	ActivityID    uuid.UUID     `json:"activity_id"`
	Status        ActivityStatus `json:"status"`
	CaptureResult *CaptureResult `json:"capture_result,omitempty"`
	TrustScore    int           `json:"trust_score"`
	Message       string        `json:"message,omitempty"`
}

type UserProfile struct {
	UserID            uuid.UUID `json:"user_id"`
	Username          string    `json:"username"`
	DisplayName       string    `json:"display_name"`
	AvatarURL         string    `json:"avatar_url,omitempty"`
	CellsOwned        int       `json:"cells_owned"`
	TerritoriesCaptured int     `json:"territories_captured"`
	TotalCaptureScore int       `json:"total_capture_score"`
	ActivityCount     int       `json:"activity_count"`
}

type TerritoryFeature struct {
	Type       string                 `json:"type"`
	Properties map[string]interface{} `json:"properties"`
	Geometry   map[string]interface{} `json:"geometry"`
}

type TerritoryGeoJSON struct {
	Type     string             `json:"type"`
	Features []TerritoryFeature `json:"features"`
}

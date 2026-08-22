package models

import "github.com/google/uuid"

const (
	ScoreNeutralCapture = 10
	ScoreDefend         = 2
	ScoreSteal          = 15
)

type CellOutcome string

const (
	OutcomeNeutral CellOutcome = "neutral"
	OutcomeDefend  CellOutcome = "defend"
	OutcomeSteal   CellOutcome = "steal"
)

// ClassifyCell decides capture vs steal vs defend from existing ownership.
// existingOwner nil/uuid.Nil means unowned.
func ClassifyCell(existingOwner *uuid.UUID, actor uuid.UUID) CellOutcome {
	if existingOwner == nil || *existingOwner == uuid.Nil {
		return OutcomeNeutral
	}
	if *existingOwner == actor {
		return OutcomeDefend
	}
	return OutcomeSteal
}

func ScoreForOutcome(o CellOutcome) int {
	switch o {
	case OutcomeDefend:
		return ScoreDefend
	case OutcomeSteal:
		return ScoreSteal
	default:
		return ScoreNeutralCapture
	}
}

type LeaderboardEntry struct {
	Rank                int       `json:"rank"`
	UserID              uuid.UUID `json:"user_id"`
	Username            string    `json:"username"`
	DisplayName         string    `json:"display_name"`
	CellsOwned          int       `json:"cells_owned"`
	TerritoriesCaptured int       `json:"territories_captured"`
	TerritoriesStolen   int       `json:"territories_stolen"`
	TotalCaptureScore   int64     `json:"total_capture_score"`
}

type LeaderboardResponse struct {
	Scope   string             `json:"scope"`
	Entries []LeaderboardEntry `json:"entries"`
}

type DecayResult struct {
	Neutralized int `json:"neutralized"`
}

type TerritoryEvent struct {
	ID          uuid.UUID              `json:"id"`
	EventType   string                 `json:"event_type"`
	H3ParentHex string                 `json:"h3_parent_hex,omitempty"`
	ActorID     *uuid.UUID             `json:"actor_id,omitempty"`
	VictimID    *uuid.UUID             `json:"victim_id,omitempty"`
	ActivityID  *uuid.UUID             `json:"activity_id,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
	CreatedAt   string                 `json:"created_at"`
}

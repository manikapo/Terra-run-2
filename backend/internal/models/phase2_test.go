package models

import (
	"testing"

	"github.com/google/uuid"
)

func TestClassifyCell(t *testing.T) {
	actor := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	rival := uuid.MustParse("22222222-2222-2222-2222-222222222222")

	if got := ClassifyCell(nil, actor); got != OutcomeNeutral {
		t.Fatalf("unowned: got %s", got)
	}
	nilID := uuid.Nil
	if got := ClassifyCell(&nilID, actor); got != OutcomeNeutral {
		t.Fatalf("nil uuid: got %s", got)
	}
	if got := ClassifyCell(&actor, actor); got != OutcomeDefend {
		t.Fatalf("own: got %s", got)
	}
	if got := ClassifyCell(&rival, actor); got != OutcomeSteal {
		t.Fatalf("rival: got %s", got)
	}
}

func TestScoreForOutcome(t *testing.T) {
	if ScoreForOutcome(OutcomeNeutral) != 10 {
		t.Fatal("neutral")
	}
	if ScoreForOutcome(OutcomeDefend) != 2 {
		t.Fatal("defend")
	}
	if ScoreForOutcome(OutcomeSteal) != 15 {
		t.Fatal("steal")
	}
}

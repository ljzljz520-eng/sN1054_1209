package validation

import (
	"testing"

	"example.com/familyitinerary/internal/model"
)

func TestIntakeIssues(t *testing.T) {
	in := model.Intake{Family: model.NewFamily("f", "A", 1, ""), Itinerary: model.NewItinerary("i", "f", "Trip", "Sanya", "2026-10-01", "2026-10-03"), Preferences: []string{"unknown"}}
	issues := IntakeIssues(in)
	if len(issues) != 2 {
		t.Fatalf("got %d issues", len(issues))
	}
	if Summarize(issues) == "" {
		t.Fatal("summary missing")
	}
}

func TestMessageIssue(t *testing.T) {
	if MessageIssue(model.ChatMessage{ID: "m", Itinerary: "i", Sequence: 1, Sender: "parent", Body: " hello"}) == "" {
		t.Fatal("whitespace message accepted")
	}
	if DateSuggestionIssue(model.DateSuggestion{ID: "s", Itinerary: "i", Suggested: "2026-10-01", Reason: "crowds"}) != "" {
		t.Fatal("valid suggestion rejected")
	}
}

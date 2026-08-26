package model

import "testing"

func TestRules(t *testing.T) {
	if ValidFamily(Family{}) {
		t.Fatal("empty family accepted")
	}
	if ValidItinerary(Itinerary{ID: "i", FamilyID: "f", Destination: "x", StartDate: "2026-09-03", EndDate: "2026-09-01"}) {
		t.Fatal("inverted dates accepted")
	}
	if ValidActivity(Activity{ID: "a", Itinerary: "i", Title: "x", Day: 0, Duration: 1}) {
		t.Fatal("invalid day accepted")
	}
	if NormalizeStatus("unknown") != "pending" || !IsFailure("failed") || PreferenceAllowed("unknown") {
		t.Fatal("status or preference rule failed")
	}
}

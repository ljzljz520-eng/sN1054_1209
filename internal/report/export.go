package report

import (
	"encoding/json"
	"fmt"

	"example.com/familyitinerary/internal/advisor"
)

type View struct {
	ItineraryID string `json:"itinerary_id"`
	Title       string `json:"title"`
	Destination string `json:"destination"`
	Status      string `json:"status"`
	Summary     string `json:"summary"`
	SafeCount   int    `json:"child_safe_activities"`
}

func FromReport(r advisor.Report) View {
	return View{ItineraryID: r.Intake.Itinerary.ID, Title: r.Intake.Itinerary.Title, Destination: r.Intake.Itinerary.Destination, Status: r.Intake.Itinerary.Status, Summary: r.Summary, SafeCount: r.ChildSafeActivities()}
}

func JSON(r advisor.Report) ([]byte, error) { return json.Marshal(FromReport(r)) }

func Plain(r advisor.Report) string {
	return fmt.Sprintf("%s\n%s", r.Intake.Itinerary.Title, r.Summary)
}

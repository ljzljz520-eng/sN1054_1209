package report

import (
	"testing"

	"example.com/familyitinerary/internal/advisor"
	"example.com/familyitinerary/internal/model"
	"example.com/familyitinerary/internal/store"
)

func TestJSONExport(t *testing.T) {
	repository, err := store.OpenMemory()
	if err != nil {
		t.Fatal(err)
	}
	defer repository.Close()
	family := model.NewFamily("f-report", "Sun", 1, "")
	itinerary := model.NewItinerary("i-report", family.ID, "Lake", "Hangzhou", "2026-06-01", "2026-06-02")
	intake := advisor.NewIntakeService(repository)
	if err = intake.Register(model.Intake{Family: family, Itinerary: itinerary, Activities: []model.Activity{model.NewActivity("a-report", itinerary.ID, 1, "Boat", "park", 60, true)}}); err != nil {
		t.Fatal(err)
	}
	result, err := advisor.BuildReport(repository, itinerary.ID)
	if err != nil {
		t.Fatal(err)
	}
	payload, err := JSON(result)
	if err != nil || len(payload) == 0 {
		t.Fatalf("payload=%s err=%v", payload, err)
	}
	if Plain(result) == "" {
		t.Fatal("plain report empty")
	}
}

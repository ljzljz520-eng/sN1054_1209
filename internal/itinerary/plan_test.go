package itinerary

import (
	"testing"

	"example.com/familyitinerary/internal/model"
)

func TestPlannerCapacityAndOrdering(t *testing.T) {
	planner := NewPlanner()
	if err := planner.AddActivity(model.NewActivity("b", "it", 2, "Park", "park", 90, true)); err != nil {
		t.Fatal(err)
	}
	if err := planner.AddActivity(model.NewActivity("a", "it", 1, "Aquarium", "indoor", 100, true)); err != nil {
		t.Fatal(err)
	}
	if planner.DayCount("it") != 2 || len(planner.Activities("it")) != 2 {
		t.Fatal("planner state incorrect")
	}
	if err := planner.ValidateCapacity("it", 50); err == nil {
		t.Fatal("capacity overflow accepted")
	}
}

func TestScheduleChecksChildren(t *testing.T) {
	item := model.NewItinerary("i", "f", "Trip", "Chengdu", "2026-01-01", "2026-01-02")
	schedule := BuildSchedule(item, []model.Activity{model.NewActivity("a", "i", 1, "Zoo", "outdoor", 90, false)})
	if err := schedule.CheckChildFit(1); err == nil {
		t.Fatal("unsafe schedule accepted")
	}
	if schedule.ChildSafeRatio() != 0 {
		t.Fatal("ratio incorrect")
	}
}

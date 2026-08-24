package store

import (
	"path/filepath"
	"testing"

	"example.com/familyitinerary/internal/model"
)

func TestPersistenceSurvivesReopen(t *testing.T) {
	path := filepath.Join(t.TempDir(), "family.db")
	family := model.NewFamily("fam-reopen", "Chen", 2, "stroller")
	itinerary := model.NewItinerary("it-reopen", family.ID, "Museum route", "Beijing", "2026-11-01", "2026-11-03")
	first, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	if err = first.SaveFamily(family); err != nil {
		t.Fatal(err)
	}
	if err = first.SaveItinerary(itinerary); err != nil {
		t.Fatal(err)
	}
	if err = first.SaveActivity(model.NewActivity("act-reopen", itinerary.ID, 1, "Museum", "indoor", 80, true)); err != nil {
		t.Fatal(err)
	}
	if err = first.Close(); err != nil {
		t.Fatal(err)
	}
	second, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer second.Close()
	loaded, err := second.GetItinerary(itinerary.ID)
	if err != nil {
		t.Fatal(err)
	}
	if loaded.Destination != "Beijing" {
		t.Fatalf("got %s", loaded.Destination)
	}
	activities, err := second.ListActivities(itinerary.ID)
	if err != nil || len(activities) != 1 {
		t.Fatalf("activities=%d err=%v", len(activities), err)
	}
}

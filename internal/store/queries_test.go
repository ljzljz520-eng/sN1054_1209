package store

import (
	"testing"

	"example.com/familyitinerary/internal/model"
)

func TestStoreQueries(t *testing.T) {
	repository, err := OpenMemory()
	if err != nil {
		t.Fatal(err)
	}
	defer repository.Close()
	family := model.NewFamily("fam-q", "Wang", 1, "")
	itinerary := model.NewItinerary("it-q", family.ID, "Garden", "Suzhou", "2026-12-01", "2026-12-02")
	if err = repository.SaveFamily(family); err != nil {
		t.Fatal(err)
	}
	if err = repository.SaveItinerary(itinerary); err != nil {
		t.Fatal(err)
	}
	if err = repository.SavePreferences(itinerary.ID, []string{"short-transfers", "indoor-backup"}); err != nil {
		t.Fatal(err)
	}
	if got, _ := repository.LoadPreferences(itinerary.ID); len(got) != 2 {
		t.Fatal("preferences were not stored")
	}
	if err = repository.SaveSuggestion(model.NewSuggestion("s-q", itinerary.ID, "2026-12-01", "school break")); err != nil {
		t.Fatal(err)
	}
	if err = repository.AcceptSuggestion("s-q"); err != nil {
		t.Fatal(err)
	}
	items, err := repository.ListSuggestions(itinerary.ID)
	if err != nil || len(items) != 1 || !items[0].Accepted {
		t.Fatal("suggestion query failed")
	}
}

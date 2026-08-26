package model

import "testing"

func TestConstructors(t *testing.T) {
	family := NewFamily("fam-1", "Lin", 2, "nap after lunch")
	itinerary := NewItinerary("it-1", family.ID, "Ocean weekend", "Qingdao", "2026-09-01", "2026-09-03")
	activity := NewActivity("act-1", itinerary.ID, 1, "Aquarium", "indoor", 90, true)
	message := NewMessage("msg-1", itinerary.ID, 1, "parent", "Can we move the beach day?")
	suggestion := NewSuggestion("sug-1", itinerary.ID, "2026-09-02", "lower crowd")
	if family.Children != 2 || itinerary.Status != "draft" || !activity.ChildSafe || message.Status != "pending" || suggestion.Accepted {
		t.Fatal("constructors did not set domain defaults")
	}
}

package model

import "fmt"

var itineraryTransitions = map[string]map[string]bool{
	"draft":    {"accepted": true, "rejected": true},
	"accepted": {"draft": true, "rejected": true},
	"rejected": {"draft": true},
}

func CanTransitionItinerary(from, to string) bool {
	if from == to {
		return true
	}
	if next, ok := itineraryTransitions[from]; ok {
		return next[to]
	}
	return false
}

func TransitionItinerary(item Itinerary, to string) (Itinerary, error) {
	if !CanTransitionItinerary(NormalizeStatus(item.Status), NormalizeStatus(to)) {
		return item, fmt.Errorf("cannot move itinerary from %s to %s", item.Status, to)
	}
	item.Status = NormalizeStatus(to)
	return item, nil
}

func MessageLabel(message ChatMessage) string {
	if IsFailure(message.Status) {
		return "Needs attention: " + message.Failure
	}
	if message.Status == "sent" {
		return "Delivered: " + message.Body
	}
	return "Pending: " + message.Body
}

func SortMessages(messages []ChatMessage) []ChatMessage {
	ordered := append([]ChatMessage(nil), messages...)
	for index := 1; index < len(ordered); index++ {
		current := ordered[index]
		position := index - 1
		for position >= 0 && ordered[position].Sequence > current.Sequence {
			ordered[position+1] = ordered[position]
			position--
		}
		ordered[position+1] = current
	}
	return ordered
}

func StatusSummary(messages []ChatMessage) (sent, failed, pending int) {
	for _, message := range messages {
		switch message.Status {
		case "sent":
			sent++
		case "failed", "rejected":
			failed++
		default:
			pending++
		}
	}
	return sent, failed, pending
}

func ItineraryReady(item Itinerary, activities []Activity) bool {
	if !ValidItinerary(item) || len(activities) == 0 {
		return false
	}
	for _, activity := range activities {
		if !ValidActivity(activity) || activity.Itinerary != item.ID {
			return false
		}
	}
	return true
}

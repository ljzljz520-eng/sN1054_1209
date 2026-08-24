package advisor

import (
	"fmt"
	"sort"

	"example.com/familyitinerary/internal/model"
	"example.com/familyitinerary/internal/store"
)

type Dashboard struct {
	Families       int
	Itineraries    int
	Messages       int
	FailedMessages int
	AcceptedDates  int
}

func BuildDashboard(repository *store.Store, familyID string) (Dashboard, error) {
	families, err := repository.ListFamilies()
	if err != nil {
		return Dashboard{}, err
	}
	itineraries, err := repository.ListItineraries(familyID)
	if err != nil {
		return Dashboard{}, err
	}
	dashboard := Dashboard{Families: len(families), Itineraries: len(itineraries)}
	for _, item := range itineraries {
		messages, listErr := repository.ListMessages(item.ID)
		if listErr != nil {
			return Dashboard{}, listErr
		}
		dashboard.Messages += len(messages)
		for _, message := range messages {
			if model.IsFailure(message.Status) {
				dashboard.FailedMessages++
			}
		}
		suggestions, listErr := repository.ListSuggestions(item.ID)
		if listErr != nil {
			return Dashboard{}, listErr
		}
		for _, suggestion := range suggestions {
			if suggestion.Accepted {
				dashboard.AcceptedDates++
			}
		}
	}
	return dashboard, nil
}

func DashboardLabel(dashboard Dashboard) string {
	return fmt.Sprintf("%d families, %d itineraries, %d messages, %d failures, %d accepted dates", dashboard.Families, dashboard.Itineraries, dashboard.Messages, dashboard.FailedMessages, dashboard.AcceptedDates)
}

func RankDestinations(items []model.Itinerary) []string {
	counts := make(map[string]int)
	for _, item := range items {
		counts[item.Destination]++
	}
	type pair struct {
		name  string
		count int
	}
	pairs := make([]pair, 0, len(counts))
	for name, count := range counts {
		pairs = append(pairs, pair{name: name, count: count})
	}
	sort.Slice(pairs, func(i, j int) bool {
		if pairs[i].count == pairs[j].count {
			return pairs[i].name < pairs[j].name
		}
		return pairs[i].count > pairs[j].count
	})
	result := make([]string, 0, len(pairs))
	for _, item := range pairs {
		result = append(result, item.name)
	}
	return result
}

func StatusBreakdown(items []model.Itinerary) map[string]int {
	breakdown := make(map[string]int)
	for _, item := range items {
		breakdown[item.Status]++
	}
	return breakdown
}

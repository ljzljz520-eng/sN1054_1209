package advisor

import (
	"sort"
	"strings"

	"example.com/familyitinerary/internal/model"
	"example.com/familyitinerary/internal/store"
)

type SearchQuery struct {
	FamilyID    string
	Destination string
	Status      string
	Text        string
}

type SearchResult struct {
	Itinerary   model.Itinerary
	MatchReason string
}

func SearchItineraries(repository *store.Store, query SearchQuery) ([]SearchResult, error) {
	items, err := repository.ListItineraries(query.FamilyID)
	if err != nil {
		return nil, err
	}
	results := make([]SearchResult, 0)
	for _, item := range items {
		if query.Destination != "" && !strings.EqualFold(item.Destination, query.Destination) {
			continue
		}
		if query.Status != "" && item.Status != query.Status {
			continue
		}
		reason := "family match"
		if query.Text != "" {
			if !strings.Contains(strings.ToLower(item.Title), strings.ToLower(query.Text)) && !strings.Contains(strings.ToLower(item.Destination), strings.ToLower(query.Text)) {
				continue
			}
			reason = "text match"
		}
		results = append(results, SearchResult{Itinerary: item, MatchReason: reason})
	}
	sort.SliceStable(results, func(i, j int) bool { return results[i].Itinerary.StartDate < results[j].Itinerary.StartDate })
	return results, nil
}

func SearchMessages(repository *store.Store, itineraryID, text string) ([]model.ChatMessage, error) {
	messages, err := repository.ListMessages(itineraryID)
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(text) == "" {
		return messages, nil
	}
	filtered := make([]model.ChatMessage, 0)
	for _, message := range messages {
		if strings.Contains(strings.ToLower(message.Body), strings.ToLower(text)) {
			filtered = append(filtered, message)
		}
	}
	return filtered, nil
}

func CountByStatus(items []SearchResult) map[string]int {
	counts := make(map[string]int)
	for _, item := range items {
		counts[item.Itinerary.Status]++
	}
	return counts
}

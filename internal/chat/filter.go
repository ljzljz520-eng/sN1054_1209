package chat

import (
	"sort"
	"strings"

	"example.com/familyitinerary/internal/model"
)

type Filter struct {
	Status          string
	Sender          string
	Contains        string
	IncludeFailures bool
}

func (f Filter) Match(message model.ChatMessage) bool {
	if f.Status != "" && message.Status != f.Status {
		return false
	}
	if f.Sender != "" && message.Sender != f.Sender {
		return false
	}
	if f.Contains != "" && !strings.Contains(strings.ToLower(message.Body), strings.ToLower(f.Contains)) {
		return false
	}
	if !f.IncludeFailures && model.IsFailure(message.Status) {
		return false
	}
	return true
}

func FilterMessages(messages []model.ChatMessage, filter Filter) []model.ChatMessage {
	items := make([]model.ChatMessage, 0)
	for _, message := range messages {
		if filter.Match(message) {
			items = append(items, message)
		}
	}
	return model.SortMessages(items)
}

func GroupByStatus(messages []model.ChatMessage) map[string][]model.ChatMessage {
	groups := make(map[string][]model.ChatMessage)
	for _, message := range messages {
		groups[message.Status] = append(groups[message.Status], message)
	}
	for status := range groups {
		sort.Slice(groups[status], func(i, j int) bool { return groups[status][i].Sequence < groups[status][j].Sequence })
	}
	return groups
}

func FailureReasons(messages []model.ChatMessage) []string {
	reasons := make([]string, 0)
	for _, message := range messages {
		if model.IsFailure(message.Status) && message.Failure != "" {
			reasons = append(reasons, message.Failure)
		}
	}
	return reasons
}

func LatestSuccessful(messages []model.ChatMessage) (model.ChatMessage, bool) {
	ordered := model.SortMessages(messages)
	for index := len(ordered) - 1; index >= 0; index-- {
		if ordered[index].Status == "sent" {
			return ordered[index], true
		}
	}
	return model.ChatMessage{}, false
}

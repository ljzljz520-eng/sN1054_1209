package chat

import (
	"fmt"
	"strings"

	"example.com/familyitinerary/internal/model"
)

type Card struct {
	ID       string
	Sequence int
	Sender   string
	Body     string
	Status   string
	Failure  string
	Tone     string
}

func ToCard(message model.ChatMessage) Card {
	tone := "neutral"
	if message.Status == "sent" {
		tone = "positive"
	}
	if model.IsFailure(message.Status) {
		tone = "negative"
	}
	return Card{ID: message.ID, Sequence: message.Sequence, Sender: message.Sender, Body: message.Body, Status: message.Status, Failure: message.Failure, Tone: tone}
}

func Cards(messages []model.ChatMessage) []Card {
	ordered := model.SortMessages(messages)
	cards := make([]Card, 0, len(ordered))
	for _, message := range ordered {
		cards = append(cards, ToCard(message))
	}
	return cards
}

func TranscriptText(messages []model.ChatMessage) string {
	lines := make([]string, 0, len(messages))
	for _, card := range Cards(messages) {
		line := fmt.Sprintf("%d %s [%s]: %s", card.Sequence, card.Sender, card.Status, card.Body)
		if card.Failure != "" {
			line += " (" + card.Failure + ")"
		}
		lines = append(lines, line)
	}
	return strings.Join(lines, "\n")
}

func CountVisibleFailures(messages []model.ChatMessage) int {
	count := 0
	for _, message := range messages {
		if model.IsFailure(message.Status) && strings.TrimSpace(message.Failure) != "" {
			count++
		}
	}
	return count
}

func ReplaceMessage(messages []model.ChatMessage, replacement model.ChatMessage) []model.ChatMessage {
	updated := append([]model.ChatMessage(nil), messages...)
	for index := range updated {
		if updated[index].ID == replacement.ID {
			updated[index] = replacement
			return updated
		}
	}
	return append(updated, replacement)
}

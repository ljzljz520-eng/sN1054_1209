package chat

import (
	"fmt"
	"strings"

	"example.com/familyitinerary/internal/model"
)

type Delivery struct {
	Message model.ChatMessage
	Visible bool
}

func PrepareUploadMessage(itineraryID string, sequence int, title string) model.ChatMessage {
	return model.NewMessage(fmt.Sprintf("upload-%s-%d", itineraryID, sequence), itineraryID, sequence, "advisor", "Uploaded itinerary: "+title)
}

func PrepareDateMessage(itineraryID string, sequence int, date string) model.ChatMessage {
	return model.NewMessage(fmt.Sprintf("date-%s-%d", itineraryID, sequence), itineraryID, sequence, "advisor", "Suggested travel date: "+date)
}

func EvaluateUpload(message model.ChatMessage, accepted bool) model.ChatMessage {
	if accepted {
		message.Status = "sent"
		message.Failure = ""
	} else {
		message.Status = "failed"
		message.Failure = "itinerary file could not be parsed"
	}
	return message
}

func EvaluateDateSuggestion(message model.ChatMessage, available bool) model.ChatMessage {
	if available {
		message.Status = "sent"
		message.Failure = ""
	} else {
		message.Status = "failed"
		message.Failure = "suggested date conflicts with school calendar"
	}
	return message
}

func VisibleMessages(messages []model.ChatMessage) []Delivery {
	deliveries := make([]Delivery, 0, len(messages))
	for _, message := range messages {
		visible := message.Body != ""
		if model.IsFailure(message.Status) && strings.TrimSpace(message.Failure) == "" {
			visible = false
		}
		deliveries = append(deliveries, Delivery{Message: message, Visible: visible})
	}
	return deliveries
}

package chat

import (
	"fmt"
	"strings"

	"example.com/familyitinerary/internal/model"
	"example.com/familyitinerary/internal/store"
)

type Engine struct{ store *store.Store }

func NewEngine(repository *store.Store) *Engine { return &Engine{store: repository} }

func (e *Engine) Transcript(itineraryID string) (model.ChatTranscript, error) {
	messages, err := e.store.ListMessages(itineraryID)
	if err != nil {
		return model.ChatTranscript{}, err
	}
	return model.ChatTranscript{Itinerary: itineraryID, Messages: messages}, nil
}

func (e *Engine) Send(itineraryID, sender, body string) (model.ChatMessage, error) {
	if strings.TrimSpace(body) == "" {
		return model.ChatMessage{}, fmt.Errorf("message body is required")
	}
	sequence, err := e.store.NextMessageSequence(itineraryID)
	if err != nil {
		return model.ChatMessage{}, err
	}
	message := model.NewMessage(fmt.Sprintf("msg-%s-%d", itineraryID, sequence), itineraryID, sequence, sender, body)
	if err := e.store.SaveMessage(message); err != nil {
		return model.ChatMessage{}, err
	}
	return message, nil
}

func (e *Engine) SetStatus(messageID, status, failure string) error {
	if model.IsFailure(status) && failure == "" {
		return fmt.Errorf("failure reason is required")
	}
	return e.store.UpdateMessageStatus(messageID, status, failure)
}

func (e *Engine) LastMessage(itineraryID string) (model.ChatMessage, error) {
	transcript, err := e.Transcript(itineraryID)
	if err != nil {
		return model.ChatMessage{}, err
	}
	if len(transcript.Messages) == 0 {
		return model.ChatMessage{}, fmt.Errorf("no messages")
	}
	return transcript.Messages[len(transcript.Messages)-1], nil
}

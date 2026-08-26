package chat

import (
	"fmt"

	"example.com/familyitinerary/internal/model"
	"example.com/familyitinerary/internal/store"
)

type Conversation struct{ Engine *Engine }

func NewConversation(repository *store.Store) *Conversation {
	return &Conversation{Engine: NewEngine(repository)}
}

func (c *Conversation) UploadItinerary(itineraryID, title string, accepted bool) (model.ChatMessage, error) {
	sequence, err := c.Engine.store.NextMessageSequence(itineraryID)
	if err != nil {
		return model.ChatMessage{}, err
	}
	message := EvaluateUpload(PrepareUploadMessage(itineraryID, sequence, title), accepted)
	if err := c.Engine.store.SaveMessage(message); err != nil {
		return model.ChatMessage{}, err
	}
	return message, nil
}

func (c *Conversation) SendDateSuggestion(itineraryID, date string, available bool) (model.ChatMessage, error) {
	sequence, err := c.Engine.store.NextMessageSequence(itineraryID)
	if err != nil {
		return model.ChatMessage{}, err
	}
	message := EvaluateDateSuggestion(PrepareDateMessage(itineraryID, sequence, date), available)
	previous, previousErr := c.Engine.LastMessage(itineraryID)
	if previousErr == nil && !available {
		message.Status = previous.Status
	}
	if message.Status == "" {
		return model.ChatMessage{}, fmt.Errorf("message status was not assigned")
	}
	if err := c.Engine.store.SaveMessage(message); err != nil {
		return model.ChatMessage{}, err
	}
	return message, nil
}

func (c *Conversation) History(itineraryID string) ([]Delivery, error) {
	transcript, err := c.Engine.Transcript(itineraryID)
	if err != nil {
		return nil, err
	}
	return VisibleMessages(transcript.Messages), nil
}

func (c *Conversation) MarkAccepted(messageID string) error {
	return c.Engine.SetStatus(messageID, "sent", "")
}

func (c *Conversation) MarkFailed(messageID, reason string) error {
	return c.Engine.SetStatus(messageID, "failed", reason)
}

package advisor

import (
	"fmt"
	"strings"

	"example.com/familyitinerary/internal/chat"
	"example.com/familyitinerary/internal/model"
	"example.com/familyitinerary/internal/store"
)

type ConversationService struct {
	conversation *chat.Conversation
	store        *store.Store
}

func NewConversationService(repository *store.Store) *ConversationService {
	return &ConversationService{conversation: chat.NewConversation(repository), store: repository}
}

func (s *ConversationService) Upload(itineraryID, title string, accepted bool) (model.ChatMessage, error) {
	if strings.TrimSpace(title) == "" {
		return model.ChatMessage{}, fmt.Errorf("title is required")
	}
	return s.conversation.UploadItinerary(itineraryID, title, accepted)
}

func (s *ConversationService) SuggestDate(itineraryID, date string, available bool) (model.ChatMessage, error) {
	if date == "" {
		return model.ChatMessage{}, fmt.Errorf("date is required")
	}
	return s.conversation.SendDateSuggestion(itineraryID, date, available)
}

func (s *ConversationService) History(itineraryID string) ([]chat.Delivery, error) {
	return s.conversation.History(itineraryID)
}

func (s *ConversationService) AcceptDate(suggestionID string) error {
	return s.store.AcceptSuggestion(suggestionID)
}

func (s *ConversationService) AddSuggestion(suggestion model.DateSuggestion) error {
	if suggestion.Reason == "" {
		return fmt.Errorf("reason is required")
	}
	return s.store.SaveSuggestion(suggestion)
}

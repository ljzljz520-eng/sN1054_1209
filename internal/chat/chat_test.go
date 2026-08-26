package chat

import (
	"testing"

	"example.com/familyitinerary/internal/model"
	"example.com/familyitinerary/internal/store"
)

func TestItineraryChatRetainsStatus(t *testing.T) {
	repository, err := store.OpenMemory()
	if err != nil {
		t.Fatal(err)
	}
	defer repository.Close()
	conversation := NewConversation(repository)
	if _, err = conversation.UploadItinerary("it-chat", "Family upload", true); err != nil {
		t.Fatal(err)
	}
	failed, err := conversation.SendDateSuggestion("it-chat", "2026-09-12", false)
	if err != nil {
		t.Fatal(err)
	}
	if failed.Status != "failed" {
		t.Fatalf("new message status=%s", failed.Status)
	}
	if failed.Failure != "suggested date conflicts with school calendar" {
		t.Fatalf("failure=%s", failed.Failure)
	}
	history, err := conversation.History("it-chat")
	if err != nil {
		t.Fatal(err)
	}
	if len(history) != 2 {
		t.Fatalf("history=%d", len(history))
	}
	if history[0].Message.Status != "sent" {
		t.Fatalf("upload status=%s", history[0].Message.Status)
	}
	if history[1].Message.Status != "failed" || history[1].Message.Failure == "" {
		t.Fatal("date error was not retained")
	}
}

func TestChatVisibility(t *testing.T) {
	items := VisibleMessages([]model.ChatMessage{{ID: "m1", Body: "ok", Status: "sent"}, {ID: "m2", Body: "bad", Status: "failed", Failure: "reason"}, {ID: "m3", Body: "bad", Status: "failed"}})
	if len(items) != 3 || !items[0].Visible || !items[1].Visible || items[2].Visible {
		t.Fatal("visibility rule failed")
	}
}

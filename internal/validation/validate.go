package validation

import (
	"fmt"
	"strings"

	"example.com/familyitinerary/internal/model"
)

type Issue struct {
	Field   string
	Message string
}

func IntakeIssues(in model.Intake) []Issue {
	issues := make([]Issue, 0)
	if !model.ValidFamily(in.Family) {
		issues = append(issues, Issue{Field: "family", Message: "family details are incomplete"})
	}
	if !model.ValidItinerary(in.Itinerary) {
		issues = append(issues, Issue{Field: "itinerary", Message: "itinerary dates or identity are invalid"})
	}
	if len(in.Activities) == 0 {
		issues = append(issues, Issue{Field: "activities", Message: "at least one activity is required"})
	}
	for index, activity := range in.Activities {
		if !model.ValidActivity(activity) {
			issues = append(issues, Issue{Field: fmt.Sprintf("activities[%d]", index), Message: "activity is invalid"})
		}
	}
	for index, preference := range in.Preferences {
		if !model.PreferenceAllowed(preference) {
			issues = append(issues, Issue{Field: fmt.Sprintf("preferences[%d]", index), Message: "preference is not supported"})
		}
	}
	return issues
}

func MessageIssue(message model.ChatMessage) string {
	if !model.ValidMessage(message) {
		return "message requires an itinerary, sender, sequence, and body"
	}
	if strings.TrimSpace(message.Body) != message.Body {
		return "message body cannot have surrounding whitespace"
	}
	return ""
}

func DateSuggestionIssue(suggestion model.DateSuggestion) string {
	if suggestion.ID == "" || suggestion.Itinerary == "" {
		return "suggestion identity is required"
	}
	if suggestion.Suggested == "" {
		return "suggested date is required"
	}
	if suggestion.Reason == "" {
		return "suggestion reason is required"
	}
	return ""
}

func Summarize(issues []Issue) string {
	if len(issues) == 0 {
		return ""
	}
	parts := make([]string, 0, len(issues))
	for _, issue := range issues {
		parts = append(parts, issue.Field+": "+issue.Message)
	}
	return strings.Join(parts, "; ")
}

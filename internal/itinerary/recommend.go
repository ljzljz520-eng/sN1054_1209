package itinerary

import (
	"fmt"
	"strings"

	"example.com/familyitinerary/internal/model"
)

type Recommendation struct {
	Activity    model.Activity
	Explanation string
}

func Recommend(activities []model.Activity, preferences []string) []Recommendation {
	items := make([]Recommendation, 0)
	for _, activity := range activities {
		reasons := make([]string, 0)
		if activity.ChildSafe {
			reasons = append(reasons, "child-safe")
		}
		for _, preference := range preferences {
			if preference == "indoor-backup" && activity.Category == "indoor" {
				reasons = append(reasons, "indoor backup")
			}
			if preference == "short-transfers" && activity.Duration <= 90 {
				reasons = append(reasons, "short transfer")
			}
			if preference == "stroller-access" && activity.Category == "park" {
				reasons = append(reasons, "stroller access")
			}
		}
		if len(reasons) > 0 {
			items = append(items, Recommendation{Activity: activity, Explanation: strings.Join(reasons, ", ")})
		}
	}
	return items
}

func ValidateDateWindow(start, end string, activities []model.Activity) error {
	if start == "" || end == "" {
		return fmt.Errorf("date window is required")
	}
	if start > end {
		return fmt.Errorf("start date is after end date")
	}
	if len(activities) == 0 {
		return fmt.Errorf("date window has no planned activities")
	}
	return nil
}

func SuggestDate(start, end string, reason string) (model.DateSuggestion, error) {
	if start == "" || end == "" {
		return model.DateSuggestion{}, fmt.Errorf("date range is required")
	}
	if start > end {
		return model.DateSuggestion{}, fmt.Errorf("date range is inverted")
	}
	if reason == "" {
		reason = "balanced travel pace"
	}
	return model.NewSuggestion("suggestion-"+start, "", start, reason), nil
}

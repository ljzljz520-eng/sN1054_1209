package itinerary

import (
	"sort"

	"example.com/familyitinerary/internal/model"
)

type Score struct {
	ActivityID string
	Value      int
	Reasons    []string
}

func ScoreActivities(activities []model.Activity, preferences []string) []Score {
	scores := make([]Score, 0, len(activities))
	for _, activity := range activities {
		value := 0
		reasons := make([]string, 0)
		if activity.ChildSafe {
			value += 40
			reasons = append(reasons, "child-safe")
		}
		if activity.Duration <= 90 {
			value += 20
			reasons = append(reasons, "manageable duration")
		} else {
			value -= 10
		}
		for _, preference := range preferences {
			switch preference {
			case "indoor-backup":
				if activity.Category == "indoor" {
					value += 25
					reasons = append(reasons, "indoor")
				}
			case "short-transfers":
				if activity.Duration <= 90 {
					value += 15
					reasons = append(reasons, "short")
				}
			case "stroller-access":
				if activity.Category == "park" {
					value += 15
					reasons = append(reasons, "stroller route")
				}
			case "early-dinner":
				if activity.Category == "food" {
					value += 10
					reasons = append(reasons, "early meal")
				}
			}
		}
		scores = append(scores, Score{ActivityID: activity.ID, Value: value, Reasons: reasons})
	}
	sort.SliceStable(scores, func(i, j int) bool {
		if scores[i].Value == scores[j].Value {
			return scores[i].ActivityID < scores[j].ActivityID
		}
		return scores[i].Value > scores[j].Value
	})
	return scores
}

func TopActivities(activities []model.Activity, preferences []string, limit int) []model.Activity {
	if limit <= 0 {
		return []model.Activity{}
	}
	scores := ScoreActivities(activities, preferences)
	byID := make(map[string]model.Activity, len(activities))
	for _, activity := range activities {
		byID[activity.ID] = activity
	}
	result := make([]model.Activity, 0, limit)
	for _, score := range scores {
		if activity, ok := byID[score.ActivityID]; ok {
			result = append(result, activity)
		}
		if len(result) == limit {
			break
		}
	}
	return result
}

func ScheduleScore(schedule Schedule, preferences []string) int {
	total := 0
	for _, score := range ScoreActivities(schedule.AllActivities(), preferences) {
		total += score.Value
	}
	if len(schedule.Days) > 4 {
		total -= (len(schedule.Days) - 4) * 5
	}
	return total
}

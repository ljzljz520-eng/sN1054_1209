package itinerary

import (
	"fmt"

	"example.com/familyitinerary/internal/model"
)

type Constraint struct {
	Name     string
	Limit    int
	Required bool
}

func DefaultConstraints() []Constraint {
	return []Constraint{{Name: "daily-minutes", Limit: 600, Required: true}, {Name: "minimum-child-safe", Limit: 1, Required: true}, {Name: "maximum-transitions", Limit: 3, Required: false}}
}

func CheckConstraints(schedule Schedule, children int, constraints []Constraint) []string {
	issues := make([]string, 0)
	for _, constraint := range constraints {
		switch constraint.Name {
		case "daily-minutes":
			if err := scheduleCapacity(schedule, constraint.Limit); err != nil && constraint.Required {
				issues = append(issues, err.Error())
			}
		case "minimum-child-safe":
			if childSafeCount(schedule) < constraint.Limit && constraint.Required {
				issues = append(issues, fmt.Sprintf("at least %d child-safe activities are required", constraint.Limit))
			}
		case "maximum-transitions":
			if schedule.Itinerary.Status == "accepted" && len(schedule.Days) > constraint.Limit {
				issues = append(issues, "accepted trip has too many transition days")
			}
		default:
			issues = append(issues, "unknown constraint: "+constraint.Name)
		}
	}
	if children == 0 && len(schedule.AllActivities()) == 0 {
		issues = append(issues, "empty schedule")
	}
	return issues
}

func scheduleCapacity(schedule Schedule, limit int) error {
	for day := range schedule.Days {
		minutes := 0
		for _, activity := range schedule.Days[day] {
			minutes += activity.Duration
		}
		if minutes > limit {
			return fmt.Errorf("day %d exceeds %d minutes", day, limit)
		}
	}
	return nil
}

func childSafeCount(schedule Schedule) int {
	count := 0
	for _, activity := range schedule.AllActivities() {
		if activity.ChildSafe {
			count++
		}
	}
	return count
}

func FilterChildSafe(activities []model.Activity) []model.Activity {
	filtered := make([]model.Activity, 0)
	for _, activity := range activities {
		if activity.ChildSafe {
			filtered = append(filtered, activity)
		}
	}
	return filtered
}

package itinerary

import (
	"fmt"
	"strings"

	"example.com/familyitinerary/internal/model"
)

type Schedule struct {
	Itinerary model.Itinerary
	Days      map[int][]model.Activity
}

func BuildSchedule(item model.Itinerary, activities []model.Activity) Schedule {
	days := make(map[int][]model.Activity)
	for _, activity := range activities {
		days[activity.Day] = append(days[activity.Day], activity)
	}
	return Schedule{Itinerary: item, Days: days}
}

func (s Schedule) Day(day int) []model.Activity { return append([]model.Activity(nil), s.Days[day]...) }

func (s Schedule) ChildSafeRatio() float64 {
	if len(s.AllActivities()) == 0 {
		return 0
	}
	safe := 0
	for _, activity := range s.AllActivities() {
		if activity.ChildSafe {
			safe++
		}
	}
	return float64(safe) / float64(len(s.AllActivities()))
}

func (s Schedule) AllActivities() []model.Activity {
	items := make([]model.Activity, 0)
	for _, day := range s.Days {
		items = append(items, day...)
	}
	return items
}

func (s Schedule) CheckChildFit(children int) error {
	if children < 0 {
		return fmt.Errorf("children cannot be negative")
	}
	if children == 0 {
		return nil
	}
	if s.ChildSafeRatio() < 0.5 {
		return fmt.Errorf("fewer than half of activities are child safe")
	}
	return nil
}

func (s Schedule) Describe() string {
	parts := []string{s.Itinerary.Title + " in " + s.Itinerary.Destination}
	for day := 1; day <= len(s.Days); day++ {
		activities := s.Days[day]
		if len(activities) == 0 {
			continue
		}
		labels := make([]string, 0, len(activities))
		for _, activity := range activities {
			labels = append(labels, activity.Title)
		}
		parts = append(parts, fmt.Sprintf("Day %d: %s", day, strings.Join(labels, ", ")))
	}
	return strings.Join(parts, " | ")
}

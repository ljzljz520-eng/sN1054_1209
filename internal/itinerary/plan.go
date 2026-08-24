package itinerary

import (
	"fmt"
	"sort"

	"example.com/familyitinerary/internal/model"
)

type Planner struct {
	activities map[string][]model.Activity
}

func NewPlanner() *Planner {
	return &Planner{activities: make(map[string][]model.Activity)}
}

func (p *Planner) AddActivity(activity model.Activity) error {
	if !model.ValidActivity(activity) {
		return fmt.Errorf("invalid activity %q", activity.ID)
	}
	items := p.activities[activity.Itinerary]
	for _, existing := range items {
		if existing.ID == activity.ID {
			return fmt.Errorf("activity %q already exists", activity.ID)
		}
	}
	p.activities[activity.Itinerary] = append(items, activity)
	return nil
}

func (p *Planner) ReplaceActivity(activity model.Activity) error {
	if !model.ValidActivity(activity) {
		return fmt.Errorf("invalid activity %q", activity.ID)
	}
	items := p.activities[activity.Itinerary]
	for index, existing := range items {
		if existing.ID == activity.ID {
			items[index] = activity
			p.activities[activity.Itinerary] = items
			return nil
		}
	}
	return p.AddActivity(activity)
}

func (p *Planner) Activities(itineraryID string) []model.Activity {
	items := append([]model.Activity(nil), p.activities[itineraryID]...)
	sort.SliceStable(items, func(i, j int) bool {
		if items[i].Day == items[j].Day {
			return items[i].ID < items[j].ID
		}
		return items[i].Day < items[j].Day
	})
	return items
}

func (p *Planner) DayCount(itineraryID string) int {
	maxDay := 0
	for _, activity := range p.activities[itineraryID] {
		if activity.Day > maxDay {
			maxDay = activity.Day
		}
	}
	return maxDay
}

func (p *Planner) ValidateCapacity(itineraryID string, maxMinutes int) error {
	perDay := make(map[int]int)
	for _, activity := range p.activities[itineraryID] {
		perDay[activity.Day] += activity.Duration
	}
	for day, minutes := range perDay {
		if minutes > maxMinutes {
			return fmt.Errorf("day %d has %d minutes, limit is %d", day, minutes, maxMinutes)
		}
	}
	return nil
}

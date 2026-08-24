package itinerary

import (
	"fmt"
	"sort"

	"example.com/familyitinerary/internal/model"
)

type Change struct {
	Day     int
	Added   []model.Activity
	Removed []string
	Notes   []string
}

func DiffActivities(before, after []model.Activity) Change {
	oldByID := make(map[string]model.Activity, len(before))
	newByID := make(map[string]model.Activity, len(after))
	for _, activity := range before {
		oldByID[activity.ID] = activity
	}
	for _, activity := range after {
		newByID[activity.ID] = activity
	}
	change := Change{Added: make([]model.Activity, 0), Removed: make([]string, 0), Notes: make([]string, 0)}
	for id, activity := range newByID {
		if _, exists := oldByID[id]; !exists {
			change.Added = append(change.Added, activity)
			change.Notes = append(change.Notes, fmt.Sprintf("added %s", activity.Title))
		}
	}
	for id, activity := range oldByID {
		if _, exists := newByID[id]; !exists {
			change.Removed = append(change.Removed, id)
			change.Notes = append(change.Notes, fmt.Sprintf("removed %s", activity.Title))
		}
	}
	if len(change.Added) > 0 {
		change.Day = change.Added[0].Day
	}
	sort.Strings(change.Removed)
	sort.Strings(change.Notes)
	return change
}

func ApplyChange(activities []model.Activity, change Change) []model.Activity {
	removed := make(map[string]bool, len(change.Removed))
	for _, id := range change.Removed {
		removed[id] = true
	}
	merged := make([]model.Activity, 0, len(activities)+len(change.Added))
	for _, activity := range activities {
		if !removed[activity.ID] {
			merged = append(merged, activity)
		}
	}
	merged = append(merged, change.Added...)
	sort.SliceStable(merged, func(i, j int) bool {
		if merged[i].Day == merged[j].Day {
			return merged[i].ID < merged[j].ID
		}
		return merged[i].Day < merged[j].Day
	})
	return merged
}

func HasConflict(activities []model.Activity) bool {
	byDay := make(map[int]int)
	for _, activity := range activities {
		byDay[activity.Day] += activity.Duration
		if byDay[activity.Day] > 600 {
			return true
		}
	}
	return false
}

func ResolveConflict(activities []model.Activity) ([]model.Activity, []model.Activity) {
	safe := make([]model.Activity, 0, len(activities))
	deferred := make([]model.Activity, 0)
	byDay := make(map[int]int)
	for _, activity := range activities {
		if byDay[activity.Day]+activity.Duration <= 600 {
			safe = append(safe, activity)
			byDay[activity.Day] += activity.Duration
		} else {
			deferred = append(deferred, activity)
		}
	}
	return safe, deferred
}

func ActivityTitles(activities []model.Activity) map[int][]string {
	titles := make(map[int][]string)
	for _, activity := range activities {
		titles[activity.Day] = append(titles[activity.Day], activity.Title)
	}
	for day := range titles {
		sort.Strings(titles[day])
	}
	return titles
}

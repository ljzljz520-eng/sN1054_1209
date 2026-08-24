package report

import (
	"sort"
	"strings"

	"example.com/familyitinerary/internal/advisor"
	"example.com/familyitinerary/internal/model"
)

type ActivityLine struct {
	Day       int
	Title     string
	Category  string
	ChildSafe bool
}

func ActivityLines(r advisor.Report) []ActivityLine {
	items := make([]ActivityLine, 0)
	for _, activity := range r.Schedule.AllActivities() {
		items = append(items, ActivityLine{Day: activity.Day, Title: activity.Title, Category: activity.Category, ChildSafe: activity.ChildSafe})
	}
	sort.SliceStable(items, func(i, j int) bool {
		if items[i].Day == items[j].Day {
			return items[i].Title < items[j].Title
		}
		return items[i].Day < items[j].Day
	})
	return items
}

func Markdown(r advisor.Report) string {
	lines := []string{"# " + r.Intake.Itinerary.Title, "", "Destination: " + r.Intake.Itinerary.Destination, "Status: " + r.Intake.Itinerary.Status, ""}
	for _, line := range ActivityLines(r) {
		marker := ""
		if line.ChildSafe {
			marker = " [child-safe]"
		}
		lines = append(lines, "- Day "+itoa(line.Day)+": "+line.Title+" ("+line.Category+")"+marker)
	}
	if len(r.Recommendations) > 0 {
		lines = append(lines, "", "Recommendations:")
		for _, recommendation := range r.Recommendations {
			lines = append(lines, "- "+recommendation.Activity.Title+": "+recommendation.Explanation)
		}
	}
	return strings.Join(lines, "\n")
}

func itoa(value int) string {
	if value == 0 {
		return "0"
	}
	digits := make([]byte, 0, 4)
	for value > 0 {
		digits = append([]byte{byte('0' + value%10)}, digits...)
		value /= 10
	}
	return string(digits)
}

func StatusCounts(messages []model.ChatMessage) map[string]int {
	counts := make(map[string]int)
	for _, message := range messages {
		counts[message.Status]++
	}
	return counts
}

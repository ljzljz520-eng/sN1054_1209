package itinerary

import (
	"fmt"
	"strconv"
	"strings"

	"example.com/familyitinerary/internal/model"
)

type DateParts struct{ Year, Month, Day int }

func ParseDate(value string) (DateParts, error) {
	parts := strings.Split(value, "-")
	if len(parts) != 3 {
		return DateParts{}, fmt.Errorf("date %q must be YYYY-MM-DD", value)
	}
	year, err := strconv.Atoi(parts[0])
	if err != nil {
		return DateParts{}, fmt.Errorf("invalid year")
	}
	month, err := strconv.Atoi(parts[1])
	if err != nil {
		return DateParts{}, fmt.Errorf("invalid month")
	}
	day, err := strconv.Atoi(parts[2])
	if err != nil {
		return DateParts{}, fmt.Errorf("invalid day")
	}
	if year < 2000 || month < 1 || month > 12 || day < 1 || day > 31 {
		return DateParts{}, fmt.Errorf("date %q is out of range", value)
	}
	return DateParts{Year: year, Month: month, Day: day}, nil
}

func DateDistance(start, end string) (int, error) {
	left, err := ParseDate(start)
	if err != nil {
		return 0, err
	}
	right, err := ParseDate(end)
	if err != nil {
		return 0, err
	}
	leftIndex := left.Year*372 + left.Month*31 + left.Day
	rightIndex := right.Year*372 + right.Month*31 + right.Day
	if rightIndex < leftIndex {
		return 0, fmt.Errorf("end date precedes start date")
	}
	return rightIndex - leftIndex, nil
}

func DayLabels(item model.Itinerary) ([]string, error) {
	distance, err := DateDistance(item.StartDate, item.EndDate)
	if err != nil {
		return nil, err
	}
	labels := make([]string, 0, distance+1)
	start, _ := ParseDate(item.StartDate)
	for offset := 0; offset <= distance; offset++ {
		day := start.Day + offset
		month := start.Month
		year := start.Year
		for day > 31 {
			day -= 31
			month++
		}
		for month > 12 {
			month -= 12
			year++
		}
		labels = append(labels, fmt.Sprintf("%04d-%02d-%02d", year, month, day))
	}
	return labels, nil
}

func ActivityByDay(schedule Schedule, day int) []model.Activity {
	activities := schedule.Day(day)
	if len(activities) < 2 {
		return activities
	}
	for index := 0; index < len(activities)-1; index++ {
		for next := index + 1; next < len(activities); next++ {
			if activities[next].Duration < activities[index].Duration {
				activities[index], activities[next] = activities[next], activities[index]
			}
		}
	}
	return activities
}

func IsSchoolDay(date string, blocked []string) bool {
	for _, candidate := range blocked {
		if candidate == date {
			return true
		}
	}
	return false
}

func AvailableDate(date string, blocked []string) bool {
	if _, err := ParseDate(date); err != nil {
		return false
	}
	return !IsSchoolDay(date, blocked)
}

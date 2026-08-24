package advisor

import (
	"fmt"
	"strings"

	"example.com/familyitinerary/internal/itinerary"
	"example.com/familyitinerary/internal/model"
	"example.com/familyitinerary/internal/store"
)

type Report struct {
	Intake          model.Intake
	Schedule        itinerary.Schedule
	Recommendations []itinerary.Recommendation
	Summary         string
}

func BuildReport(repository *store.Store, itineraryID string) (Report, error) {
	service := NewIntakeService(repository)
	in, err := service.Load(itineraryID)
	if err != nil {
		return Report{}, err
	}
	schedule := itinerary.BuildSchedule(in.Itinerary, in.Activities)
	if err := schedule.CheckChildFit(in.Family.Children); err != nil {
		return Report{}, err
	}
	recommendations := itinerary.Recommend(in.Activities, in.Preferences)
	return Report{Intake: in, Schedule: schedule, Recommendations: recommendations, Summary: BuildSummary(in, schedule, recommendations)}, nil
}

func BuildSummary(in model.Intake, schedule itinerary.Schedule, recommendations []itinerary.Recommendation) string {
	parts := []string{fmt.Sprintf("Family %s has %d children", in.Family.Name, in.Family.Children), schedule.Describe()}
	if len(recommendations) > 0 {
		parts = append(parts, fmt.Sprintf("%d activities match preferences", len(recommendations)))
	}
	if schedule.Itinerary.Status != "draft" {
		parts = append(parts, "Status: "+schedule.Itinerary.Status)
	}
	return strings.Join(parts, "\n")
}

func (r Report) ChildSafeActivities() int {
	count := 0
	for _, activity := range r.Schedule.AllActivities() {
		if activity.ChildSafe {
			count++
		}
	}
	return count
}

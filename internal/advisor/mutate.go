package advisor

import (
	"fmt"

	"example.com/familyitinerary/internal/itinerary"
	"example.com/familyitinerary/internal/model"
	"example.com/familyitinerary/internal/store"
)

type ChangeService struct {
	repository *store.Store
	planner    *itinerary.Planner
}

func NewChangeService(repository *store.Store) *ChangeService {
	return &ChangeService{repository: repository, planner: itinerary.NewPlanner()}
}

func (s *ChangeService) ReplaceActivity(activity model.Activity) error {
	if !model.ValidActivity(activity) {
		return fmt.Errorf("invalid activity")
	}
	if err := s.planner.ReplaceActivity(activity); err != nil {
		return err
	}
	if err := s.repository.SaveActivity(activity); err != nil {
		return err
	}
	return s.repository.RecordAudit(store.AuditEntry{Entity: "activity", EntityID: activity.ID, Action: "replace", Detail: activity.Title})
}

func (s *ChangeService) RemoveActivity(itineraryID, activityID string) error {
	if itineraryID == "" || activityID == "" {
		return fmt.Errorf("itinerary and activity are required")
	}
	if err := s.repository.DeleteActivity(activityID); err != nil {
		return err
	}
	return s.repository.RecordAudit(store.AuditEntry{Entity: "activity", EntityID: activityID, Action: "remove", Detail: itineraryID})
}

func (s *ChangeService) AcceptSuggestion(suggestionID string) error {
	if suggestionID == "" {
		return fmt.Errorf("suggestion is required")
	}
	if err := s.repository.AcceptSuggestion(suggestionID); err != nil {
		return err
	}
	return s.repository.RecordAudit(store.AuditEntry{Entity: "date_suggestion", EntityID: suggestionID, Action: "accept", Detail: "family approved date"})
}

func (s *ChangeService) ValidateChange(itineraryID string, activities []model.Activity) []string {
	schedule, err := s.repository.GetItinerary(itineraryID)
	if err != nil {
		return []string{err.Error()}
	}
	issues := CheckChangedSchedule(schedule, activities)
	return issues
}

func CheckChangedSchedule(item model.Itinerary, activities []model.Activity) []string {
	issues := make([]string, 0)
	schedule := itinerary.BuildSchedule(item, activities)
	issues = append(issues, itinerary.CheckConstraints(schedule, 1, itinerary.DefaultConstraints())...)
	if len(activities) == 0 {
		issues = append(issues, "changed schedule must retain an activity")
	}
	return issues
}

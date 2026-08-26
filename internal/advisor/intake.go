package advisor

import (
	"fmt"

	"example.com/familyitinerary/internal/itinerary"
	"example.com/familyitinerary/internal/model"
	"example.com/familyitinerary/internal/store"
	"example.com/familyitinerary/internal/validation"
)

type IntakeService struct {
	store   *store.Store
	planner *itinerary.Planner
}

func NewIntakeService(repository *store.Store) *IntakeService {
	return &IntakeService{store: repository, planner: itinerary.NewPlanner()}
}

func (s *IntakeService) Store() *store.Store { return s.store }

func (s *IntakeService) Register(in model.Intake) error {
	if issue := validation.Summarize(validation.IntakeIssues(in)); issue != "" {
		return fmt.Errorf("invalid intake: %s", issue)
	}
	for _, activity := range in.Activities {
		if err := s.planner.AddActivity(activity); err != nil {
			return err
		}
	}
	if err := s.planner.ValidateCapacity(in.Itinerary.ID, 600); err != nil {
		return err
	}
	return s.store.SaveIntake(in)
}

func (s *IntakeService) Load(itineraryID string) (model.Intake, error) {
	it, err := s.store.GetItinerary(itineraryID)
	if err != nil {
		return model.Intake{}, err
	}
	family, err := s.store.GetFamily(it.FamilyID)
	if err != nil {
		return model.Intake{}, err
	}
	activities, err := s.store.ListActivities(itineraryID)
	if err != nil {
		return model.Intake{}, err
	}
	preferences, err := s.store.LoadPreferences(itineraryID)
	if err != nil {
		return model.Intake{}, err
	}
	return model.Intake{Family: family, Itinerary: it, Activities: activities, Preferences: preferences}, nil
}

func (s *IntakeService) ChangeStatus(itineraryID, status string) error {
	if status == "" {
		return fmt.Errorf("status is required")
	}
	return s.store.UpdateItineraryStatus(itineraryID, status)
}

func (s *IntakeService) AddActivity(activity model.Activity) error {
	if err := s.planner.ReplaceActivity(activity); err != nil {
		return err
	}
	return s.store.SaveActivity(activity)
}

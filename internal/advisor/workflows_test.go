package advisor

import (
	"testing"

	"example.com/familyitinerary/internal/model"
	"example.com/familyitinerary/internal/store"
)

func workflowIntake() model.Intake {
	family := model.NewFamily("fam-flow", "Zhao", 2, "nap")
	itinerary := model.NewItinerary("it-flow", family.ID, "Island days", "Xiamen", "2026-08-01", "2026-08-04")
	activities := []model.Activity{model.NewActivity("act-flow-1", itinerary.ID, 1, "Beach", "park", 120, true), model.NewActivity("act-flow-2", itinerary.ID, 2, "Science hall", "indoor", 90, true)}
	return model.Intake{Family: family, Itinerary: itinerary, Activities: activities, Preferences: []string{"indoor-backup", "stroller-access"}}
}

func TestWorkflowOne(t *testing.T) {
	repository, err := store.OpenMemory()
	if err != nil {
		t.Fatal(err)
	}
	defer repository.Close()
	service := NewIntakeService(repository)
	if err = service.Register(workflowIntake()); err != nil {
		t.Fatal(err)
	}
	loaded, err := service.Load("it-flow")
	if err != nil {
		t.Fatal(err)
	}
	if loaded.Family.Name != "Zhao" || len(loaded.Activities) != 2 {
		t.Fatal("intake workflow failed")
	}
}

func TestWorkflowTwo(t *testing.T) {
	repository, err := store.OpenMemory()
	if err != nil {
		t.Fatal(err)
	}
	defer repository.Close()
	service := NewIntakeService(repository)
	in := workflowIntake()
	if err = service.Register(in); err != nil {
		t.Fatal(err)
	}
	if err = service.ChangeStatus(in.Itinerary.ID, "accepted"); err != nil {
		t.Fatal(err)
	}
	report, err := BuildReport(repository, in.Itinerary.ID)
	if err != nil {
		t.Fatal(err)
	}
	if report.Intake.Itinerary.Status != "accepted" || report.ChildSafeActivities() != 2 {
		t.Fatal("report workflow failed")
	}
}

func TestWorkflowThree(t *testing.T) {
	repository, err := store.OpenMemory()
	if err != nil {
		t.Fatal(err)
	}
	defer repository.Close()
	intake := NewIntakeService(repository)
	conversation := NewConversationService(repository)
	in := workflowIntake()
	if err = intake.Register(in); err != nil {
		t.Fatal(err)
	}
	if _, err = conversation.Upload(in.Itinerary.ID, in.Itinerary.Title, true); err != nil {
		t.Fatal(err)
	}
	if _, err = conversation.SuggestDate(in.Itinerary.ID, "2026-08-02", true); err != nil {
		t.Fatal(err)
	}
	history, err := conversation.History(in.Itinerary.ID)
	if err != nil || len(history) != 2 {
		t.Fatalf("history=%d err=%v", len(history), err)
	}
}

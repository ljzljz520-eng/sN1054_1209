package httpapi

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"example.com/familyitinerary/internal/advisor"
	"example.com/familyitinerary/internal/model"
	"example.com/familyitinerary/internal/store"
)

func TestHealthAndReportEndpoints(t *testing.T) {
	repository, err := store.OpenMemory()
	if err != nil {
		t.Fatal(err)
	}
	defer repository.Close()
	family := model.NewFamily("f-http", "He", 1, "")
	itinerary := model.NewItinerary("i-http", family.ID, "City", "Nanjing", "2026-07-01", "2026-07-02")
	intake := advisor.NewIntakeService(repository)
	if err = intake.Register(model.Intake{Family: family, Itinerary: itinerary, Activities: []model.Activity{model.NewActivity("a-http", itinerary.ID, 1, "Museum", "indoor", 80, true)}}); err != nil {
		t.Fatal(err)
	}
	handler := New(intake, advisor.NewConversationService(repository)).Handler()
	health := httptest.NewRecorder()
	handler.ServeHTTP(health, httptest.NewRequest(http.MethodGet, "/health", nil))
	if health.Code != http.StatusOK {
		t.Fatalf("health=%d", health.Code)
	}
	report := httptest.NewRecorder()
	handler.ServeHTTP(report, httptest.NewRequest(http.MethodGet, "/itineraries/i-http", nil))
	if report.Code != http.StatusOK || report.Body.Len() == 0 {
		t.Fatalf("report code=%d body=%s", report.Code, report.Body.String())
	}
}

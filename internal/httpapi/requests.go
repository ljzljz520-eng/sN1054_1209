package httpapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"example.com/familyitinerary/internal/model"
)

type intakeRequest struct {
	Family      model.Family     `json:"family"`
	Itinerary   model.Itinerary  `json:"itinerary"`
	Activities  []model.Activity `json:"activities"`
	Preferences []string         `json:"preferences"`
}

type statusRequest struct {
	Status string `json:"status"`
}

func decodeJSON(request *http.Request, target interface{}) error {
	if request.Body == nil {
		return fmt.Errorf("request body is required")
	}
	decoder := json.NewDecoder(request.Body)
	if err := decoder.Decode(target); err != nil {
		return fmt.Errorf("invalid json: %w", err)
	}
	return nil
}

func writeJSON(writer http.ResponseWriter, status int, value interface{}) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(status)
	json.NewEncoder(writer).Encode(value)
}

func parsePath(path string) ([]string, error) {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) < 2 || parts[0] != "itineraries" || parts[1] == "" {
		return nil, fmt.Errorf("invalid itinerary path")
	}
	return parts, nil
}

func activityPayload(in model.Intake) map[string]interface{} {
	return map[string]interface{}{"family": in.Family, "itinerary": in.Itinerary, "activities": in.Activities, "preferences": in.Preferences}
}

func methodAllowed(request *http.Request, methods ...string) bool {
	for _, method := range methods {
		if request.Method == method {
			return true
		}
	}
	return false
}

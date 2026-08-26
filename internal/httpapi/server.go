package httpapi

import (
	"encoding/json"
	"net/http"
	"strings"

	"example.com/familyitinerary/internal/advisor"
	"example.com/familyitinerary/internal/report"
)

type Server struct {
	intake       *advisor.IntakeService
	conversation *advisor.ConversationService
}

func New(intake *advisor.IntakeService, conversation *advisor.ConversationService) *Server {
	return &Server{intake: intake, conversation: conversation}
}

func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", s.health)
	mux.HandleFunc("/itineraries/", s.itinerary)
	return mux
}

func (s *Server) health(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(map[string]string{"status": "ok"})
}

func (s *Server) itinerary(writer http.ResponseWriter, request *http.Request) {
	parts := strings.Split(strings.Trim(request.URL.Path, "/"), "/")
	if len(parts) < 2 || parts[1] == "" {
		http.Error(writer, "itinerary id is required", http.StatusBadRequest)
		return
	}
	id := parts[1]
	if len(parts) == 2 && request.Method == http.MethodGet {
		s.getReport(writer, id)
		return
	}
	if len(parts) == 3 && parts[2] == "messages" && request.Method == http.MethodGet {
		s.getHistory(writer, id)
		return
	}
	http.NotFound(writer, request)
}

func (s *Server) getReport(writer http.ResponseWriter, id string) {
	r, err := advisor.BuildReport(s.intake.Store(), id)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusNotFound)
		return
	}
	payload, err := report.JSON(r)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.Write(payload)
}

func (s *Server) getHistory(writer http.ResponseWriter, id string) {
	items, err := s.conversation.History(id)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusNotFound)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(items)
}

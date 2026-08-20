package httpapi

import (
	"emergency-claim-code/internal/model"
	"emergency-claim-code/internal/service"
	"emergency-claim-code/internal/workflow"
	"encoding/json"
	"net/http"
	"strings"
)

type Server struct {
	engine  *workflow.Engine
	service *service.Service
}

func New(engine *workflow.Engine, svc *service.Service) *Server {
	return &Server{engine: engine, service: svc}
}

func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", s.health)
	mux.HandleFunc("/records", s.records)
	mux.HandleFunc("/records/", s.record)
	return mux
}

func (s *Server) health(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": s.service.Health()})
}

func (s *Server) records(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		term := r.URL.Query().Get("q")
		status := r.URL.Query().Get("status")
		items, err := s.service.ListRecords(term, status)
		if err != nil {
			writeError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, items)
	case http.MethodPost:
		var input model.Record
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			writeError(w, err)
			return
		}
		item, err := s.engine.Register(input, r.Header.Get("X-Actor"), r.Header.Get("X-At"))
		if err != nil {
			writeError(w, err)
			return
		}
		writeJSON(w, http.StatusCreated, item)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (s *Server) record(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/records/")
	if id == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	switch r.Method {
	case http.MethodGet:
		item, err := s.service.GetRecord(id)
		if err != nil {
			writeError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, item)
	case http.MethodPost:
		action := r.URL.Query().Get("action")
		item, err := s.service.ChangeStatus(id, action, r.Header.Get("X-Actor"), r.Header.Get("X-At"))
		if err != nil {
			writeError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, item)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func writeError(w http.ResponseWriter, err error) {
	writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
}

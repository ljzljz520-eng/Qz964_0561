package api

import (
	"encoding/json"
	"net/http"
	"strings"
	"timber-safety/internal/service"
	"timber-safety/internal/validation"
)

type Server struct {
	Service *service.Service
	mux     *http.ServeMux
}

func New(s *service.Service) *Server {
	a := &Server{Service: s, mux: http.NewServeMux()}
	a.routes()
	return a
}
func (a *Server) routes() {
	a.mux.HandleFunc("/health", a.health)
	a.mux.HandleFunc("/records", a.records)
	a.mux.HandleFunc("/records/", a.record)
	a.mux.HandleFunc("/reports/summary", a.summary)
}
func (a *Server) Handler() http.Handler { return a.mux }
func (a *Server) health(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
func (a *Server) records(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		q := validation.Query{State: r.URL.Query().Get("state"), Risk: r.URL.Query().Get("risk"), Source: r.URL.Query().Get("source")}
		items, err := a.Service.Search(q)
		if err != nil {
			writeError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, items)
		return
	}
	writeErrorStatus(w, http.StatusMethodNotAllowed, "method not allowed")
}
func (a *Server) record(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) < 2 {
		writeErrorStatus(w, http.StatusNotFound, "record not found")
		return
	}
	id := parts[1]
	if len(parts) == 2 && r.Method == http.MethodGet {
		item, err := a.Service.Get(id)
		if err != nil {
			writeError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, item)
		return
	}
	writeErrorStatus(w, http.StatusNotImplemented, "unsupported record action")
}
func (a *Server) summary(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeErrorStatus(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	sum, err := a.Service.Summary(validation.Query{})
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, sum)
}
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
func writeError(w http.ResponseWriter, err error) {
	writeErrorStatus(w, http.StatusBadRequest, err.Error())
}
func writeErrorStatus(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

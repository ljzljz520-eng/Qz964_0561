package api

import (
	"net/http"
	"strconv"
	"timber-safety/internal/validation"
)

func parseQuery(r *http.Request) (validation.Query, error) {
	q := validation.Query{State: r.URL.Query().Get("state"), Risk: r.URL.Query().Get("risk"), Source: r.URL.Query().Get("source")}
	if raw := r.URL.Query().Get("limit"); raw != "" {
		n, e := strconv.Atoi(raw)
		if e != nil {
			return q, e
		}
		q.Limit = n
	}
	return q, nil
}
func statusForError(err error) int {
	if err == nil {
		return http.StatusOK
	}
	return http.StatusBadRequest
}
func methodAllowed(r *http.Request, methods ...string) bool {
	for _, m := range methods {
		if r.Method == m {
			return true
		}
	}
	return false
}

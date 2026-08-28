package api

import (
	"encoding/json"
	"net/http"
	"timber-safety/internal/domain"
)

func decodeRecord(r *http.Request) (domain.Record, error) {
	var v domain.Record
	err := json.NewDecoder(r.Body).Decode(&v)
	return v, err
}
func requestID(r *http.Request) string           { return r.Header.Get("X-Request-ID") }
func writeAccepted(w http.ResponseWriter, v any) { writeJSON(w, http.StatusAccepted, v) }
func pathTail(path string) string {
	for i := len(path) - 1; i >= 0; i-- {
		if path[i] == '/' {
			return path[i+1:]
		}
	}
	return path
}

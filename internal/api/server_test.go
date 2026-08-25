package api

import (
	"net/http/httptest"
	"testing"
	"timber-safety/internal/service"
	"timber-safety/internal/store"
)

func TestHealth(t *testing.T) {
	s, _ := store.Open(t.TempDir() + "/db")
	defer s.Close()
	rr := httptest.NewRecorder()
	New(service.New(s)).Handler().ServeHTTP(rr, httptest.NewRequest("GET", "/health", nil))
	if rr.Code != 200 {
		t.Fatal(rr.Code)
	}
}

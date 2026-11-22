package controlplane

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
)


func TestCreateAndGetProject(t *testing.T) {
	store := NewProjectStore()


	var ready atomic.Bool
	ready.Store(true)

	handler := NewHandler(slog.Default(), &ready, store)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /v1/projects", handler.CreateProject)
	mux.HandleFunc("GET /v1/projects/{id}", handler.GetProject)

	body := `{"name": "Google", "timezone": "Europe/Moscow"}`

	req := httptest.NewRequest(
		http.MethodPost,
		"/v1/projects",
		strings.NewReader(body),
	)

	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", rec.Code)
	}

	var created Project

	if err := json.NewDecoder(rec.Body).Decode(&created); err != nil {
		t.Fatal(err)
	}

	if created.ID == "" {
		t.Fatal("project ID is empty")
	}

	if created.Name != "Google" {
		t.Fatalf("expected Google, got %s", created.Name)
	}

	if created.Timezone != "Europe/Moscow" {
		t.Fatalf("unexpected timezone: %s", created.Timezone)
	}

	getReq := httptest.NewRequest(
		http.MethodGet,
		"/v1/projects/"+created.ID,
		nil,
	)
	getRec := httptest.NewRecorder()
	mux.ServeHTTP(getRec, getReq)

	if getRec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", getRec.Code)
	}

	var fetched Project

	if err := json.NewDecoder(getRec.Body).Decode(&fetched); err != nil {
		t.Fatal(err)
	}

	if fetched.ID != created.ID || fetched.Name != created.Name {
		t.Fatal("fetched project does not match created project")
	}
}

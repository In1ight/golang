package controlplane

import (
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
)

func TestReadyHandler(t *testing.T) {
	var ready atomic.Bool
	ready.Store(true)

	handler := ReadyHandler(&ready)

	req := httptest.NewRequest(http.MethodGet, "/readyz", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}

	ready.Store(false)
	rec2 := httptest.NewRecorder()
	handler.ServeHTTP(rec2, req)

	if rec2.Code != http.StatusServiceUnavailable {
		t.Errorf("Expected status 503, got %d", rec2.Code)
	}
}

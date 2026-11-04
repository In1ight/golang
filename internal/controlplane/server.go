package controlplane

import (
	"fmt"
	"log/slog"
	"net/http"
	"sync/atomic"
)

func HealthHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	w.WriteHeader(http.StatusOK)
	_, err := fmt.Fprintln(w, "Health OK")
	if err != nil {
		slog.Error("Failed to write health response", "error", err)
	}
}

func ReadyHandler(ready *atomic.Bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !ready.Load() {
			w.WriteHeader(http.StatusServiceUnavailable)
			_, err := fmt.Fprintln(w, "Not Ready")
			if err != nil {
				slog.Error("Failed to write response", "error", err)
			}
			return
		}

		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprintln(w, "Ready OK")
		if err != nil {
			slog.Error("Failed to write response", "error", err)
		}
	}
}

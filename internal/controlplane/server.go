package controlplane

import (
	"fmt"
	"log/slog"
	"net/http"
	"sync/atomic"
)

type Handler struct {
	logger *slog.Logger
	ready  *atomic.Bool
}

func NewHandler(logger *slog.Logger, ready *atomic.Bool) *Handler {
	return &Handler{
		logger: logger,
		ready:  ready,
	}
}

func (h *Handler) HealthHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	w.WriteHeader(http.StatusOK)
	_, err := fmt.Fprintln(w, "Health OK")
	if err != nil {
		h.logger.Error("Failed to write health response", "error", err)
	}
}

func (h *Handler) ReadyHandler(w http.ResponseWriter, r *http.Request) {
	if !h.ready.Load() {
		w.WriteHeader(http.StatusServiceUnavailable)
		_, err := fmt.Fprintln(w, "Not Ready")
		if err != nil {
			h.logger.Error("Failed to write response", "error", err)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err := fmt.Fprintln(w, "Ready OK")
	if err != nil {
		h.logger.Error("Failed to write response", "error", err)
	}
}

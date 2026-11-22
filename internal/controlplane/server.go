package controlplane

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"sync/atomic"
	"time"
	"unicode/utf8"
)

type Handler struct {
	logger *slog.Logger
	ready  *atomic.Bool
	store  *ProjectStore
}

func NewHandler(logger *slog.Logger, ready *atomic.Bool, store *ProjectStore) *Handler {
	return &Handler{
		logger: logger,
		ready:  ready,
		store:  store,
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

func (h *Handler) CreateProject(w http.ResponseWriter, r *http.Request) {
	var input CreateProjectInput

	err := json.NewDecoder(r.Body).Decode(&input)

	if err != nil {
		h.writeJSONError(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	name := strings.TrimSpace(input.Name)

	if name == "" || utf8.RuneCountInString(input.Name) > 100 {
		h.writeJSONError(w, http.StatusUnprocessableEntity, "invalid name")
		return
	}

	if input.Timezone == "" {
		input.Timezone = "UTC"
	}

	_, err = time.LoadLocation(input.Timezone)

	if err != nil {
		h.writeJSONError(w, http.StatusUnprocessableEntity, "invalide timezone")
		return
	}

	project := Project{
		Name:     input.Name,
		Timezone: input.Timezone,
	}
	created := h.store.Create(project)

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Location", "/v1/projects/"+created.ID)
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(created); err != nil {
		h.logger.Error("Failed to encode project", "error", err)
	}
}

func (h *Handler) DeleteProject(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	deleted := h.store.Delete(id)
	if !deleted {
		h.writeJSONError(w, http.StatusNotFound, "Project not found")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) PatchProject(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	var input UpdateProjectInput

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.writeJSONError(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	if input.Name != nil {
		name := strings.TrimSpace(*input.Name)

		if name == "" || utf8.RuneCountInString(*input.Name) > 100 {
			h.writeJSONError(w, http.StatusUnprocessableEntity, "invalid name")
			return
		}
	}

	if input.Timezone != nil {
		if *input.Timezone == "" {
			h.writeJSONError(w, http.StatusUnprocessableEntity, "invalid timezone")
			return
		}

		if _, err := time.LoadLocation(*input.Timezone); err != nil {
			h.writeJSONError(w, http.StatusUnprocessableEntity, "invalid timezone")
			return
		}
	}

	updated, ok := h.store.Update(id, input)

	if !ok {
		h.writeJSONError(w, http.StatusNotFound, "project not found")
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(updated); err != nil {
		h.logger.Error("Failed to encode updated project", "error", err)
	}

}

func (h *Handler) GetProject(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	project, ok := h.store.Get(id)

	if !ok {
		h.writeJSONError(w, http.StatusNotFound, "project not found")
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(project); err != nil {
		h.logger.Error("Failed to encode project", "error", err)
	}
}

func (h *Handler) writeJSONError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(map[string]string{
		"error": message,
	}); err != nil {
		h.logger.Error("Failed to encode error", "error", err)
	}
}

func (h *Handler) ListProjects(w http.ResponseWriter, r *http.Request) {
	items := h.store.List()

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(map[string]any{
		"items": items,
	}); err != nil {
		h.logger.Error("Failed to encode projects", "error", err)
	}
}

// Package api implements the generated.ServerInterface with thin HTTP
// handlers. Business logic (checking services, normalizing status) lives in
// internal/monitoring — handlers only translate between HTTP and that
// package, per service/AGENTS.md.
package api

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"pulse/generated"
	"pulse/internal/monitoring"
)

// Handler implements generated.ServerInterface.
type Handler struct {
	checker *monitoring.Checker
	logger  *slog.Logger
}

// NewHandler builds a Handler backed by the given Checker.
func NewHandler(checker *monitoring.Checker, logger *slog.Logger) *Handler {
	if logger == nil {
		logger = slog.Default()
	}
	return &Handler{checker: checker, logger: logger}
}

var _ generated.ServerInterface = (*Handler)(nil)

// ListServices implements GET /services.
func (h *Handler) ListServices(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, h.checker.List())
}

// GetService implements GET /services/{id}.
func (h *Handler) GetService(w http.ResponseWriter, r *http.Request, id generated.ServiceId) {
	svc, err := h.checker.Get(id)
	if err != nil {
		h.writeNotFound(w, id, err)
		return
	}
	writeJSON(w, http.StatusOK, svc)
}

// CheckService implements POST /services/{id}/check.
func (h *Handler) CheckService(w http.ResponseWriter, r *http.Request, id generated.ServiceId) {
	svc, err := h.checker.Check(r.Context(), id)
	if err != nil {
		h.writeNotFound(w, id, err)
		return
	}
	writeJSON(w, http.StatusOK, svc)
}

func (h *Handler) writeNotFound(w http.ResponseWriter, id string, err error) {
	if !errors.Is(err, monitoring.ErrNotFound) {
		h.logger.Error("unexpected checker error", "service_id", id, "error", err)
		writeJSON(w, http.StatusInternalServerError, generated.ErrorResponse{
			Code:    "internal_error",
			Message: "an unexpected error occurred",
		})
		return
	}
	writeJSON(w, http.StatusNotFound, generated.ErrorResponse{
		Code:    "service_not_found",
		Message: "no monitored service exists with id '" + id + "'",
	})
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

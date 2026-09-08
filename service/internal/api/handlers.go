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
	"pulse/internal/plugins"
)

// Handler implements generated.ServerInterface.
type Handler struct {
	checker  *monitoring.Checker
	registry *plugins.Registry
	logger   *slog.Logger
}

// NewHandler builds a Handler backed by the given Checker and plugin
// Registry (the latter powers GET /plugin-types — see
// ../../adr/0002-factory-based-plugin-model-for-checks.md and
// ../../../api/adr/0002-plugin-type-discovery-endpoint.md).
func NewHandler(checker *monitoring.Checker, registry *plugins.Registry, logger *slog.Logger) *Handler {
	if logger == nil {
		logger = slog.Default()
	}
	return &Handler{checker: checker, registry: registry, logger: logger}
}

var _ generated.ServerInterface = (*Handler)(nil)

// ListPluginTypes implements GET /plugin-types.
func (h *Handler) ListPluginTypes(w http.ResponseWriter, r *http.Request) {
	descriptors := h.registry.List()
	out := make([]generated.PluginType, 0, len(descriptors))
	for _, d := range descriptors {
		var inputSchema, outputSchema map[string]interface{}
		if err := json.Unmarshal(d.InputSchema, &inputSchema); err != nil {
			h.logger.Error("invalid plugin input schema", "type", d.Type, "error", err)
			writeJSON(w, http.StatusInternalServerError, generated.ErrorResponse{
				Code:    "internal_error",
				Message: "an unexpected error occurred",
			})
			return
		}
		if err := json.Unmarshal(d.OutputSchema, &outputSchema); err != nil {
			h.logger.Error("invalid plugin output schema", "type", d.Type, "error", err)
			writeJSON(w, http.StatusInternalServerError, generated.ErrorResponse{
				Code:    "internal_error",
				Message: "an unexpected error occurred",
			})
			return
		}
		out = append(out, generated.PluginType{
			Type:         d.Type,
			InputSchema:  inputSchema,
			OutputSchema: outputSchema,
		})
	}
	writeJSON(w, http.StatusOK, out)
}

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

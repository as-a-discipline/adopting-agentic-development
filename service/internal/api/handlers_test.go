package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"pulse/generated"
	"pulse/internal/config"
	"pulse/internal/monitoring"
)

func newTestServer(t *testing.T) http.Handler {
	t.Helper()
	checker := monitoring.NewChecker([]config.Service{
		{ID: "svc-a", Name: "Service A", URL: "http://example.invalid"},
	}, nil)
	handler := NewHandler(checker, nil)
	return generated.Handler(handler)
}

func TestListServices(t *testing.T) {
	srv := newTestServer(t)

	req := httptest.NewRequest(http.MethodGet, "/services", nil)
	rec := httptest.NewRecorder()
	srv.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}

	var services []generated.MonitoredService
	if err := json.Unmarshal(rec.Body.Bytes(), &services); err != nil {
		t.Fatalf("invalid JSON response: %v", err)
	}
	if len(services) != 1 || services[0].Id != "svc-a" {
		t.Errorf("expected one service 'svc-a', got %+v", services)
	}
}

func TestGetService_Found(t *testing.T) {
	srv := newTestServer(t)

	req := httptest.NewRequest(http.MethodGet, "/services/svc-a", nil)
	rec := httptest.NewRecorder()
	srv.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestGetService_NotFound(t *testing.T) {
	srv := newTestServer(t)

	req := httptest.NewRequest(http.MethodGet, "/services/does-not-exist", nil)
	rec := httptest.NewRecorder()
	srv.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d: %s", rec.Code, rec.Body.String())
	}

	var errResp generated.ErrorResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &errResp); err != nil {
		t.Fatalf("invalid JSON error response: %v", err)
	}
	if errResp.Code != "service_not_found" {
		t.Errorf("expected code 'service_not_found', got %q", errResp.Code)
	}
}

func TestCheckService_NotFound(t *testing.T) {
	srv := newTestServer(t)

	req := httptest.NewRequest(http.MethodPost, "/services/does-not-exist/check", nil)
	rec := httptest.NewRecorder()
	srv.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestCheckService_Found(t *testing.T) {
	srv := newTestServer(t)

	req := httptest.NewRequest(http.MethodPost, "/services/svc-a/check", nil)
	rec := httptest.NewRecorder()
	srv.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}

	var svc generated.MonitoredService
	if err := json.Unmarshal(rec.Body.Bytes(), &svc); err != nil {
		t.Fatalf("invalid JSON response: %v", err)
	}
	// example.invalid never resolves, so the check must report "down".
	if svc.Status != generated.Down {
		t.Errorf("expected status %q after checking an unreachable URL, got %q", generated.Down, svc.Status)
	}
}

package monitoring

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"pulse/generated"
	"pulse/internal/config"
	"pulse/internal/plugins"
	"pulse/internal/plugins/httpcheck"
)

func testRegistry() *plugins.Registry {
	r := plugins.NewRegistry()
	r.Register(httpcheck.Type, httpcheck.Factory)
	return r
}

func newTestChecker(t *testing.T, url string) *Checker {
	t.Helper()
	services := []config.Service{
		{ID: "svc-a", Name: "Service A", URL: url},
	}
	return NewChecker(services, testRegistry())
}

func TestNewChecker_SeedsUnknownStatus(t *testing.T) {
	c := newTestChecker(t, "http://example.invalid")

	svc, err := c.Get("svc-a")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if svc.Status != generated.Unknown {
		t.Errorf("expected status %q, got %q", generated.Unknown, svc.Status)
	}
	if svc.LastChecked != nil {
		t.Errorf("expected LastChecked to be nil before any check, got %v", svc.LastChecked)
	}
}

func TestGet_UnknownID(t *testing.T) {
	c := newTestChecker(t, "http://example.invalid")

	_, err := c.Get("does-not-exist")
	if !errors.Is(err, ErrNotFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestCheck_UnknownID(t *testing.T) {
	c := newTestChecker(t, "http://example.invalid")

	_, err := c.Check(context.Background(), "does-not-exist")
	if !errors.Is(err, ErrNotFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestCheck_Healthy(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	c := newTestChecker(t, server.URL)

	svc, err := c.Check(context.Background(), "svc-a")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if svc.Status != generated.Healthy {
		t.Errorf("expected status %q, got %q (message: %v)", generated.Healthy, svc.Status, svc.Message)
	}
	if svc.LastChecked == nil {
		t.Error("expected LastChecked to be set after a check")
	}
	if svc.ResponseTimeMs == nil {
		t.Error("expected ResponseTimeMs to be set after a check")
	}
	if svc.Message != nil {
		t.Errorf("expected no message on a healthy check, got %q", *svc.Message)
	}
}

func TestCheck_Degraded(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(HealthyThreshold + 200*time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	c := newTestChecker(t, server.URL)

	svc, err := c.Check(context.Background(), "svc-a")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if svc.Status != generated.Degraded {
		t.Errorf("expected status %q, got %q", generated.Degraded, svc.Status)
	}
	if svc.ResponseTimeMs == nil || *svc.ResponseTimeMs < HealthyThreshold.Milliseconds() {
		t.Errorf("expected ResponseTimeMs above the healthy threshold, got %v", svc.ResponseTimeMs)
	}
}

func TestCheck_DownOnErrorStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	c := newTestChecker(t, server.URL)

	svc, err := c.Check(context.Background(), "svc-a")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if svc.Status != generated.Down {
		t.Errorf("expected status %q, got %q", generated.Down, svc.Status)
	}
	if svc.Message == nil || *svc.Message == "" {
		t.Error("expected a message explaining the down status")
	}
}

func TestCheck_DownOnNetworkFailure(t *testing.T) {
	// A closed server guarantees connection refused rather than depending on
	// external network state.
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	url := server.URL
	server.Close()

	c := newTestChecker(t, url)

	svc, err := c.Check(context.Background(), "svc-a")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if svc.Status != generated.Down {
		t.Errorf("expected status %q, got %q", generated.Down, svc.Status)
	}
}

func TestList_PreservesSeedOrder(t *testing.T) {
	services := []config.Service{
		{ID: "b", Name: "B", URL: "http://example.invalid"},
		{ID: "a", Name: "A", URL: "http://example.invalid"},
	}
	c := NewChecker(services, testRegistry())

	list := c.List()
	if len(list) != 2 || list[0].Id != "b" || list[1].Id != "a" {
		t.Errorf("expected seed order [b, a], got %+v", list)
	}
}

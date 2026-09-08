package httpcheck

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"pulse/internal/plugins"
)

func TestPlugin_ImplementsInterface(t *testing.T) {
	p := Factory()
	if p.Type() != Type {
		t.Errorf("expected type %q, got %q", Type, p.Type())
	}
}

func TestPlugin_InputSchema_ValidatesRequiredURL(t *testing.T) {
	p := Factory()

	if err := plugins.Validate(p.InputSchema(), json.RawMessage(`{"url":"http://example.invalid"}`)); err != nil {
		t.Errorf("expected valid input to pass schema validation: %v", err)
	}
	if err := plugins.Validate(p.InputSchema(), json.RawMessage(`{}`)); err == nil {
		t.Error("expected missing url to fail schema validation")
	}
	if err := plugins.Validate(p.InputSchema(), json.RawMessage(`{"url":"x","extra":true}`)); err == nil {
		t.Error("expected an unknown property to fail schema validation (additionalProperties: false)")
	}
}

func TestPlugin_Check_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	p := Factory()
	rawOut, err := p.Check(context.Background(), json.RawMessage(`{"url":"`+server.URL+`"}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := plugins.Validate(p.OutputSchema(), rawOut); err != nil {
		t.Errorf("output failed schema validation: %v", err)
	}

	var out struct {
		Success    bool  `json:"success"`
		StatusCode int   `json:"statusCode"`
		ElapsedMs  int64 `json:"elapsedMs"`
	}
	if err := json.Unmarshal(rawOut, &out); err != nil {
		t.Fatalf("invalid output JSON: %v", err)
	}
	if !out.Success {
		t.Error("expected success=true for a 200 response")
	}
	if out.StatusCode != http.StatusOK {
		t.Errorf("expected statusCode 200, got %d", out.StatusCode)
	}
}

func TestPlugin_Check_FailureStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	p := Factory()
	rawOut, err := p.Check(context.Background(), json.RawMessage(`{"url":"`+server.URL+`"}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var out struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
	}
	if err := json.Unmarshal(rawOut, &out); err != nil {
		t.Fatalf("invalid output JSON: %v", err)
	}
	if out.Success {
		t.Error("expected success=false for a 500 response")
	}
	if out.Message == "" {
		t.Error("expected a non-empty message explaining the failure")
	}
}

func TestPlugin_Check_NetworkFailure(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	url := server.URL
	server.Close()

	p := Factory()
	rawOut, err := p.Check(context.Background(), json.RawMessage(`{"url":"`+url+`"}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var out struct {
		Success bool `json:"success"`
	}
	if err := json.Unmarshal(rawOut, &out); err != nil {
		t.Fatalf("invalid output JSON: %v", err)
	}
	if out.Success {
		t.Error("expected success=false when the connection is refused")
	}
}

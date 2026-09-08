// Package httpcheck implements Pulse's "http" check-plugin: an HTTP GET
// (or HEAD) request against a URL, reporting success, elapsed time, and an
// optional message. It is behaviorally identical to the hardcoded check
// internal/monitoring performed before
// ../../../adr/0002-factory-based-plugin-model-for-checks.md — this package
// exists to prove the plugin model works, not to change check behavior.
package httpcheck

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"pulse/internal/plugins"
)

// Type is this plugin's registered type name.
const Type = "http"

// inputSchema describes the {url, method} shape Check expects. method is
// optional; when absent, GET is used (see Check).
const inputSchema = `{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "type": "object",
  "required": ["url"],
  "additionalProperties": false,
  "properties": {
    "url": {
      "type": "string",
      "format": "uri",
      "description": "URL to send the health-check request to"
    },
    "method": {
      "type": "string",
      "enum": ["GET", "HEAD"],
      "description": "HTTP method to use; defaults to GET"
    }
  }
}`

// outputSchema describes the {success, statusCode, elapsedMs, message}
// shape Check returns.
const outputSchema = `{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "type": "object",
  "required": ["success", "elapsedMs"],
  "additionalProperties": false,
  "properties": {
    "success": {
      "type": "boolean",
      "description": "Whether the request completed with a 2xx status code"
    },
    "statusCode": {
      "type": "integer",
      "description": "HTTP status code received, if any"
    },
    "elapsedMs": {
      "type": "integer",
      "minimum": 0,
      "description": "Request duration in milliseconds"
    },
    "message": {
      "type": "string",
      "description": "Human-readable detail, present when success is false"
    }
  }
}`

// input mirrors the JSON shape described by inputSchema.
type input struct {
	URL    string `json:"url"`
	Method string `json:"method,omitempty"`
}

// output mirrors the JSON shape described by outputSchema.
type output struct {
	Success    bool   `json:"success"`
	StatusCode int    `json:"statusCode,omitempty"`
	ElapsedMs  int64  `json:"elapsedMs"`
	Message    string `json:"message,omitempty"`
}

// Plugin implements plugins.Plugin for the "http" check type.
type Plugin struct {
	client *http.Client
}

var _ plugins.Plugin = (*Plugin)(nil)

// Factory constructs a new *Plugin using http.DefaultClient. The caller
// (internal/monitoring.Checker) applies its own per-check timeout via the
// context passed to Check, so this plugin does not need its own client
// timeout configuration.
func Factory() plugins.Plugin { return &Plugin{client: http.DefaultClient} }

// Type implements plugins.Plugin.
func (p *Plugin) Type() string { return Type }

// InputSchema implements plugins.Plugin.
func (p *Plugin) InputSchema() []byte { return []byte(inputSchema) }

// OutputSchema implements plugins.Plugin.
func (p *Plugin) OutputSchema() []byte { return []byte(outputSchema) }

// Check implements plugins.Plugin: performs the HTTP request and reports
// the result. rawInput is assumed to already be validated against
// InputSchema by the caller.
func (p *Plugin) Check(ctx context.Context, rawInput json.RawMessage) (json.RawMessage, error) {
	var in input
	if err := json.Unmarshal(rawInput, &in); err != nil {
		return nil, fmt.Errorf("httpcheck: invalid input: %w", err)
	}
	method := in.Method
	if method == "" {
		method = http.MethodGet
	}

	req, err := http.NewRequestWithContext(ctx, method, in.URL, nil)
	if err != nil {
		return marshalOutput(output{Success: false, Message: "invalid URL: " + err.Error()})
	}

	start := time.Now()
	resp, err := p.client.Do(req)
	elapsed := time.Since(start)
	if err != nil {
		return marshalOutput(output{
			Success:   false,
			ElapsedMs: elapsed.Milliseconds(),
			Message:   "request failed: " + err.Error(),
		})
	}
	defer resp.Body.Close()

	success := resp.StatusCode >= 200 && resp.StatusCode < 300
	out := output{
		Success:    success,
		StatusCode: resp.StatusCode,
		ElapsedMs:  elapsed.Milliseconds(),
	}
	if !success {
		out.Message = fmt.Sprintf("unexpected status code %d", resp.StatusCode)
	}
	return marshalOutput(out)
}

func marshalOutput(out output) (json.RawMessage, error) {
	b, err := json.Marshal(out)
	if err != nil {
		return nil, fmt.Errorf("httpcheck: marshal output: %w", err)
	}
	return b, nil
}

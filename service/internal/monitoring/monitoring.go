// Package monitoring implements Pulse's deterministic service health-check
// behavior: seeded in-memory state, plugin-based checks, and status
// normalization. See ../../adr/0001-no-persistence-in-baseline.md — there
// is no persistent store; state lives only for the process lifetime. See
// ../../adr/0002-factory-based-plugin-model-for-checks.md — the actual
// check mechanism for each service is a plugins.Plugin, resolved by type
// from a *plugins.Registry, not a hardcoded HTTP call.
package monitoring

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"pulse/generated"
	"pulse/internal/config"
	"pulse/internal/plugins"
)

// HealthyThreshold is the maximum response time for a successful check to be
// considered "healthy" rather than "degraded". Kept simple and documented
// here rather than made configurable, per the baseline's "keep it boring"
// design principle. Status normalization is a Pulse-level concept — it
// applies uniformly across every plugin type, not something each plugin
// implements itself.
const HealthyThreshold = 300 * time.Millisecond

// CheckTimeout bounds how long a single check (i.e. a single Plugin.Check
// call) may take. Every plugin must honor the context passed to Check.
const CheckTimeout = 5 * time.Second

// ErrNotFound indicates no monitored service exists with the requested ID.
var ErrNotFound = errors.New("service not found")

// pluginOutput mirrors the generic {success, elapsedMs, statusCode,
// message} shape every plugin's output conforms to. Plugin-specific fields
// beyond these four are not currently interpreted by the Checker — see
// ../../adr/0002-factory-based-plugin-model-for-checks.md's consequences.
type pluginOutput struct {
	Success    bool   `json:"success"`
	StatusCode int    `json:"statusCode,omitempty"`
	ElapsedMs  int64  `json:"elapsedMs"`
	Message    string `json:"message,omitempty"`
}

// checkEntry holds everything the Checker needs to run a single service's
// check: the constructed Plugin instance and its pre-built,
// schema-validated input.
type checkEntry struct {
	plugin plugins.Plugin
	input  json.RawMessage
}

// Checker holds the current in-memory status of a fixed set of monitored
// services and knows how to perform a check against any of them by
// delegating to that service's check-plugin. Checker is safe for concurrent
// use.
type Checker struct {
	mu      sync.RWMutex
	state   map[string]generated.MonitoredService
	order   []string // preserves deterministic seed order for listing
	entries map[string]checkEntry
}

// NewChecker builds a Checker seeded with the given services, resolving
// each service's check-plugin via registry (falling back to
// config.DefaultType when a service's Type is empty). Every service starts
// in the "unknown" status until its first check.
//
// NewChecker panics if a service names an unregistered plugin type or
// supplies input that fails that plugin's InputSchema — both are
// configuration defects that should fail fast at startup, not surface as a
// confusing runtime 404/500 on first check.
func NewChecker(services []config.Service, registry *plugins.Registry) *Checker {
	c := &Checker{
		state:   make(map[string]generated.MonitoredService, len(services)),
		order:   make([]string, 0, len(services)),
		entries: make(map[string]checkEntry, len(services)),
	}
	for _, s := range services {
		typeName := s.Type
		if typeName == "" {
			typeName = config.DefaultType
		}

		plugin, err := registry.New(typeName)
		if err != nil {
			panic(fmt.Sprintf("monitoring: service %q: %v", s.ID, err))
		}

		input := s.Config
		if len(input) == 0 {
			input = buildDefaultInput(s)
		}
		if err := plugins.Validate(plugin.InputSchema(), input); err != nil {
			panic(fmt.Sprintf("monitoring: service %q: invalid plugin input: %v", s.ID, err))
		}

		c.state[s.ID] = generated.MonitoredService{
			Id:     s.ID,
			Name:   s.Name,
			Url:    s.URL,
			Type:   typeName,
			Status: generated.Unknown,
		}
		c.order = append(c.order, s.ID)
		c.entries[s.ID] = checkEntry{plugin: plugin, input: input}
	}
	return c
}

// buildDefaultInput constructs the "http" plugin's input ({"url": ...})
// from a service's legacy URL field, so seed configurations written before
// the plugin model (just id/name/url) keep working unchanged.
func buildDefaultInput(s config.Service) json.RawMessage {
	b, err := json.Marshal(map[string]string{"url": s.URL})
	if err != nil {
		// s.URL is always a plain string; Marshal cannot fail here.
		panic(fmt.Sprintf("monitoring: building default input for %q: %v", s.ID, err))
	}
	return b
}

// List returns the current status of every monitored service, in
// deterministic (seed) order.
func (c *Checker) List() []generated.MonitoredService {
	c.mu.RLock()
	defer c.mu.RUnlock()
	out := make([]generated.MonitoredService, 0, len(c.order))
	for _, id := range c.order {
		out = append(out, c.state[id])
	}
	return out
}

// Get returns the current status of a single monitored service.
func (c *Checker) Get(id string) (generated.MonitoredService, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	s, ok := c.state[id]
	if !ok {
		return generated.MonitoredService{}, ErrNotFound
	}
	return s, nil
}

// Check performs an immediate check against the given service (via its
// resolved plugin), updates its in-memory status, and returns the
// resulting status. It does not persist any history (no database) — only
// the latest result is kept.
func (c *Checker) Check(ctx context.Context, id string) (generated.MonitoredService, error) {
	c.mu.RLock()
	svc, ok := c.state[id]
	entry, entryOK := c.entries[id]
	c.mu.RUnlock()
	if !ok || !entryOK {
		return generated.MonitoredService{}, ErrNotFound
	}

	ctx, cancel := context.WithTimeout(ctx, CheckTimeout)
	defer cancel()

	updated := c.performCheck(ctx, svc, entry)

	c.mu.Lock()
	c.state[id] = updated
	c.mu.Unlock()

	return updated, nil
}

func (c *Checker) performCheck(ctx context.Context, svc generated.MonitoredService, entry checkEntry) generated.MonitoredService {
	rawOutput, err := entry.plugin.Check(ctx, entry.input)
	if err != nil {
		return normalize(svc, 0, false, "plugin error: "+err.Error())
	}
	if err := plugins.Validate(entry.plugin.OutputSchema(), rawOutput); err != nil {
		return normalize(svc, 0, false, "plugin returned invalid output: "+err.Error())
	}

	var out pluginOutput
	if err := json.Unmarshal(rawOutput, &out); err != nil {
		return normalize(svc, 0, false, "plugin output unmarshal failed: "+err.Error())
	}

	return normalize(svc, time.Duration(out.ElapsedMs)*time.Millisecond, out.Success, out.Message)
}

// normalize applies Pulse's deterministic status rules:
//
//   - success is false                              -> down
//   - success is true, within HealthyThreshold       -> healthy
//   - success is true, slower than HealthyThreshold  -> degraded
func normalize(svc generated.MonitoredService, elapsed time.Duration, success bool, message string) generated.MonitoredService {
	now := time.Now().UTC()
	svc.LastChecked = &now

	ms := elapsed.Milliseconds()
	svc.ResponseTimeMs = &ms

	switch {
	case !success:
		svc.Status = generated.Down
	case elapsed <= HealthyThreshold:
		svc.Status = generated.Healthy
	default:
		svc.Status = generated.Degraded
	}

	if message != "" {
		svc.Message = &message
	} else {
		svc.Message = nil
	}
	return svc
}

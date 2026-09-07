// Package monitoring implements Pulse's deterministic service health-check
// behavior: seeded in-memory state, HTTP checks with explicit timeouts, and
// status normalization. See ../../adr/0001-no-persistence-in-baseline.md —
// there is no persistent store; state lives only for the process lifetime.
package monitoring

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"pulse/generated"
	"pulse/internal/config"
)

// HealthyThreshold is the maximum response time for a successful check to be
// considered "healthy" rather than "degraded". Kept simple and documented
// here rather than made configurable, per the baseline's "keep it boring"
// design principle.
const HealthyThreshold = 300 * time.Millisecond

// CheckTimeout bounds how long a single outbound health check may take.
// Every outbound HTTP call Pulse makes honors this timeout and the caller's
// context, per service/AGENTS.md.
const CheckTimeout = 5 * time.Second

// ErrNotFound indicates no monitored service exists with the requested ID.
var ErrNotFound = errors.New("service not found")

// Checker holds the current in-memory status of a fixed set of monitored
// services and knows how to perform an HTTP health check against any of
// them. Checker is safe for concurrent use.
type Checker struct {
	client *http.Client

	mu    sync.RWMutex
	state map[string]generated.MonitoredService
	order []string // preserves deterministic seed order for listing
}

// NewChecker builds a Checker seeded with the given services. Every service
// starts in the "unknown" status until its first check. A nil client uses a
// default *http.Client configured with CheckTimeout.
func NewChecker(services []config.Service, client *http.Client) *Checker {
	if client == nil {
		client = &http.Client{Timeout: CheckTimeout}
	}
	c := &Checker{
		client: client,
		state:  make(map[string]generated.MonitoredService, len(services)),
		order:  make([]string, 0, len(services)),
	}
	for _, s := range services {
		c.state[s.ID] = generated.MonitoredService{
			Id:     s.ID,
			Name:   s.Name,
			Url:    s.URL,
			Status: generated.Unknown,
		}
		c.order = append(c.order, s.ID)
	}
	return c
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

// Check performs an immediate HTTP health check against the given service's
// URL, updates its in-memory status, and returns the resulting status. It
// does not persist any history (no database) — only the latest result is
// kept.
func (c *Checker) Check(ctx context.Context, id string) (generated.MonitoredService, error) {
	c.mu.RLock()
	svc, ok := c.state[id]
	c.mu.RUnlock()
	if !ok {
		return generated.MonitoredService{}, ErrNotFound
	}

	ctx, cancel := context.WithTimeout(ctx, CheckTimeout)
	defer cancel()

	updated := c.performCheck(ctx, svc)

	c.mu.Lock()
	c.state[id] = updated
	c.mu.Unlock()

	return updated, nil
}

func (c *Checker) performCheck(ctx context.Context, svc generated.MonitoredService) generated.MonitoredService {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, svc.Url, nil)
	if err != nil {
		return normalize(svc, 0, false, "invalid service URL: "+err.Error())
	}

	start := time.Now()
	resp, err := c.client.Do(req)
	elapsed := time.Since(start)
	if err != nil {
		return normalize(svc, elapsed, false, "request failed: "+err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return normalize(svc, elapsed, false, fmt.Sprintf("unexpected status code %d", resp.StatusCode))
	}

	return normalize(svc, elapsed, true, "")
}

// normalize applies Pulse's deterministic status rules:
//
//   - request failed, invalid URL, or non-2xx response -> down
//   - succeeded within HealthyThreshold                -> healthy
//   - succeeded slower than HealthyThreshold            -> degraded
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

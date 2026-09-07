// Package config provides the deterministic set of monitored services Pulse
// starts with. Pulse has no persistence (see
// ../../adr/0001-no-persistence-in-baseline.md); this seed list is the only
// source of monitored-service configuration.
package config

import (
	"encoding/json"
	"log/slog"
	"os"
)

// Service is a monitored service's static configuration.
type Service struct {
	ID   string
	Name string
	URL  string
}

// seedEnvVar, when set to a JSON array of {"id","name","url"} objects,
// overrides the default seed list below. The Docker Compose integration
// environment (deploy/compose/docker-compose.yml) sets this to point at
// deterministic fake HTTP targets instead of the public internet, so
// integration tests never depend on real network state.
const seedEnvVar = "PULSE_SEED_SERVICES_JSON"

// Seed returns the deterministic, ordered list of services Pulse monitors.
//
// If PULSE_SEED_SERVICES_JSON is set to valid JSON, it is used verbatim (in
// array order). Otherwise, a default placeholder list (suitable for local,
// non-integration-tested use) is returned.
func Seed() []Service {
	if raw := os.Getenv(seedEnvVar); raw != "" {
		var services []Service
		if err := json.Unmarshal([]byte(raw), &services); err != nil {
			slog.Error("invalid "+seedEnvVar+", falling back to default seed", "error", err)
			return defaultSeed()
		}
		return services
	}
	return defaultSeed()
}

func defaultSeed() []Service {
	return []Service{
		{ID: "web-app", Name: "Web App", URL: "https://httpbin.org/status/200"},
		{ID: "billing-api", Name: "Billing API", URL: "https://httpbin.org/delay/1"},
		{ID: "legacy-service", Name: "Legacy Service", URL: "https://httpbin.org/status/500"},
	}
}

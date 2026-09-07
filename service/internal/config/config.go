// Package config provides the deterministic set of monitored services Pulse
// starts with. Pulse has no persistence (see
// ../../adr/0001-no-persistence-in-baseline.md); this seed list is the only
// source of monitored-service configuration.
package config

// Service is a monitored service's static configuration.
type Service struct {
	ID   string
	Name string
	URL  string
}

// Seed returns the deterministic, ordered list of services Pulse monitors.
//
// These placeholder URLs are suitable for local/manual use; the Docker
// Compose integration environment (Session 3) overrides this with
// deterministic fake HTTP targets so integration tests never depend on the
// public internet.
func Seed() []Service {
	return []Service{
		{ID: "web-app", Name: "Web App", URL: "https://httpbin.org/status/200"},
		{ID: "billing-api", Name: "Billing API", URL: "https://httpbin.org/delay/1"},
		{ID: "legacy-service", Name: "Legacy Service", URL: "https://httpbin.org/status/500"},
	}
}

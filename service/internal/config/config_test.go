package config

import (
	"reflect"
	"testing"
)

func TestSeed_DefaultWhenEnvUnset(t *testing.T) {
	t.Setenv(seedEnvVar, "")

	got := Seed()
	if len(got) == 0 {
		t.Fatalf("expected a non-empty default seed list")
	}
	if got[0].ID != "web-app" {
		t.Fatalf("expected default seed order to start with 'web-app', got %q", got[0].ID)
	}
}

func TestSeed_OverriddenByEnv(t *testing.T) {
	t.Setenv(seedEnvVar, `[{"id":"a","name":"A","url":"http://a.invalid"},{"id":"b","name":"B","url":"http://b.invalid"}]`)

	want := []Service{
		{ID: "a", Name: "A", URL: "http://a.invalid"},
		{ID: "b", Name: "B", URL: "http://b.invalid"},
	}
	got := Seed()
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Seed() = %+v, want %+v", got, want)
	}
}

func TestSeed_FallsBackOnInvalidJSON(t *testing.T) {
	t.Setenv(seedEnvVar, "not valid json")

	got := Seed()
	if len(got) == 0 || got[0].ID != "web-app" {
		t.Fatalf("expected fallback to default seed on invalid JSON, got %+v", got)
	}
}

func TestSeed_DefaultTypeOmittedFromLegacyConfig(t *testing.T) {
	// Seed configs written before the plugin model (just id/name/url) must
	// keep working unchanged: Type/Config are simply absent.
	t.Setenv(seedEnvVar, `[{"id":"a","name":"A","url":"http://a.invalid"}]`)

	got := Seed()
	if len(got) != 1 {
		t.Fatalf("expected one service, got %+v", got)
	}
	if got[0].Type != "" {
		t.Errorf("expected empty Type on a legacy seed entry, got %q", got[0].Type)
	}
	if len(got[0].Config) != 0 {
		t.Errorf("expected empty Config on a legacy seed entry, got %q", got[0].Config)
	}
}

func TestSeed_ExplicitTypeAndConfig(t *testing.T) {
	t.Setenv(seedEnvVar, `[{"id":"a","name":"A","url":"http://a.invalid","type":"http","config":{"url":"http://override.invalid"}}]`)

	got := Seed()
	if len(got) != 1 {
		t.Fatalf("expected one service, got %+v", got)
	}
	if got[0].Type != "http" {
		t.Errorf("expected Type %q, got %q", "http", got[0].Type)
	}
	if string(got[0].Config) != `{"url":"http://override.invalid"}` {
		t.Errorf("expected Config to round-trip verbatim, got %q", got[0].Config)
	}
}

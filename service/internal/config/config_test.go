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

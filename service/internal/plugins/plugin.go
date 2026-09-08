// Package plugins defines Pulse's factory-based check-plugin model: a
// Plugin interface, a Factory type for constructing them, and a Registry
// that maps a plugin type name to its Factory. See
// ../../adr/0002-factory-based-plugin-model-for-checks.md for the
// architectural decision this package implements.
//
// Every Plugin must ship JSON Schema documents describing its input
// (configuration) and output (result) shapes — this is what makes a check
// type self-describing and mechanically validatable, not just documented in
// Go source. Validation itself (input before Check runs, output before the
// result is trusted) is performed by the caller (internal/monitoring) using
// the Validator in this package, so every plugin gets the same enforcement
// without having to implement it itself.
package plugins

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
)

// Plugin implements a single check mechanism (e.g. an HTTP GET). A Plugin
// instance is stateless between calls to Check — any per-check
// configuration arrives via the input parameter, not stored on the Plugin.
type Plugin interface {
	// Type returns this plugin's registered type name, e.g. "http". It must
	// match the name it was registered under.
	Type() string

	// InputSchema returns the JSON Schema (as raw bytes) describing the
	// shape Check's input parameter must conform to.
	InputSchema() []byte

	// OutputSchema returns the JSON Schema (as raw bytes) describing the
	// shape Check's returned output conforms to.
	OutputSchema() []byte

	// Check performs a single check using the given input (already
	// validated by the caller against InputSchema) and returns a result
	// conforming to OutputSchema. Check must honor ctx cancellation/timeout.
	Check(ctx context.Context, input json.RawMessage) (json.RawMessage, error)
}

// Factory constructs a new Plugin instance. Factories are registered by
// type name with a Registry; the Registry is the "factory" in this
// package's factory-based plugin model — it holds named constructors and
// builds instances on demand.
type Factory func() Plugin

// Registry maps plugin type names to Factories and constructs Plugin
// instances by name. Registry is safe for concurrent use.
type Registry struct {
	mu        sync.RWMutex
	factories map[string]Factory
}

// NewRegistry returns an empty Registry.
func NewRegistry() *Registry {
	return &Registry{factories: make(map[string]Factory)}
}

// Register adds a Factory under the given type name. Registration is
// explicit — callers (e.g. cmd/pulse/main.go) must call Register themselves;
// this package never registers a plugin as a side effect of being imported
// (no init()-based self-registration), per
// ../../adr/0002-factory-based-plugin-model-for-checks.md.
//
// Register panics if typeName is empty or already registered — both are
// programmer errors caught at startup, not runtime conditions to recover
// from.
func (r *Registry) Register(typeName string, factory Factory) {
	if typeName == "" {
		panic("plugins: Register called with empty type name")
	}
	if factory == nil {
		panic("plugins: Register called with nil factory for type " + typeName)
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.factories[typeName]; exists {
		panic("plugins: type already registered: " + typeName)
	}
	r.factories[typeName] = factory
}

// New constructs a fresh Plugin instance for the given type name, or an
// error if no such type is registered.
func (r *Registry) New(typeName string) (Plugin, error) {
	r.mu.RLock()
	factory, ok := r.factories[typeName]
	r.mu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("plugins: unknown plugin type %q", typeName)
	}
	return factory(), nil
}

// List returns every registered plugin type's descriptor (type name, input
// schema, output schema), sorted by type name. Used by
// internal/api.Handler.ListPluginTypes to serve GET /plugin-types.
func (r *Registry) List() []Descriptor {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]Descriptor, 0, len(r.factories))
	for typeName, factory := range r.factories {
		p := factory()
		out = append(out, Descriptor{
			Type:         typeName,
			InputSchema:  p.InputSchema(),
			OutputSchema: p.OutputSchema(),
		})
	}
	sortDescriptors(out)
	return out
}

// Descriptor describes a registered plugin type without requiring callers
// to construct or hold a live Plugin instance.
type Descriptor struct {
	Type         string
	InputSchema  []byte
	OutputSchema []byte
}

func sortDescriptors(d []Descriptor) {
	for i := 1; i < len(d); i++ {
		for j := i; j > 0 && d[j-1].Type > d[j].Type; j-- {
			d[j-1], d[j] = d[j], d[j-1]
		}
	}
}

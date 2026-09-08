package plugins

import (
	"context"
	"encoding/json"
	"testing"
)

// fakePlugin is a minimal Plugin used only to exercise Registry behavior
// without depending on internal/plugins/httpcheck.
type fakePlugin struct{ typeName string }

func (p *fakePlugin) Type() string         { return p.typeName }
func (p *fakePlugin) InputSchema() []byte  { return []byte(`{"type":"object"}`) }
func (p *fakePlugin) OutputSchema() []byte { return []byte(`{"type":"object"}`) }
func (p *fakePlugin) Check(ctx context.Context, input json.RawMessage) (json.RawMessage, error) {
	return json.RawMessage(`{}`), nil
}

func TestRegistry_RegisterAndNew(t *testing.T) {
	r := NewRegistry()
	r.Register("fake", func() Plugin { return &fakePlugin{typeName: "fake"} })

	p, err := r.New("fake")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.Type() != "fake" {
		t.Errorf("expected type %q, got %q", "fake", p.Type())
	}
}

func TestRegistry_New_UnknownType(t *testing.T) {
	r := NewRegistry()

	if _, err := r.New("does-not-exist"); err == nil {
		t.Error("expected an error for an unregistered type, got nil")
	}
}

func TestRegistry_Register_PanicsOnEmptyType(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Error("expected a panic for an empty type name")
		}
	}()
	NewRegistry().Register("", func() Plugin { return &fakePlugin{} })
}

func TestRegistry_Register_PanicsOnNilFactory(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Error("expected a panic for a nil factory")
		}
	}()
	NewRegistry().Register("fake", nil)
}

func TestRegistry_Register_PanicsOnDuplicate(t *testing.T) {
	r := NewRegistry()
	r.Register("fake", func() Plugin { return &fakePlugin{typeName: "fake"} })

	defer func() {
		if recover() == nil {
			t.Error("expected a panic when registering a duplicate type")
		}
	}()
	r.Register("fake", func() Plugin { return &fakePlugin{typeName: "fake"} })
}

func TestRegistry_List_SortedByType(t *testing.T) {
	r := NewRegistry()
	r.Register("zeta", func() Plugin { return &fakePlugin{typeName: "zeta"} })
	r.Register("alpha", func() Plugin { return &fakePlugin{typeName: "alpha"} })

	got := r.List()
	if len(got) != 2 || got[0].Type != "alpha" || got[1].Type != "zeta" {
		t.Fatalf("expected sorted [alpha, zeta], got %+v", got)
	}
}

func TestValidate_Success(t *testing.T) {
	schema := []byte(`{"type":"object","required":["url"],"properties":{"url":{"type":"string"}}}`)
	data := json.RawMessage(`{"url":"http://example.invalid"}`)

	if err := Validate(schema, data); err != nil {
		t.Errorf("expected valid data to pass, got error: %v", err)
	}
}

func TestValidate_Failure(t *testing.T) {
	schema := []byte(`{"type":"object","required":["url"],"properties":{"url":{"type":"string"}}}`)
	data := json.RawMessage(`{}`)

	if err := Validate(schema, data); err == nil {
		t.Error("expected missing required field to fail validation, got nil")
	}
}

func TestValidate_InvalidJSON(t *testing.T) {
	schema := []byte(`{"type":"object"}`)
	data := json.RawMessage(`not json`)

	if err := Validate(schema, data); err == nil {
		t.Error("expected invalid JSON data to fail validation, got nil")
	}
}

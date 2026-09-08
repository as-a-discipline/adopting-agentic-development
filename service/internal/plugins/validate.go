package plugins

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/santhosh-tekuri/jsonschema/v5"
)

// Validate checks that data conforms to the given JSON Schema document. It
// is used by internal/monitoring to validate a plugin's input before
// Check runs and its output before the result is trusted — real,
// runtime-enforced validation, not just documentation, per
// ../../adr/0002-factory-based-plugin-model-for-checks.md.
func Validate(schema []byte, data json.RawMessage) error {
	compiler := jsonschema.NewCompiler()
	if err := compiler.AddResource("schema.json", bytes.NewReader(schema)); err != nil {
		return fmt.Errorf("plugins: invalid schema: %w", err)
	}
	compiled, err := compiler.Compile("schema.json")
	if err != nil {
		return fmt.Errorf("plugins: invalid schema: %w", err)
	}

	var v any
	if err := json.Unmarshal(data, &v); err != nil {
		return fmt.Errorf("plugins: invalid JSON: %w", err)
	}
	if err := compiled.Validate(v); err != nil {
		return fmt.Errorf("plugins: validation failed: %w", err)
	}
	return nil
}

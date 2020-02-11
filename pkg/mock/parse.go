package mock

import (
	"bytes"
	"strings"
)

// replacePathVars takes a path and matches it to an API endpoint. It iterates over
// currently configured Transformers, and looks for a match between an existing Transformer
// value and a value in the path. If it finds a match, it replaces that value in the path
// with the API symbol, e.g. `:org`, `:username`.
func replacePathVars(i string, ms *MockServer) string {
	parts := strings.Split(i, "/")

	for _, t := range ms.transformers {
		switch tr := t.(type) {
		case *VariableSubstitution:
			for idx, p := range parts {
				if p == tr.value {
					var buffer bytes.Buffer
					buffer.WriteString(":")
					buffer.WriteString(strings.ToLower(tr.key))
					parts[idx] = buffer.String()
				}
			}
		}
	}

	return strings.Join(parts, "/")
}

package mock

import (
	"fmt"
	"net/url"
	"strings"
)

// replacePathVars takes a path and matches it to an API endpoint. It iterates over
// currently configured Transformers, and looks for a match between an existing Transformer
// value and a value in the path. If it finds a match, it replaces that value in the path
// with the API symbol, e.g. `:org`, `:username`.
func (ms *MockServer) replacePathVars(u *url.URL) string {
	parts := strings.Split(u.EscapedPath(), "/")

	for _, t := range ms.transformers {
		switch tr := t.(type) {
		case *VariableSubstitution:
			search := tr.value
			for idx, p := range parts {
				p, _ = url.PathUnescape(p)
				if p == search {
					parts[idx] = fmt.Sprintf(":%s", strings.ToLower(tr.key))
				}
			}
		}
	}

	return strings.Join(parts, "/")
}

package mock

import (
	"io"
	"io/ioutil"
	"strings"
	"text/template"
	templateparse "text/template/parse"
)

// VariableSubstitution represents a single Golang type template value to be
// replaced with a given value.
type VariableSubstitution struct {
	key   string
	value string
}

// NewVariableSubstitution is a creator for a new VariableSubstitution.
func NewVariableSubstitution(key, value string) (*VariableSubstitution, error) {
	return &VariableSubstitution{key: key, value: value}, nil
}

// Transform is used to implement the Transformer interface. It takes an input
// Reader, substitutes the "key" with the "value" using Golang templates and
// returns a Reader that has that substitution performed.
func (vs *VariableSubstitution) Transform(in io.Reader) (io.Reader, error) {
	subMap := map[string]string{
		vs.key: vs.value,
	}

	b, err := ioutil.ReadAll(in)
	if err != nil {
		return nil, err
	}

	tmpl, err := template.New("var-substitution").Option("missingkey=error").Parse(string(b))
	if err != nil {
		return nil, err
	}

	// I'm pretty surprised this isn't a template option by default, but:
	// If the template has a substitution that is not in the map, just output
	// the original substitution string, allowing these transforms to be
	// chained together.
	for _, n := range tmpl.Root.Nodes {
		if n.Type() == templateparse.NodeAction {
			key := strings.TrimFunc(n.String(), func(r rune) bool {
				return r == '{' || r == '}' || r == ' ' || r == '.'
			})

			_, ok := subMap[key]
			if !ok {
				subMap[key] = n.String()
			}
		}
	}

	pr, pw := io.Pipe()
	go func() {
		_ = pw.CloseWithError(tmpl.Execute(pw, subMap))
	}()
	return pr, nil
}

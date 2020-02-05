package mock

import (
	"io"
	"io/ioutil"
	"strings"
	"text/template"
	templateparse "text/template/parse"
)

type VariableSubstitution struct {
	key   string
	value string
}

func NewVariableSubstitution(key, value string) (*VariableSubstitution, error) {
	return &VariableSubstitution{key: key, value: value}, nil
}

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

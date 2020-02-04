package mock

import (
	"io"
	"io/ioutil"
	"text/template"
)

type VariableSubstitution struct {
	key   string
	value string
}

func NewVariableSubstitution(key, value string) (*VariableSubstitution, error) {
	return &VariableSubstitution{key: key, value: value}, nil
}

func (vs *VariableSubstitution) Transform(in io.ReadCloser) (io.ReadCloser, error) {
	subMap := map[string]string{
		vs.key: vs.value,
	}

	b, err := ioutil.ReadAll(in)
	if err != nil {
		return nil, err
	}

	tmpl, err := template.New("var-substitution").Parse(string(b))
	if err != nil {
		return nil, err
	}

	pr, pw := io.Pipe()
	go func() {
		pw.CloseWithError(tmpl.Execute(pw, subMap))
	}()
	return pr, nil
}

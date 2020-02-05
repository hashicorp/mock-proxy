package mock

import (
	"io"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewVariableSubstitution(t *testing.T) {
	tcs := []struct {
		name  string
		key   string
		value string
		want  *VariableSubstitution
	}{
		{
			name:  "simple",
			key:   "test_key",
			value: "test_value",
			want: &VariableSubstitution{
				key:   "test_key",
				value: "test_value",
			},
		},
	}

	for _, tc := range tcs {
		tc := tc // capture range variable
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got, err := NewVariableSubstitution(tc.key, tc.value)
			require.Nil(t, err)

			assert.Equal(t, tc.want, got)
		})
	}
}

func TestVariableSubstitutionTransform(t *testing.T) {
	tcs := []struct {
		name  string
		vs    []*VariableSubstitution
		input string
		want  string
	}{
		{
			name: "no handlebars in input",
			vs: []*VariableSubstitution{
				&VariableSubstitution{
					key:   "test_key",
					value: "test_value",
				},
			},
			input: "just some input!",
			want:  "just some input!",
		},
		{
			name: "do a substitution",
			vs: []*VariableSubstitution{
				&VariableSubstitution{
					key:   "test_key",
					value: "output",
				},
			},
			input: "just some {{ .test_key }}!",
			want:  "just some output!",
		},
		{
			name: "missing key",
			vs: []*VariableSubstitution{
				&VariableSubstitution{key: "a", value: "b"},
			},
			input: "just some {{.cool_key}}!",
			want:  "just some {{.cool_key}}!",
		},
		{
			name: "chain some transforms together",
			vs: []*VariableSubstitution{
				&VariableSubstitution{key: "a", value: "b"},
				&VariableSubstitution{key: "c", value: "d"},
			},
			input: "transform {{ .a }} and {{ .c }}",
			want:  "transform b and d",
		},
	}

	for _, tc := range tcs {
		tc := tc // capture range variable
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ir := strings.NewReader(tc.input)

			var gotReader io.Reader
			gotReader = ir

			for _, transform := range tc.vs {
				gr, err := transform.Transform(gotReader)
				require.Nil(t, err)
				gotReader = gr
			}

			gotBytes, err := ioutil.ReadAll(gotReader)
			require.Nil(t, err)

			got := string(gotBytes)
			assert.Equal(t, tc.want, got)
		})
	}
}

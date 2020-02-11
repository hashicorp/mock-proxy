package mock

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReplacePathVars(t *testing.T) {
	tcs := []struct {
		name    string
		input   string
		options []Option
		want    string
	}{
		{
			name:  "even",
			input: "/orgs/hashicorp/repos/atlas",
			options: []Option{
				WithMockRoot("testdata/"),
				WithDefaultVariables(
					&VariableSubstitution{key: "Org", value: "hashicorp"},
					&VariableSubstitution{key: "Repo", value: "atlas"},
				),
			},
			want: "/orgs/:org/repos/:repo",
		},
		{
			name:  "odd",
			input: "/orgs/hashicorp/repos",
			options: []Option{
				WithMockRoot("testdata/"),
				WithDefaultVariables(
					&VariableSubstitution{key: "Org", value: "hashicorp"},
				),
			},
			want: "/orgs/:org/repos",
		},
		{
			name:  "sequential",
			input: "/orgs/rae/atlas/teams",
			options: []Option{
				WithMockRoot("testdata/"),
				WithDefaultVariables(
					&VariableSubstitution{key: "Owner", value: "rae"},
					&VariableSubstitution{key: "Repo", value: "atlas"},
				),
			},
			want: "/orgs/:owner/:repo/teams",
		},
		{
			name:  "none",
			input: "/user/repositories",
			options: []Option{
				WithMockRoot("testdata/"),
				WithDefaultVariables(),
			},
			want: "/user/repositories",
		},
	}

	for _, tc := range tcs {
		tc := tc // capture range variable
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ms, err := NewMockServer(tc.options...)
			require.Nil(t, err)

			got := ms.replacePathVars(tc.input)

			assert.Equal(t, tc.want, got)
		})
	}
}

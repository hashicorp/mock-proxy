package mock

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParsePath(t *testing.T) {
	tcs := []struct {
		name     string
		input    string
		want     string
		wantSubs []*VariableSubstitution
	}{
		{
			name:  "even",
			input: "/orgs/~~org~hashicorp~~/repos/~~repo~atlas~~",
			want:  "/orgs/:org/repos/:repo",
			wantSubs: []*VariableSubstitution{
				{key: "org", value: "hashicorp"},
				{key: "repo", value: "atlas"},
			},
		},
		{
			name:  "odd",
			input: "/orgs/~~org~hashicorp~~/repos",
			want:  "/orgs/:org/repos",
			wantSubs: []*VariableSubstitution{
				{key: "org", value: "hashicorp"},
			},
		},
		{
			name:  "sequential",
			input: "/orgs/~~owner~rae~~/~~repo~atlas~~/teams",
			want:  "/orgs/:owner/:repo/teams",
			wantSubs: []*VariableSubstitution{
				{key: "owner", value: "rae"},
				{key: "repo", value: "atlas"},
			},
		},
		{
			name:     "none",
			input:    "/user/repositories",
			want:     "/user/repositories",
			wantSubs: []*VariableSubstitution{},
		},
	}

	for _, tc := range tcs {
		tc := tc // capture range variable
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ms, err := NewMockServer(WithMockRoot("testdata/"))
			require.Nil(t, err)

			got, gotSubs := ms.parsePath(tc.input)

			assert.Equal(t, tc.want, got)
			assert.Equal(t, tc.wantSubs, gotSubs)
		})
	}
}

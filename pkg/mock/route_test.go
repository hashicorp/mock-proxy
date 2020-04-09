package mock

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseURL(t *testing.T) {
	tcs := []struct {
		name             string
		route            *Route
		url              string
		wantPath         string
		wantTransformers []Transformer
	}{
		{
			name: "index",
			route: &Route{
				Host: "example.com",
				Path: "/",
				Type: "http",
			},
			url:              "http://example.com",
			wantPath:         "example.com/index.mock",
			wantTransformers: []Transformer{},
		},
		{
			name: "simple",
			route: &Route{
				Host: "example.com",
				Path: "/foo/bar",
				Type: "http",
			},
			url:              "http://example.com/foo/bar",
			wantPath:         "example.com/foo/bar.mock",
			wantTransformers: []Transformer{},
		},
		{
			name: "with one transform",
			route: &Route{
				Host: "example.com",
				Path: "/users/:foo/bars",
				Type: "http",
			},
			url:      "http://example.com/users/russell/bars",
			wantPath: "example.com/users/:foo/bars.mock",
			wantTransformers: []Transformer{
				&VariableSubstitution{key: "foo", value: "russell"},
			},
		},
		{
			name: "with multiple transforms",
			route: &Route{
				Host: "example.com",
				Path: "/users/:user/settings/:setting/value",
				Type: "http",
			},
			url:      "http://example.com/users/russell/settings/locale/value",
			wantPath: "example.com/users/:user/settings/:setting/value.mock",
			wantTransformers: []Transformer{
				&VariableSubstitution{key: "user", value: "russell"},
				&VariableSubstitution{key: "setting", value: "locale"},
			},
		},
		{
			name: "git",
			route: &Route{
				Host: "github.com",
				Path: "example-repo",
				Type: "git",
			},
			url:              "http://github.com/example-repo",
			wantPath:         "/git/github.com/example-repo/.git",
			wantTransformers: []Transformer{},
		},
	}

	for _, tc := range tcs {
		tc := tc // capture range variable
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			testURL, err := url.Parse(tc.url)
			require.Nil(t, err)

			gotPath, gotTransformers, err := tc.route.ParseURL(testURL)
			require.Nil(t, err)

			assert.Equal(t, tc.wantPath, gotPath)
			assert.Equal(t, tc.wantTransformers, gotTransformers)
		})
	}
}

func TestMatchRoute(t *testing.T) {
	tcs := []struct {
		name        string
		routeConfig RouteConfig
		url         string
		want        *Route
		wantErr     string
	}{
		{
			name:        "simple miss",
			routeConfig: []*Route{},
			url:         "http://example.com",
			want:        nil,
		},
		{
			name: "host miss",
			routeConfig: []*Route{
				{Host: "example.com", Path: "", Type: "http"},
			},
			url:  "http://mycoolwebsite.biz",
			want: nil,
		},
		{
			name: "simple hit",
			routeConfig: []*Route{
				{Host: "example.com", Path: "", Type: "http"},
			},
			url: "http://example.com",
			want: &Route{
				Host: "example.com",
				Path: "",
				Type: "http",
			},
		},
		{
			name: "too many hits, overlaps are bad",
			routeConfig: []*Route{
				{Host: "example.com", Path: "", Type: "http"},
				{Host: "example.com", Path: "", Type: "http"},
			},
			url:     "http://example.com",
			wantErr: "multiple routes matched input",
		},
		{
			name: "hosts and paths must both match for simple cases",
			routeConfig: []*Route{
				{Host: "example.com", Path: "/foo/bar", Type: "http"},
			},
			url:  "http://example.com/baz/bing",
			want: nil,
		},
		{
			name: "git paths match known git request patterns",
			routeConfig: []*Route{
				{Host: "github.com", Path: "/example-repo", Type: "git"},
			},
			url: "http://github.com/example-repo/info/refs?service=git-upload-pack",
			want: &Route{
				Host: "github.com",
				Path: "/example-repo",
				Type: "git",
			},
		},
		{
			name: "but not unknown ones",
			routeConfig: []*Route{
				{Host: "github.com", Path: "/example-repo", Type: "git"},
			},
			url:  "http://github.com/example-repo/otherinfo",
			want: nil,
		},
		{
			name: "or the wrong repo",
			routeConfig: []*Route{
				{Host: "github.com", Path: "/example-repo", Type: "git"},
			},
			url:  "http://github.com/other-repo/info/refs?service=git-upload-pack",
			want: nil,
		},
		{
			name: "http requests also work with substitutions logic",
			routeConfig: []*Route{
				{Host: "example.com", Path: "/users/:user/settings/:setting/value", Type: "http"},
			},
			url: "http://example.com/users/russell/settings/locale/value",
			want: &Route{
				Host: "example.com",
				Path: "/users/:user/settings/:setting/value",
				Type: "http",
			},
		},
		{
			name: "legal overlap",
			routeConfig: []*Route{
				{Host: "api.github.com", Path: "/orgs/:org", Type: "http"},
				{Host: "api.github.com", Path: "/orgs/:org/repos", Type: "http"},
				{Host: "api.github.com", Path: "/orgs/:org/repos/tree", Type: "http"},
			},
			url: "http://api.github.com/orgs/hashicorp/repos",
			want: &Route{
				Host: "api.github.com",
				Path: "/orgs/:org/repos",
				Type: "http",
			},
		},
	}

	for _, tc := range tcs {
		tc := tc // capture range variable
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			inURL, err := url.Parse(tc.url)
			require.Nil(t, err)

			got, err := tc.routeConfig.MatchRoute(inURL)
			if tc.wantErr == "" {
				require.Nil(t, err)
				assert.Equal(t, tc.want, got)
			} else {
				require.NotNil(t, err)
				assert.Contains(t, err.Error(), tc.wantErr)
			}
		})
	}
}

func TestParseRoutes(t *testing.T) {
	tcs := []struct {
		name  string
		input string
		want  RouteConfig
	}{
		{
			name:  "simple",
			input: "testdata/routes.hcl",
			want: []*Route{
				{Host: "example.com", Path: "/simple", Type: "http"},
				{Host: "example.com", Path: "/substitutions", Type: "http"},
				{Host: "example.com", Path: "/users/:name", Type: "http"},
			},
		},
	}

	for _, tc := range tcs {
		tc := tc // capture range variable
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got, err := ParseRoutes(tc.input)
			require.Nil(t, err)

			assert.Equal(t, tc.want, got)
		})
	}
}

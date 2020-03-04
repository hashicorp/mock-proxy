package mock

import (
	"bytes"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/vcs-mock-proxy/internal/cachedfs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewMockServer(t *testing.T) {
	tcs := []struct {
		name    string
		options []Option
		want    *MockServer
	}{
		{
			name: "simple",
			options: []Option{
				WithMockRoot("testdata/"),
			},
			want: &MockServer{
				apiPort:  80,
				icapPort: 11344,

				mockFilesRoot: "testdata/",

				cachedFS: &cachedfs.CachedFS{},
			},
		},
	}

	for _, tc := range tcs {
		tc := tc // capture range variable
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got, err := NewMockServer(tc.options...)
			require.Nil(t, err)

			assert.Equal(t, tc.want, got)
		})
	}
}

func TestMockServerMockHandler(t *testing.T) {
	tcs := []struct {
		name    string
		options []Option
		url     string
		want    string
	}{
		{
			name: "simple",
			url:  "example.com/simple",
			options: []Option{
				WithMockRoot("testdata/"),
			},
			want: "Hello, World!\n",
		},
		{
			name: "substitutions",
			url:  "example.com/substitutions",
			options: []Option{
				WithMockRoot("testdata/"),
				WithDefaultVariables(
					&VariableSubstitution{key: "name", value: "Davenport"},
				),
			},
			want: "Hello, Davenport!\n",
		},
	}

	for _, tc := range tcs {
		tc := tc // capture range variable
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ms, err := NewMockServer(tc.options...)
			require.Nil(t, err)

			req, err := http.NewRequest(http.MethodGet, tc.url, nil)
			require.Nil(t, err)

			recorder := httptest.NewRecorder()

			ms.mockHandler(recorder, req)

			assert.Equal(t, http.StatusOK, recorder.Result().StatusCode)

			gotBytes, err := ioutil.ReadAll(recorder.Result().Body)
			require.Nil(t, err)
			got := string(gotBytes)

			assert.Equal(t, tc.want, got)
		})
	}
}

func TestMockServerSubstitutionVariableHandler_GET(t *testing.T) {
	tcs := []struct {
		name    string
		options []Option
		want    string
	}{
		{
			name: "simple",
			options: []Option{
				WithMockRoot("testdata/"),
				WithDefaultVariables(
					&VariableSubstitution{key: "name", value: "Davenport"},
				),
			},
			want: `[{"key":"name","value":"Davenport"}]`,
		},
		{
			name: "multi",
			options: []Option{
				WithMockRoot("testdata/"),
				WithDefaultVariables(
					&VariableSubstitution{key: "name", value: "Davenport"},
					&VariableSubstitution{key: "name", value: "Barry"},
					&VariableSubstitution{key: "foo", value: "bar"},
				),
			},
			want: `[{"key":"name","value":"Barry"},{"key":"foo","value":"bar"}]`,
		},
	}

	for _, tc := range tcs {
		tc := tc // capture range variable
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ms, err := NewMockServer(tc.options...)
			require.Nil(t, err)

			req, err := http.NewRequest(http.MethodGet, "", nil)
			require.Nil(t, err)

			recorder := httptest.NewRecorder()

			ms.substitutionVariableHandler(recorder, req)

			assert.Equal(t, http.StatusOK, recorder.Result().StatusCode)

			gotBytes, err := ioutil.ReadAll(recorder.Result().Body)
			require.Nil(t, err)
			got := string(gotBytes)

			assert.Equal(t, tc.want, got)
		})
	}
}

func TestMockServerSubstitutionVariableHandler_POST(t *testing.T) {
	tcs := []struct {
		name    string
		options []Option
		key     string
		value   string
		want    []Transformer
	}{
		{
			name:  "simple",
			key:   "name",
			value: "Davenport",
			options: []Option{
				WithMockRoot("testdata/"),
			},
			want: []Transformer{
				&VariableSubstitution{key: "name", value: "Davenport"},
			},
		},
		{
			name:  "replace",
			key:   "name",
			value: "Barry",
			options: []Option{
				WithMockRoot("testdata/"),
				WithDefaultVariables(
					&VariableSubstitution{key: "name", value: "Davenport"},
				),
			},
			want: []Transformer{
				&VariableSubstitution{key: "name", value: "Barry"},
			},
		},
		{
			name:  "add",
			key:   "foo",
			value: "bar",
			options: []Option{
				WithMockRoot("testdata/"),
				WithDefaultVariables(
					&VariableSubstitution{key: "name", value: "Davenport"},
				),
			},
			want: []Transformer{
				&VariableSubstitution{key: "name", value: "Davenport"},
				&VariableSubstitution{key: "foo", value: "bar"},
			},
		},
	}

	for _, tc := range tcs {
		tc := tc // capture range variable
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ms, err := NewMockServer(tc.options...)
			require.Nil(t, err)

			var formBody bytes.Buffer
			formWriter := multipart.NewWriter(&formBody)
			_ = formWriter.WriteField("key", tc.key)
			_ = formWriter.WriteField("value", tc.value)
			formWriter.Close()

			req, err := http.NewRequest(http.MethodPost, "", &formBody)
			require.Nil(t, err)
			req.Header.Set("Content-Type", formWriter.FormDataContentType())

			recorder := httptest.NewRecorder()

			ms.substitutionVariableHandler(recorder, req)

			if !assert.Equal(t, http.StatusOK, recorder.Result().StatusCode) {
				resBytes, err := ioutil.ReadAll(recorder.Result().Body)
				require.Nil(t, err)
				res := string(resBytes)
				require.Fail(t, res)
			}

			got := ms.transformers
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestMockServerAddVariableSubstitution(t *testing.T) {
	tcs := []struct {
		name          string
		options       []Option
		substitutions []*VariableSubstitution
		want          []Transformer
	}{
		{
			name: "simple",
			options: []Option{
				WithMockRoot("testdata/"),
			},
			substitutions: []*VariableSubstitution{
				{key: "foo", value: "bar"},
			},
			want: []Transformer{
				&VariableSubstitution{key: "foo", value: "bar"},
			},
		},
		{
			name: "adding with different key adds",
			options: []Option{
				WithMockRoot("testdata/"),
			},
			substitutions: []*VariableSubstitution{
				{key: "foo", value: "bar"},
				{key: "bing", value: "baz"},
			},
			want: []Transformer{
				&VariableSubstitution{key: "foo", value: "bar"},
				&VariableSubstitution{key: "bing", value: "baz"},
			},
		},
		{
			name: "adding with same key overrides",
			options: []Option{
				WithMockRoot("testdata/"),
			},
			substitutions: []*VariableSubstitution{
				{key: "foo", value: "bar"},
				{key: "foo", value: "baz"},
			},
			want: []Transformer{
				&VariableSubstitution{key: "foo", value: "baz"},
			},
		},
	}

	for _, tc := range tcs {
		tc := tc // capture range variable
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ms, err := NewMockServer(tc.options...)
			require.Nil(t, err)

			for _, s := range tc.substitutions {
				ms.addVariableSubstitution(s)
			}

			got := ms.transformers
			assert.Equal(t, tc.want, got)
		})
	}
}

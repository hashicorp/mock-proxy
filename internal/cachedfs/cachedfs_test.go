package cachedfs

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPathExists(t *testing.T) {
	tcs := []struct {
		name         string
		path         string
		want         bool
		wantMapState map[string]bool
	}{
		{
			name: "simple",
			path: "testdata/testdir",
			want: true,
			wantMapState: map[string]bool{
				"testdata/testdir": true,
			},
		},
		{
			name:         "simple miss",
			path:         "testdata/notarealdirectorythough",
			want:         false,
			wantMapState: nil,
		},
	}

	for _, tc := range tcs {
		tc := tc // capture range variable
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			cf, err := NewCachedFS()
			require.Nil(t, err)

			got := cf.PathExists(tc.path)
			assert.Equal(t, tc.want, got)
			assert.Equal(t, tc.wantMapState, cf.hits)
		})
	}
}

func TestAddPath(t *testing.T) {
	tcs := []struct {
		name  string
		paths []string
		want  map[string]bool
	}{
		{
			name:  "add one",
			paths: []string{"a"},
			want: map[string]bool{
				"a": true,
			},
		},
		{
			name:  "add two",
			paths: []string{"a", "b"},
			want: map[string]bool{
				"a": true,
				"b": true,
			},
		},
		{
			name:  "add two repeat",
			paths: []string{"a", "a"},
			want: map[string]bool{
				"a": true,
			},
		},
	}

	for _, tc := range tcs {
		tc := tc // capture range variable
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			cf, err := NewCachedFS()
			require.Nil(t, err)

			for _, p := range tc.paths {
				cf.AddPath(p)
			}

			assert.Equal(t, tc.want, cf.hits)
		})
	}
}

func BenchmarkUncachedExists(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = os.Stat("testdata/testdir")
	}
}

func BenchmarkCachedExists(b *testing.B) {
	cf, _ := NewCachedFS()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cf.PathExists("testdata/testdir")
	}
}

func BenchmarkUncachedNotExists(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = os.Stat("testdata/notarealdirectory")
	}
}

func BenchmarkCachedNotExists(b *testing.B) {
	cf, _ := NewCachedFS()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cf.PathExists("testdata/notarealdirectory")
	}
}

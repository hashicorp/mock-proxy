package cachedfs

import (
	"os"
	"testing"
	"time"

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
			name: "simple miss",
			path: "testdata/notarealdirectorythough",
			want: false,
			wantMapState: map[string]bool{
				"testdata/notarealdirectorythough": false,
			},
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
			cf.m.RLock()
			assert.Equal(t, tc.wantMapState, cf.hits)
			cf.m.RUnlock()
		})
	}
}

func TestAddPath(t *testing.T) {
	tcs := []struct {
		name   string
		paths  []string
		values []bool
		want   map[string]bool
	}{
		{
			name:   "add one",
			paths:  []string{"a"},
			values: []bool{true},
			want: map[string]bool{
				"a": true,
			},
		},
		{
			name:   "add two",
			paths:  []string{"a", "b"},
			values: []bool{true, true},
			want: map[string]bool{
				"a": true,
				"b": true,
			},
		},
		{
			name:   "add two repeat",
			paths:  []string{"a", "a"},
			values: []bool{true, true},
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

			for i, p := range tc.paths {
				cf.AddPath(p, tc.values[i])
			}

			cf.m.RLock()
			assert.Equal(t, tc.want, cf.hits)
			cf.m.RUnlock()
		})
	}
}

func TestInvalidateFunc(t *testing.T) {
	cf, err := NewCachedFS(WithCacheExpiry(1 * time.Second))
	require.Nil(t, err)

	want := map[string]bool{"foo": true}
	cf.AddPath("foo", true)
	cf.m.RLock()
	assert.Equal(t, want, cf.hits)
	cf.m.RUnlock()

	want = map[string]bool{}
	time.Sleep(2 * time.Second)
	cf.m.RLock()
	assert.Equal(t, want, cf.hits)
	cf.m.RUnlock()
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

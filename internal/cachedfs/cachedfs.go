package cachedfs

import (
	"os"
	"sync"
	"time"
)

// invalidationFunc is a type alias for functions that can clear the cache.
type invalidationFunc func(*CachedFS)

// Option is a configuration option for passing to the CachedFS constructor.
// This is used to implement the "Functional Options" pattern:
//    https://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis
type Option func(*CachedFS) error

// CachedFS is a simple cache structure that allows you to cache whether a
// given directory exists.
type CachedFS struct {
	m                sync.RWMutex
	hits             map[string]bool
	invalidationFunc invalidationFunc
}

// NewCachedFS is the constructor for CachedFS.
func NewCachedFS(options ...Option) (*CachedFS, error) {
	cf := &CachedFS{
		hits: map[string]bool{},
	}

	for _, o := range options {
		if err := o(cf); err != nil {
			return nil, err
		}
	}

	if cf.invalidationFunc != nil {
		go cf.invalidationFunc(cf)
	}

	return cf, nil
}

// WithSimpleCacheExpiry is a functional option which will delete the entire
// cache after a set amount of time. This is very basic, but covers the easy
// case.
func WithSimpleCacheExpiry(ttl time.Duration) Option {
	return func(cf *CachedFS) error {
		cf.invalidationFunc = func(cf *CachedFS) {
			ticker := time.NewTicker(ttl)

			for range ticker.C {
				cf.m.Lock()
				cf.hits = map[string]bool{}
				cf.m.Unlock()
			}
		}
		return nil
	}
}

// PathExists checks if the given path exists and caches the lookup so it won't
// be repeated.
func (cf *CachedFS) PathExists(path string) bool {
	cf.m.RLock()

	val, hit := cf.hits[path]
	if !hit {
		// On miss, check if filepath exists.
		var exists bool
		if _, err := os.Stat(path); err == nil {
			exists = true
		}

		// cache information about whether the path exists.
		cf.m.RUnlock()
		cf.addPath(path, exists)
		return exists
	}

	cf.m.RUnlock()
	return val
}

// addPath adds a path to the cache.
func (cf *CachedFS) addPath(path string, val bool) {
	cf.m.Lock()
	defer cf.m.Unlock()

	if cf.hits == nil {
		cf.hits = map[string]bool{}
	}

	cf.hits[path] = val
}

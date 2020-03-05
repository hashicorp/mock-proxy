package cachedfs

import (
	"os"
	"sync"
	"time"
)

type invalidationFunc func(*CachedFS)

type Option func(*CachedFS) error

type CachedFS struct {
	m                sync.RWMutex
	hits             map[string]bool
	invalidationFunc invalidationFunc
}

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

func WithCacheExpiry(ttl time.Duration) Option {
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
		cf.AddPath(path, exists)
		return exists
	}

	cf.m.RUnlock()
	return val
}

func (cf *CachedFS) AddPath(path string, val bool) {
	cf.m.Lock()
	defer cf.m.Unlock()

	if cf.hits == nil {
		cf.hits = map[string]bool{}
	}

	cf.hits[path] = val
}

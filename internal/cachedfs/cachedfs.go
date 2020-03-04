package cachedfs

import (
	"os"
	"sync"
)

type CachedFS struct {
	m    sync.RWMutex
	hits map[string]bool
}

func NewCachedFS() (*CachedFS, error) {
	return &CachedFS{}, nil
}

func (cf *CachedFS) PathExists(path string) bool {
	cf.m.RLock()
	_, hit := cf.hits[path]
	if !hit {
		// On miss, check if filepath exists.
		if _, err := os.Stat(path); err == nil {
			// if it does, cache that, and return true
			cf.m.RUnlock()
			cf.AddPath(path)
			return true
		}
		// otherwise, return false, but cache nothing, we'll always test
		// misses against the filesystem.
		cf.m.RUnlock()
		return false
	}

	cf.m.RUnlock()
	return true
}

func (cf *CachedFS) AddPath(path string) {
	cf.m.Lock()
	defer cf.m.Unlock()

	if cf.hits == nil {
		cf.hits = map[string]bool{}
	}

	cf.hits[path] = true
}

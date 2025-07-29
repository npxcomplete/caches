package thread_safe

import (
	"sync"

	"github.com/npxcomplete/caches/src"
)

// NewGuardedCache wraps the provided cache with a simple mutex to make it safe
// for concurrent use.
func NewGuardedCache[K comparable, V any](p caches.Interface[K, V]) caches.Interface[K, V] {
	return &guardedLRU[K, V]{
		generic: p,
	}
}

// guardedLRU serialises access to the underlying cache using a mutex.
type guardedLRU[K comparable, V any] struct {
	mut     sync.Mutex
	generic caches.Interface[K, V]
}

// see caches.Interface for contract
func (cache *guardedLRU[K, V]) Keys() []K {
	cache.mut.Lock()
	defer cache.mut.Unlock()
	return cache.generic.Keys()
}

// see caches.Interface for contract
func (cache *guardedLRU[K, V]) Put(key K, value V) V {
	cache.mut.Lock()
	defer cache.mut.Unlock()
	return cache.generic.Put(key, value)
}

// see caches.Interface for contract
func (cache *guardedLRU[K, V]) Get(key K) (result V, err error) {
	cache.mut.Lock()
	defer cache.mut.Unlock()
	return cache.generic.Get(key)
}

// see caches.Interface for contract
func (cache *guardedLRU[K, V]) Range(f func(K, V) bool) {
	cache.mut.Lock()
	defer cache.mut.Unlock()
	cache.generic.Range(f)
}

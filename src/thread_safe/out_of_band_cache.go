package thread_safe

import (
	"context"

	"github.com/npxcomplete/caches/src"
)

// NewOutOfBandCache provides a cache implementation that performs all
// operations on a single goroutine. It is much slower than the guarded
// implementation but is included for completeness.
func NewOutOfBandCache[K comparable, V any](ctx context.Context, inner caches.Interface[K, V]) (safe caches.Interface[K, V]) {
	universalConstruction := make(chan func())

	safe = &outOfBan[K, V]{
		generic:               inner,
		universalConstruction: universalConstruction,
	}

	go func(ctx context.Context, inner caches.Interface[K, V]) {
		for {
			select {
			case f := <-universalConstruction:
				f()
			case <-ctx.Done():
				return
			}
		}
	}(ctx, inner)

	return
}

// outOfBan serialises all cache operations on a dedicated goroutine.
type outOfBan[K comparable, V any] struct {
	universalConstruction chan func()
	generic               caches.Interface[K, V]
}

func (cache *outOfBan[K, V]) Keys() []K {
	ret := make(chan []K)
	cache.universalConstruction <- func() {
		ret <- cache.generic.Keys()
	}
	return <-ret
}

// see caches.Interface for contract
func (cache *outOfBan[K, V]) Put(key K, value V) V {
	ret := make(chan V)
	cache.universalConstruction <- func() {
		ret <- cache.generic.Put(key, value)
	}
	return <-ret
}

// see caches.Interface for contract
func (cache *outOfBan[K, V]) Get(key K) (result V, err error) {
	ret := make(chan func() (V, error))
	cache.universalConstruction <- func() {
		val, err := cache.generic.Get(key)
		ret <- func() (V, error) { return val, err }
	}
	return (<-ret)()
}

// see caches.Interface for contract
func (cache *outOfBan[K, V]) Range(f func(K, V) bool) {
	done := make(chan struct{})
	cache.universalConstruction <- func() {
		cache.generic.Range(f)
		close(done)
	}
	<-done
}

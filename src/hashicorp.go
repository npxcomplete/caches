package caches

import hashicorp "github.com/hashicorp/golang-lru"

// hashi2Q is a thin generic wrapper around the hashicorp 2Q implementation
// which operates on interface{} types.
type hashi2Q[K comparable, V any] struct {
	inner *hashicorp.TwoQueueCache
}

func (h hashi2Q[K, V]) Keys() []K {
	raw := h.inner.Keys()
	keys := make([]K, len(raw))
	for i, k := range raw {
		keys[i] = k.(K)
	}
	return keys
}

func (h hashi2Q[K, V]) Put(k K, v V) V {
	h.inner.Add(k, v)
	var zero V
	return zero
}

func (h hashi2Q[K, V]) Get(k K) (get V, err error) {
	raw, ok := h.inner.Get(k)
	if !ok {
		err = MissingValueError
		var zero V
		return zero, err
	}
	get, _ = raw.(V)
	return
}

// Range invokes f for every key/value pair stored in the cache. Iteration
// stops if f returns false.
func (h hashi2Q[K, V]) Range(f func(K, V) bool) {
	for _, k := range h.inner.Keys() {
		v, ok := h.inner.Peek(k)
		if !ok {
			continue
		}
		if !f(k.(K), v.(V)) {
			return
		}
	}
}

func New2Q[K comparable, V any](capacity int) hashi2Q[K, V] {
	q, err := hashicorp.New2Q(capacity)
	if err != nil {
		panic("Illegal argument, capacity was negative.")
	}
	return hashi2Q[K, V]{inner: q}
}

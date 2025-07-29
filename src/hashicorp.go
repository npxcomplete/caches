package caches

import hashicorp "github.com/hashicorp/golang-lru"

type hashi2Q struct {
	inner *hashicorp.TwoQueueCache
}

func (h hashi2Q) Keys() []Key {
	raw := h.inner.Keys()
	keys := make([]Key, len(raw))
	for i, k := range raw {
		keys[i] = k
	}
	return keys
}

func (h hashi2Q) Put(k Key, v Value) Value {
	h.inner.Add(k, v)
	return nil
}

func (h hashi2Q) Get(k Key) (get Value, err error) {
	get, ok := h.inner.Get(k)
	if !ok {
		err = MissingValueError
	}
	return
}

// Range invokes f for every key/value pair stored in the cache. Iteration
// stops if f returns false.
func (h hashi2Q) Range(f func(Key, Value) bool) {
	for _, k := range h.inner.Keys() {
		v, ok := h.inner.Peek(k)
		if !ok {
			continue
		}
		if !f(k, v) {
			return
		}
	}
}

func New2Q(capacity int) hashi2Q {
	q, err := hashicorp.New2Q(capacity)
	if err != nil {
		panic("Illegal argument, capacity was negative.")
	}
	return hashi2Q{inner: q}
}

package caches

import hashicorp "github.com/hashicorp/golang-lru"

type hashi2Q struct {
	inner *hashicorp.TwoQueueCache
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

func New2Q(capacity int) hashi2Q {
	q, err := hashicorp.New2Q(capacity)
	if err != nil {
		panic("Illegal argument, capacity was negative.")
	}
	return hashi2Q{inner: q}
}

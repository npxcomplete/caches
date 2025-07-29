package caches

// Interface is the common contract implemented by all caches. K is the key
// type and must be comparable so that it can be used in maps. V represents the
// stored value type.
type Interface[K comparable, V any] interface {
	// Keys returns a slice of every key currently stored in the cache.
	Keys() []K

	// Put inserts the given key/value pair into the cache returning the
	// evicted value if any, else the zero value for V.
	Put(K, V) V

	// Get retrieves the corresponding value if present, else returns
	// `caches.MissingValueError`.
	Get(K) (V, error)

	// Range will invoke the provided function for each key/value pair in the
	// cache. Iteration stops if the function returns false.
	Range(func(K, V) bool)
}

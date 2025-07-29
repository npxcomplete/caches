package caches

type Key interface{}
type Value interface{}

type Interface interface {
	Keys() []Key

	// Insert the given key/value pair into the cache returning the evicted value if any, else return nil.
	Put(Key, Value) Value

	// Retrieve the corrosponding value if present, else return `caches.MissingValueError`
	Get(Key) (Value, error)

	// Range will invoke the provided function for each key/value pair in the
	// cache. Iteration stops if the function returns false.
	Range(func(Key, Value) bool)
}

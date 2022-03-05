package caches

type Key interface{}
type Value interface{}

type Interface interface{
	Keys() []Key

	// Insert the given key/value pair into the cache returning the evicted value if any, else return nil.
	Put(Key, Value) Value

	// Retrieve the corrosponding value if present, else return `caches.MissingValueError`
	Get(Key) (Value, error)
}

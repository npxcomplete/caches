package caches

// lruNode is a single entry in the cache.
// 48 bytes (16 per interface, 8 per concrete pointer) in the original
// implementation. With generics the size is similar and depends on the key and
// value types.
type lruNode[K comparable, V any] struct {
	next *lruNode[K, V]
	prev *lruNode[K, V]

	key   K
	value V
}

type lruCache[K comparable, V any] struct {
	// add to head
	head *lruNode[K, V]

	// fast access
	store map[K]*lruNode[K, V]

	// full if stack == nil
	stack *lruNode[K, V]

	// v_node pool, to prevent GC churn
	nodes []lruNode[K, V]
}

// NewLRUCache returns a new LRU cache with the provided capacity.
// ~48 bytes per entry in the original implementation.
func NewLRUCache[K comparable, V any](capacity int) *lruCache[K, V] {
	memoryPool := make([]lruNode[K, V], capacity)
	for i := 0; i < capacity-1; i++ {
		memoryPool[i].next = &memoryPool[i+1]
	}

	// simplified nil checks
	dummy := &lruNode[K, V]{}

	dummy.next = dummy
	dummy.prev = dummy

	return &lruCache[K, V]{
		store: make(map[K]*lruNode[K, V], capacity),
		stack: &memoryPool[0],
		head:  dummy,
		nodes: memoryPool,
	}
}

func (lru *lruCache[K, V]) Keys() []K {
	keys := make([]K, 0, len(lru.nodes))
	for _, node := range lru.nodes {
		keys = append(keys, node.key)
	}
	return keys
}

func (lru *lruCache[K, V]) Put(key K, value V) (evictedValue V) {
	var node *lruNode[K, V]
	var ok bool

	if node, ok = lru.store[key]; ok {
		// key already present, evict present lruNode to replace
	} else if lru.stack == nil {
		// cache full, evict the tail
		node = lru.head.prev
	}

	// do eviction
	if node != nil {
		node.prev.next = node.next
		node.next.prev = node.prev

		node.prev = nil
		node.next = nil

		// stack push
		node.next = lru.stack
		lru.stack = node

		evictedValue = node.value
		delete(lru.store, node.key)
	}

	// stack must be non-empty by now
	node = lru.stack
	lru.stack = lru.stack.next

	node.key = key
	node.value = value

	// insert at head
	node.next = lru.head.next
	node.next.prev = node
	node.prev = lru.head
	lru.head.next = node

	lru.store[key] = node

	return
}

func (lru *lruCache[K, V]) Get(key K) (value V, err error) {
	var node *lruNode[K, V]
	var ok bool
	if node, ok = lru.store[key]; !ok {
		err = MissingValueError
		var zero V
		return zero, err
	}

	value = node.value
	// Put implicitly resets eviction priority
	lru.Put(key, value)
	return
}

// Range iterates over each entry in the cache. Iteration stops if f returns
// false.

func (lru *lruCache[K, V]) Range(f func(K, V) bool) {
	for k, node := range lru.store {
		if !f(k, node.value) {
			return
		}
	}
}

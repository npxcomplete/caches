package thread_safe

import (
	"sync"

	"github.com/npxcomplete/caches/src"
)

// 48 bytes (16 per interface, 8 per concrete pointer)
type lruNode struct {
	mutex sync.Mutex

	next *lruNode
	prev *lruNode

	key   caches.Key
	value caches.Value
}

type lruCache struct {
	// MRU is found at head.next
	// LRU is found at head.prev
	head *lruNode

	// provide access for bounded constant time search & removal
	store sync.Map

	// full if stack == nil
	stack      *lruNode
	stackGuard sync.Mutex

	// v_node pool, to prevent GC churn
	nodes []lruNode
}

// ~ 48 bytes per entry.
func New_Broken_LRUCache(capacity int) *lruCache {
	memoryPool := make([]lruNode, capacity)
	for i := 0; i < capacity-1; i++ {
		memoryPool[i].next = &memoryPool[i+1]
	}

	// simplified nil checks
	dummy := &lruNode{}

	dummy.next = dummy
	dummy.prev = dummy

	return &lruCache{
		store: sync.Map{},
		stack: &memoryPool[0],
		head:  dummy,
		nodes: memoryPool,
	}
}

func (lru *lruCache) Put(key caches.Key, value caches.Value) (evictedValue caches.Value) {
	var node *lruNode = nil
	var ok bool

	untyped, ok := lru.store.Load(key) // here
	if ok && untyped != nil {
		node = untyped.(*lruNode)
	} else if lru.stack == nil {
		// cache full, evict the tail
		node = lru.head.prev
	} // else // key already present, evict present lruNode to replace

	// do eviction
	if node != nil {
		node.prev.next = node.next // here
		node.next.prev = node.prev

		node.prev = nil
		node.next = nil

		// stack push
		node.next = lru.stack
		lru.stack = node

		evictedValue = node.value
		lru.store.Delete(node.key)
	}

	// stack must be non-empty by now
	node = lru.stack
	lru.stack = lru.stack.next

	node.key = key
	node.value = value

	// insert at head
	node.next = lru.head.next
	node.next.prev = node // her
	node.prev = lru.head
	lru.head.next = node

	lru.store.Store(key, node)

	return
}

func (lru *lruCache) Get(key caches.Key) (value caches.Value, err error) {
	var node *lruNode
	var ok bool
	untyped, ok := lru.store.Load(key)
	if !ok {
		err = caches.MissingValueError
		return
	}
	node = untyped.(*lruNode)
	value = node.value
	// Put implicitly resets eviction priority
	lru.Put(key, value)
	return
}

func (lru *lruCache) evictOrAlloc() (node *lruNode) {
	func() {
		lru.stackGuard.Lock()
		defer lru.stackGuard.Unlock()

		if lru.stack != nil {
			node = lru.stack
			lru.stack = lru.stack.next
		}
	}()

	if node != nil {
		// memory available, eviction not required
		return
	}

	//remove node from the center, retaining it's lock
	lru.head.mutex.Lock()
	node = lru.head.prev
	node.mutex.Lock()
	node.prev.mutex.Lock()

	lru.head.prev = node.prev
	node.prev.next = lru.head

	lru.head.mutex.Unlock()
	node.prev.mutex.Unlock()
	////////////////////////////////////

	node.prev = nil
	node.next = nil
	lru.store.Delete(node.key)
	node.mutex.Unlock()
	return
}

func (lru *lruCache) pushNode(node *lruNode) {
	// node should be unused and therefore should not be
	// pointed at by any thread, ergo it should not be under contention
	node.mutex.Lock()

	lru.head.mutex.Lock()
	oldMRU := lru.head.next
	oldMRU.mutex.Lock()

	lru.head.next = node
	lru.head.mutex.Unlock()

	node.prev = lru.head
	node.next = oldMRU
	oldMRU.prev = node
	oldMRU.mutex.Unlock()

	lru.store.Store(node.key, node)
	node.mutex.Unlock()
	return
}

func (lru *lruCache) removeFromCenter(key caches.Key) (node *lruNode) {
	return
}

package caches_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	caches "github.com/npxcomplete/caches/src"
	thread_safe "github.com/npxcomplete/caches/src/thread_safe"
)

func Test_AddToEmptyCache(t *testing.T) {
	var kvLRU caches.Interface = caches.NewLRUCache(2)

	kvLRU.Put(5, "hello")

	v, err := kvLRU.Get(5)
	assert.Equal(t, "hello", v)
	assert.Nil(t, err)

	v, err = kvLRU.Get(6)
	assert.NotNil(t, err)
}

func Test_AddToFullCache(t *testing.T) {
	var kvLRU caches.Interface = caches.NewLRUCache(2)

	kvLRU.Put(5, "hello")
	kvLRU.Put(6, "world")
	kvLRU.Put(7, "!")

	var value interface{}
	var err error

	value, err = kvLRU.Get(6)
	assert.Equal(t, "world", value)
	assert.Nil(t, err)

	value, err = kvLRU.Get(7)
	assert.Equal(t, "!", value)
	assert.Nil(t, err)

	value, err = kvLRU.Get(5)
	assert.NotNil(t, err, "earlier value must be evicted")
}

func Test_QueriedValuesReset(t *testing.T) {
	var kvLRU caches.Interface = caches.NewLRUCache(2)

	kvLRU.Put(5, "hello")
	kvLRU.Put(6, "world")
	kvLRU.Get(5)

	kvLRU.Put(7, "!")

	var value interface{}
	var err error

	value, err = kvLRU.Get(5)
	assert.Equal(t, "hello", value)
	assert.Nil(t, err)

	value, err = kvLRU.Get(7)
	assert.Equal(t, "!", value)
	assert.Nil(t, err)

	value, err = kvLRU.Get(6)
	assert.NotNil(t, err)
}

func Test_ValueNeverPresent(t *testing.T) {
	var kvLRU caches.Interface = caches.NewLRUCache(2)

	kvLRU.Put(5, "hello")
	kvLRU.Put(6, "world")

	var err error

	_, err = kvLRU.Get(12)
	assert.NotNil(t, err)
}

func Test_RangeVisitsAll(t *testing.T) {
	var kvLRU caches.Interface = caches.NewLRUCache(2)

	kvLRU.Put(1, "a")
	kvLRU.Put(2, "b")

	visited := make(map[int]string)
	kvLRU.Range(func(k caches.Key, v caches.Value) bool {
		visited[k.(int)] = v.(string)
		return true
	})

	assert.Equal(t, 2, len(visited))
	assert.Equal(t, "a", visited[1])
	assert.Equal(t, "b", visited[2])
}

func Test_RangeEarlyExit(t *testing.T) {
	var kvLRU caches.Interface = caches.NewLRUCache(3)
	kvLRU.Put(1, "a")
	kvLRU.Put(2, "b")
	kvLRU.Put(3, "c")

	count := 0
	kvLRU.Range(func(k caches.Key, v caches.Value) bool {
		count++
		return false
	})

	assert.Equal(t, 1, count)
}

func Test_RangeThreadSafeWrappers(t *testing.T) {
	cache := caches.NewLRUCache(2)
	cache.Put(1, "a")
	cache.Put(2, "b")

	guard := thread_safe.NewGuardedCache(cache)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	band := thread_safe.NewOutOfBandCache(ctx, cache)

	countGuard := 0
	guard.Range(func(k caches.Key, v caches.Value) bool {
		countGuard++
		return true
	})

	countBand := 0
	band.Range(func(k caches.Key, v caches.Value) bool {
		countBand++
		return true
	})

	assert.Equal(t, 2, countGuard)
	assert.Equal(t, 2, countBand)
}

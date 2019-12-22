package thread_safe

import (
	"context"
	"math/rand"
	"sync"
	"testing"

	caches "github.com/npxcomplete/caches/src"
)

var capacity = 10
var threads = 4

func Benchmark_GuardedLRU_Get(b *testing.B) {
	cache := NewGuardedCache(caches.NewLRUCache(capacity))
	benchmark_Get(b, cache)
}

func Benchmark_OutOfBan_Get(b *testing.B) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	cache := NewOutOfBandCache(ctx, caches.NewLRUCache(capacity))
	benchmark_Get(b, cache)
}

func Benchmark_GuardedLRU_Put(b *testing.B) {
	cache := NewGuardedCache(caches.NewLRUCache(capacity))
	benchmark_Put(b, cache)
}

func Benchmark_OutOfBan_Put(b *testing.B) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	cache := NewOutOfBandCache(ctx, caches.NewLRUCache(capacity))
	benchmark_Put(b, cache)
}

func Benchmark_HashiCorp_GuardedLRU_parallel_Get(b *testing.B) {
	cache := caches.New2Q(capacity)
	benchmark_parallel_Get(b, cache)
}

func Benchmark_HashiCorp_LRU_parallel_Put(b *testing.B) {
	cache := caches.New2Q(capacity)
	benchmark_parallel_Put(b, cache)
}

func Benchmark_GuardedLRU_parallel_Get(b *testing.B) {
	cache := NewGuardedCache(caches.NewLRUCache(capacity))
	benchmark_parallel_Get(b, cache)
}

func Benchmark_GuardedLRU_parallel_Put(b *testing.B) {
	cache := NewGuardedCache(caches.NewLRUCache(capacity))
	benchmark_parallel_Put(b, cache)
}

func Benchmark_OutOfBan_parallel_Get(b *testing.B) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	cache := NewOutOfBandCache(ctx, caches.NewLRUCache(capacity))
	benchmark_parallel_Get(b, cache)
}

func Benchmark_OutOfBan_parallel_Put(b *testing.B) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	cache := NewOutOfBandCache(ctx, caches.NewLRUCache(capacity))
	benchmark_parallel_Put(b, cache)
}

func benchmark_Get(b *testing.B, cache caches.Interface) {
	keys := genKeys()
	for i := 0; i < len(keys); i++ {
		cache.Put(keys[i], 12)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.Get(keys[i%len(keys)])
	}
}

func benchmark_Put(b *testing.B, cache caches.Interface) {
	keys := genKeys()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.Put(keys[i%len(keys)], 91)
	}
}


func benchmark_parallel_Get(b *testing.B, cache caches.Interface) {
	keys := genKeys()
	for i := 0; i < len(keys); i++ {
		cache.Put(keys[i], 12)
	}

	b.ResetTimer()
	var wg sync.WaitGroup
	for t := 0; t < threads; t++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < b.N/threads; i++ {
				cache.Get(keys[rand.Int31n(int32(len(keys)))])
			}
		}()
	}
	wg.Wait()
}

func benchmark_parallel_Put(b *testing.B, cache caches.Interface) {
	keys := genKeys()

	b.ResetTimer()
	var wg sync.WaitGroup
	for t := 0; t < threads; t++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < b.N/threads; i++ {
				cache.Put(keys[rand.Int31n(int32(len(keys)))], 91)
			}
		}()
	}
	wg.Wait()
}

func genKeys() []int {
	alloc := make([]int, capacity*2)
	for i := 0; i < len(alloc); i++ {
		alloc[i] = rand.Int()
	}
	return alloc
}


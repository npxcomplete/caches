package thread_safe

import (
	"math/rand"
	"testing"

	"github.com/npxcomplete/random/src/strings"
)

func Benchmark_LRUCachewMapInterface_Put_single_key(b *testing.B) {
	lru := NewLRUCache(2)

	for i := 0; i < b.N; i++ {
		lru.Put("hello", "world")
	}
}

func Benchmark_LRUCachewMapInterface_Put_multi_key_with_high_eviction(b *testing.B) {
	lru := NewLRUCache(2)

	gen := random_strings.ByteStringGenerator{
		Alphabet:  random_strings.EnglishAlphabet,
		RandomGen: rand.New(rand.NewSource(0)),
	}

	for i := 0; i < b.N; i++ {
		lru.Put(gen.String(5), "world")
	}
}

func Benchmark_LRUCachewMapInterface_Put_multi_key_with_low_eviction(b *testing.B) {
	lru := NewLRUCache(10000)

	gen := random_strings.ByteStringGenerator{
		Alphabet:  random_strings.EnglishAlphabet,
		RandomGen: rand.New(rand.NewSource(0)),
	}

	for i := 0; i < b.N; i++ {
		lru.Put(gen.String(4), "world")
	}
}

func Benchmark_LRUCachewMapInterface_Get_MRU(b *testing.B) {
	lru := NewLRUCache(10)

	keys := []string {
		"hello1",
		"hello2",
		"hello3",
		"hello4",
		"hello5",
	}

	for _, key := range keys {
		lru.Put(key, "world")
	}

	for i := 0; i < b.N; i++ {
		lru.Get(keys[2])
	}
}

func Benchmark_LRUCachewMapInterface_Get_cycle(b *testing.B) {
	lru := NewLRUCache(10)

	keys := []string {
		"hello1",
		"hello2",
		"hello3",
		"hello4",
		"hello5",
	}

	for _, key := range keys {
		lru.Put(key, "world")
	}

	for i := 0; i < b.N; i++ {
		lru.Get(keys[i % len(keys)])
	}
}

func Benchmark_LRUCachewMapInterface_Put_with_eviction_and_no_gen(b *testing.B) {
	lru := NewLRUCache(4)

	keys := []string{
		"aaaa",
		"bbbb",
		"cccc",
		"dddd",
		"eeee",
	}

	for i := 0; i < b.N; i++ {
		lru.Put(keys[i % len(keys)], "world")
	}
}
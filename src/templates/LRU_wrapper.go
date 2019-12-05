package templates

import (
	"github.com/cheekybits/genny/generic"

	"github.com/npxcomplete/caches/src"
)

type GenericKey generic.Type
type GenericValue generic.Type

func NewGenericKeyGenericValueLRUCache(capacity int) privateGenericKeyGenericValueLRUCache {
	return privateGenericKeyGenericValueLRUCache{caches.NewLRUCache(capacity)}
}

// genny is case sensitive even though this has other meanings in go, so we prefix the intent.
type privateGenericKeyGenericValueLRUCache struct {
	generic caches.Interface
}

// see caches.Interface for contract
func (cache privateGenericKeyGenericValueLRUCache) Put(key GenericKey, value GenericValue) GenericValue {
	result, _ := cache.generic.Put(key, value).(GenericValue)
	return result
}

// see caches.Interface for contract
func (cache privateGenericKeyGenericValueLRUCache) Get(key GenericKey) (result GenericValue, err error) {
	value, err := cache.generic.Get(key)
	result, _ = value.(GenericValue)
	return
}


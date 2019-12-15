// This file was automatically generated by genny.
// Any changes will be lost if this file is regenerated.
// see https://github.com/cheekybits/genny

package concrete_caches

import (
	"github.com/npxcomplete/caches/examples"
	caches "github.com/npxcomplete/caches/src"
)

func NewExamplesMyCustomTypeExamplesMyCustomTypeLRUCache(capacity int) privateExamplesMyCustomTypeExamplesMyCustomTypeLRUCache {
	return privateExamplesMyCustomTypeExamplesMyCustomTypeLRUCache{caches.NewLRUCache(capacity)}
}

// genny is case sensitive even though this has other meanings in go, so we prefix the intent.
type privateExamplesMyCustomTypeExamplesMyCustomTypeLRUCache struct {
	generic caches.Interface
}

// see caches.Interface for contract
func (cache privateExamplesMyCustomTypeExamplesMyCustomTypeLRUCache) Put(key examples.MyCustomType, value examples.MyCustomType) examples.MyCustomType {
	result, _ := cache.generic.Put(key, value).(examples.MyCustomType)
	return result
}

// see caches.Interface for contract
func (cache privateExamplesMyCustomTypeExamplesMyCustomTypeLRUCache) Get(key examples.MyCustomType) (result examples.MyCustomType, err error) {
	value, err := cache.generic.Get(key)
	result, _ = value.(examples.MyCustomType)
	return
}

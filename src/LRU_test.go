package src

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func Test_AddToEmptyCache(t *testing.T) {
	var kvLRU KeyTypeValueTypeCache = NewKeyTypeValueTypeLRU(2)

	kvLRU.Put(5, "hello")

	v, err := kvLRU.Get(5)
	assert.Equal(t, "hello", v)
	assert.Nil(t, err)

	v, err = kvLRU.Get(6)
	assert.NotNil(t, err)
}

func Test_AddToFullCache(t *testing.T) {
	var kvLRU KeyTypeValueTypeCache = NewKeyTypeValueTypeLRU(2)

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
	var kvLRU KeyTypeValueTypeCache = NewKeyTypeValueTypeLRU(2)

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
	var kvLRU KeyTypeValueTypeCache = NewKeyTypeValueTypeLRU(2)

	kvLRU.Put(5, "hello")
	kvLRU.Put(6, "world")

	var err error

	_, err = kvLRU.Get(12)
	assert.NotNil(t, err)
}

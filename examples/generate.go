// specify this one once
//go:generate genny -pkg=primitives -in=../src/single.go   -out=./primitives/constants.go             gen "KeyType= ValueType="

// specify this one once for each key / value pair
//go:generate genny -pkg=primitives -in=../src/common.go   -out=./primitives/int_string_type.go       gen "KeyType=int ValueType=string"

// specify this one once for each key / value / implementation triple
//go:generate genny -pkg=primitives -in=../src/LRU.go      -out=./primitives/int_string_cache.go      gen "KeyType=int ValueType=string"

package examples

// specify this one once for each key / value / implementation triple
//go:generate genny -pkg=concrete_caches -in=../src/templates/LRU_wrapper.go -out=./generated/int_string_cache.go  gen "GenericKey=int GenericValue=string"
//go:generate genny -pkg=concrete_caches -in=../src/templates/LRU_wrapper.go -out=./generated/custom_custom_cache.go  gen "GenericKey=examples.MyCustomType GenericValue=examples.MyCustomType"

package examples

type MyCustomType struct {}
package examples

import (
	"fmt"
	caches "github.com/npxcomplete/caches/src"
	"os"
)

func main() {
	cache := caches.NewLRUCache[int, []string](10)
	cache.Put(0, []string{"Hello, world!"})
	value, err := cache.Get(0)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	fmt.Println(value[0])
}

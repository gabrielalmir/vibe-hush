package main

import (
	"fmt"
	"time"

	"github.com/gabrielalmir/vibe-hush/cache"
)

func main() {
	cache := cache.NewCache(3, 5*time.Second, cache.LRU{})
	cache.Set("a", 1)
	cache.Set("b", 2)
	cache.Set("c", 3)
	fmt.Println("Initial Cache:", cache.Get("a"), cache.Get("b"), cache.Get("c"))
	cache.Set("d", 4)
	fmt.Println("After Eviction:", cache.Get("a"), cache.Get("b"), cache.Get("c"), cache.Get("d"))
}

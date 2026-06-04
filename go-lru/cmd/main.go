package main

import (
	"fmt"

	"github.com/arunhiremath92/go-lru/cache/lru"
)

func main() {
	lru, err := lru.NewCacheInstance(10, nil)
	if err != nil {
		panic(err)
	}
	lru.Put("string1", "10")
	lru.Put("string2", "20")

	lru.Dump()
	fmt.Println(lru.Get("string1"))
	fmt.Println(lru.Get("string2"))
	lru.Put("string3", "30")
	lru.Dump()
	fmt.Println(lru.Get("string3"))
	fmt.Println(lru.Get("string1"))
	lru.Dump()
}

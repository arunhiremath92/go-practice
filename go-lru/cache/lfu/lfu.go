package lfu

import (
	"container/heap"

	"fmt"

	"github.com/arunhiremath92/go-lru/cache"
)

type lfunode struct {
	count int
	value string
	key   string
	index int
}

type lfuheap []lfunode

type lfucache struct {
	size        int
	removedFunc cache.RemovedFunc
	cache       lfuheap
	freqMap     map[string]lfunode
}

func (h lfuheap) Len() int {
	return len(h)
}

func (h lfuheap) Less(i, j int) bool {
	return h[i].count < h[j].count
}

func (h lfuheap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
	h[i].index = i
	h[j].index = j
}

func (h *lfuheap) Push(x any) {
	// Push and Pop use pointer receivers because they modify the slice's length,
	// not just its contents.
	node := x.(lfunode)
	*h = append(*h, node)
}

func (h *lfuheap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func (cache *lfucache) Get(key string) interface{} {
	if _, exists := cache.freqMap[key]; !exists {
		fmt.Println("No Value Found in the Cache")
		return nil
	}

	node := cache.freqMap[key]
	cache.cache[node.index].count = cache.cache[node.index].count + 1
	heap.Fix(&cache.cache, node.index)
	return node
}

func (cache *lfucache) Put(key string, value string) interface{} {
	if node, exists := cache.freqMap[key]; exists {
		// updating the count of a element which was used again
		cache.cache[node.index].count = cache.cache[node.index].count + 1
		heap.Fix(&cache.cache, node.index)
		return key

	}

	// check if the cache limit has reached, if yes pop the least used
	if cache.cache.Len() >= cache.size {
		removedNode := heap.Pop(&cache.cache).(lfunode)
		fmt.Println("removed: ", removedNode.value, removedNode.key, " repeated: ", removedNode.count)
		delete(cache.freqMap, removedNode.key)
	}

	// make grounds for inserting new one
	lfunode := lfunode{
		count: 1,
		key:   key,
		value: value,
	}
	heap.Push(&cache.cache, lfunode)
	cache.freqMap[key] = lfunode
	return key
}

func (cache *lfucache) Delete(key string) interface{} { return string("") }

func (cache *lfucache) Dump() {
	cacheContents := []string{}
	for _, node := range cache.cache {
		cacheContents = append(cacheContents, fmt.Sprintf("%s:%s", node.key, node.value))
	}
	s := fmt.Sprintf("%s", cacheContents)
	fmt.Println(s)
}

func NewCacheInstance(cacheSize int, opts *cache.Options) (cache.Cache, error) {

	if cacheSize < 0 {
		return nil, fmt.Errorf("%s", "cache size requested is invalid")
	}
	if opts == nil {
		opts = &cache.Options{}
	}
	lfu := lfucache{
		size:        cacheSize,
		cache:       make(lfuheap, 0, opts.InitialCapacity),
		freqMap:     make(map[string]lfunode),
		removedFunc: opts.RemovedFunc,
	}
	heap.Init(&lfu.cache)
	return &lfu, nil
}

func (cache *lfucache) Size() int {
	return cache.size
}

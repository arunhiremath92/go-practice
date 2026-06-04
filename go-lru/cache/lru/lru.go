package lru

import (
	"container/list"
	"fmt"
	"sync"

	"github.com/arunhiremath92/go-lru/cache"
)

type Node struct {
	key   string
	value string
}

type lruCache struct {
	mutex     sync.Mutex
	cacheSize int
	dlist     *list.List
	nodeTable map[string]*list.Element
	rmFunc    cache.RemovedFunc
}

func NewCacheInstance(cacheSize int, opts *cache.Options) (cache.Cache, error) {
	if cacheSize < 0 {
		return nil, fmt.Errorf("%s", "cache size requested is invalid")
	}
	if opts == nil {
		opts = &cache.Options{}
	}

	return &lruCache{
		mutex:     sync.Mutex{},
		cacheSize: cacheSize,
		dlist:     list.New(),
		nodeTable: make(map[string]*list.Element),
		rmFunc:    opts.RemovedFunc,
	}, nil
}

func (cache *lruCache) Get(key string) interface{} {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()
	if element := cache.nodeTable[key]; element != nil {
		cache.dlist.MoveToFront(element)
		node := element.Value.(*Node)
		return node.value
	}
	return nil
}

func (cache *lruCache) Size() int {
	return cache.cacheSize
}

func (cache *lruCache) Dump() {
	cacheContents := []string{}
	for e := cache.dlist.Front(); e != nil; e = e.Next() {
		node := e.Value.(*Node)
		cacheContents = append(cacheContents, fmt.Sprintf("%s:%s", node.key, node.value))
	}
	s := fmt.Sprintf("%s", cacheContents)
	fmt.Println(s)
}

func (cache *lruCache) Put(key string, value string) interface{} {

	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	if element := cache.nodeTable[key]; element != nil {
		node := element.Value.(*Node)
		node.value = value
		cache.dlist.MoveToFront(element)
		return element.Value
	}
	if cache.dlist.Len() >= cache.cacheSize {
		back := cache.dlist.Back()
		node := back.Value.(*Node)
		delete(cache.nodeTable, node.key)
		cache.dlist.Remove(back)
	}
	node := &Node{key: key, value: value}
	elementLoc := cache.dlist.PushFront(node)
	cache.nodeTable[key] = elementLoc
	return elementLoc.Value
}

func (cache *lruCache) Delete(key string) interface{} {
	return nil
}

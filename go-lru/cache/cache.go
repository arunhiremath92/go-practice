package cache

const (
	DefaultCapcity = 100
)

// RemovedFunc is a type for notifying applications when an item is
// scheduled for removal from the Cache. If f is a function with the
// appropriate signature and i is the interface{} scheduled for
// deletion, Cache calls go f(i)
type RemovedFunc func(interface{})
type Options struct {
	// InitialCapacity controls the initial capacity of the cache
	InitialCapacity int

	// RemovedFunc is an optional function called when an element
	// is scheduled for deletion
	RemovedFunc RemovedFunc
}

type Cache interface {
	Get(key string) interface{}

	Put(key string, value string) interface{}

	Delete(key string) interface{}

	Dump()

	Size() int
}

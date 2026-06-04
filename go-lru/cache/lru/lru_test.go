package lru

import (
	"testing"
)

func TestLruInstance(t *testing.T) {
	lruInstance, err := NewCacheInstance(-1, nil)
	if err == nil || lruInstance != nil {
		t.Errorf("lru instance should be empty")
	}

	lru, _ := NewCacheInstance(10, nil)
	if lru == nil {
		t.Errorf("lru instance should not be nil")
	}
}

func TestLruFunctionality(t *testing.T) {
	lru, err := NewCacheInstance(2, nil)
	if err != nil {
		t.Errorf("lru instance should not be nil")
	}

	if lru.Get("non_existant_key") != nil {
		t.Errorf("lru Get on empty list should be nill")
	}

	if lru.Size() != 2 {
		t.Errorf("size is not the one created with want %d get %d", 2, lru.Size())
	}

	lru.Put("hello", "world")
	if lru.Get("hello").(string) != "world" {
		t.Errorf("put and get are not the same")
	}

	lru.Put("hello1", "world1")
	if lru.Get("hello").(string) != "world" {
		t.Errorf("put and get are not the same")
	}

	lru.Put("hello2", "world2")
	if lru.Get("hello").(string) != "world" {
		t.Errorf("put and get are not the same")
	}

	// since the capacity is full, hello would have been eliminated,
	// a get operation should fail
	if lru.Get("hello") != nil {
		t.Errorf("get should be empty")
	}

}

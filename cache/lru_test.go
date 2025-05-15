package cache

import (
	"sync"
	"testing"
)

func TestLRUCacheBasic(t *testing.T) {
	// Test non-thread-safe version
	cache := NewLRUCache[string, int](3, false)
	testBasicOperations(t, cache)

	// Test thread-safe version
	cache = NewLRUCache[string, int](3, true)
	testBasicOperations(t, cache)
}

func testBasicOperations(t *testing.T, cache *LRUCache[string, int]) {
	// Test empty cache
	if cache.Len() != 0 {
		t.Errorf("Expected length 0, got %d", cache.Len())
	}

	// Test Put and Get
	cache.Put("one", 1)
	cache.Put("two", 2)
	cache.Put("three", 3)

	// Access "one" to make it most recently used
	if val, exists := cache.Get("one"); !exists || val != 1 {
		t.Errorf("Expected value 1, got %v, exists: %v", val, exists)
	}

	// Test capacity limit - "two" should be evicted as it's now least recently used
	cache.Put("four", 4)
	if _, exists := cache.Get("two"); exists {
		t.Error("Expected 'two' to be evicted")
	}

	// Test Remove
	cache.Remove("one")
	if _, exists := cache.Get("one"); exists {
		t.Error("Expected 'one' to be removed")
	}
	if cache.Len() != 2 {
		t.Errorf("Expected length 2, got %d", cache.Len())
	}

	// Test Clear
	cache.Clear()
	if cache.Len() != 0 {
		t.Errorf("Expected length 0, got %d", cache.Len())
	}
}

func TestLRUCacheConcurrent(t *testing.T) {
	cache := NewLRUCache[string, int](100, true)
	var wg sync.WaitGroup
	iterations := 1000

	// Test concurrent Put
	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 0; i < iterations; i++ {
			cache.Put("key1", i)
		}
	}()
	go func() {
		defer wg.Done()
		for i := 0; i < iterations; i++ {
			cache.Put("key2", i)
		}
	}()
	wg.Wait()

	// Test concurrent Get
	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 0; i < iterations; i++ {
			cache.Get("key1")
		}
	}()
	go func() {
		defer wg.Done()
		for i := 0; i < iterations; i++ {
			cache.Get("key2")
		}
	}()
	wg.Wait()

	// Test concurrent Remove
	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 0; i < iterations; i++ {
			cache.Remove("key1")
		}
	}()
	go func() {
		defer wg.Done()
		for i := 0; i < iterations; i++ {
			cache.Remove("key2")
		}
	}()
	wg.Wait()
}

func TestLRUCacheEviction(t *testing.T) {
	cache := NewLRUCache[string, int](3, false)

	// Fill the cache
	cache.Put("one", 1)
	cache.Put("two", 2)
	cache.Put("three", 3)

	// Access "one" to make it most recently used
	cache.Get("one")

	// Add a new value, should evict "two" as it's now least recently used
	cache.Put("four", 4)

	// Verify "two" was evicted
	if _, exists := cache.Get("two"); exists {
		t.Error("Expected 'two' to be evicted")
	}

	// Verify other values are still present
	if val, exists := cache.Get("one"); !exists || val != 1 {
		t.Error("Expected 'one' to be present")
	}
	if val, exists := cache.Get("three"); !exists || val != 3 {
		t.Error("Expected 'three' to be present")
	}
	if val, exists := cache.Get("four"); !exists || val != 4 {
		t.Error("Expected 'four' to be present")
	}
}

func TestLRUCacheUpdate(t *testing.T) {
	cache := NewLRUCache[string, int](3, false)

	// Add initial value
	cache.Put("one", 1)

	// Update value
	cache.Put("one", 2)

	// Verify update
	if val, exists := cache.Get("one"); !exists || val != 2 {
		t.Errorf("Expected value 2, got %v, exists: %v", val, exists)
	}

	// Verify length hasn't changed
	if cache.Len() != 1 {
		t.Errorf("Expected length 1, got %d", cache.Len())
	}
}

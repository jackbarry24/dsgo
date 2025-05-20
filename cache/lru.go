package cache

import (
	"sync"

	"dsgo/linkedlist"
)

type LRUCache[K comparable, V any] struct {
	capacity   int
	cache      map[K]*linkedlist.DNode[K]
	list       *linkedlist.DoubleLinkedList[K]
	values     map[K]V
	threadSafe bool
	mu         sync.RWMutex
}

// NewLRUCache creates a new LRU cache with the specified capacity
func NewLRUCache[K comparable, V any](capacity int, threadSafe ...bool) *LRUCache[K, V] {
	isThreadSafe := true
	if len(threadSafe) > 0 {
		isThreadSafe = threadSafe[0]
	}
	return &LRUCache[K, V]{
		capacity:   capacity,
		cache:      make(map[K]*linkedlist.DNode[K]),
		list:       linkedlist.NewDoubleLinkedList[K](isThreadSafe),
		values:     make(map[K]V),
		threadSafe: isThreadSafe,
	}
}

// Get retrieves a value from the cache and marks it as most recently used
func (c *LRUCache[K, V]) Get(key K) (V, bool) {
	if c.threadSafe {
		c.mu.Lock()
		defer c.mu.Unlock()
	}

	if _, exists := c.cache[key]; exists {
		// Remove the node from its current position
		c.list.Remove(key)
		// Add it to the front (most recently used)
		c.list.PushFront(key)
		// Update the cache map with the new node
		if front, err := c.list.Front(); err == nil {
			c.cache[key] = front
		}
		return c.values[key], true
	}
	var zero V
	return zero, false
}

// Put adds or updates a value in the cache
func (c *LRUCache[K, V]) Put(key K, value V) {
	if c.threadSafe {
		c.mu.Lock()
		defer c.mu.Unlock()
	}

	// If key exists, update it
	if _, exists := c.cache[key]; exists {
		c.list.Remove(key)
	} else if c.list.Len() >= c.capacity {
		// If cache is full, remove the least recently used item
		if tail, err := c.list.Back(); err == nil {
			oldKey := tail.GetValue()
			c.list.Remove(oldKey)
			delete(c.cache, oldKey)
			delete(c.values, oldKey)
		}
	}

	// Add the new key to the front
	c.list.PushFront(key)
	if front, err := c.list.Front(); err == nil {
		c.cache[key] = front
	}
	c.values[key] = value
}

// Remove removes a key-value pair from the cache
func (c *LRUCache[K, V]) Remove(key K) {
	if c.threadSafe {
		c.mu.Lock()
		defer c.mu.Unlock()
	}

	if _, exists := c.cache[key]; exists {
		c.list.Remove(key)
		delete(c.cache, key)
		delete(c.values, key)
	}
}

// Clear removes all items from the cache
func (c *LRUCache[K, V]) Clear() {
	if c.threadSafe {
		c.mu.Lock()
		defer c.mu.Unlock()
	}

	c.list.Clear()
	c.cache = make(map[K]*linkedlist.DNode[K])
	c.values = make(map[K]V)
}

// Len returns the current number of items in the cache
func (c *LRUCache[K, V]) Len() int {
	if c.threadSafe {
		c.mu.RLock()
		defer c.mu.RUnlock()
	}
	return c.list.Len()
}

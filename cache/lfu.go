package cache

import (
	"sync"
)

type frequencyNode[K comparable] struct {
	freq  int
	items map[K]struct{}
	prev  *frequencyNode[K]
	next  *frequencyNode[K]
}

type LFUCache[K comparable, V any] struct {
	capacity   int
	cache      map[K]*frequencyNode[K]
	freqList   *frequencyNode[K]
	values     map[K]V
	threadSafe bool
	mu         sync.RWMutex
}

// NewLFUCache creates a new LFU cache with the specified capacity
func NewLFUCache[K comparable, V any](capacity int, threadSafe ...bool) *LFUCache[K, V] {
	isThreadSafe := true
	if len(threadSafe) > 0 {
		isThreadSafe = threadSafe[0]
	}
	return &LFUCache[K, V]{
		capacity:   capacity,
		cache:      make(map[K]*frequencyNode[K]),
		freqList:   nil,
		values:     make(map[K]V),
		threadSafe: isThreadSafe,
	}
}

// Get retrieves a value from the cache and increments its frequency
func (c *LFUCache[K, V]) Get(key K) (V, bool) {
	if c.threadSafe {
		c.mu.Lock()
		defer c.mu.Unlock()
	}

	if node, exists := c.cache[key]; exists {
		c.updateFrequency(key, node)
		return c.values[key], true
	}
	var zero V
	return zero, false
}

// updateFrequency moves a key to the next frequency node
func (c *LFUCache[K, V]) updateFrequency(key K, node *frequencyNode[K]) {
	// Remove from current frequency node
	delete(node.items, key)

	// If node becomes empty and it's not the head, remove it
	if len(node.items) == 0 && node != c.freqList {
		if node.prev != nil {
			node.prev.next = node.next
		}
		if node.next != nil {
			node.next.prev = node.prev
		}
	}

	// Create or get next frequency node
	nextFreq := node.freq + 1
	var nextNode *frequencyNode[K]

	if node.next != nil && node.next.freq == nextFreq {
		nextNode = node.next
	} else {
		nextNode = &frequencyNode[K]{
			freq:  nextFreq,
			items: make(map[K]struct{}),
			prev:  node,
			next:  node.next,
		}
		if node.next != nil {
			node.next.prev = nextNode
		}
		node.next = nextNode
	}

	// Add to next frequency node
	nextNode.items[key] = struct{}{}
	c.cache[key] = nextNode
}

// Put adds or updates a value in the cache
func (c *LFUCache[K, V]) Put(key K, value V) {
	if c.threadSafe {
		c.mu.Lock()
		defer c.mu.Unlock()
	}

	// If key exists, update it
	if node, exists := c.cache[key]; exists {
		c.updateFrequency(key, node)
	} else {
		// If cache is full, remove the least frequently used item
		if len(c.values) >= c.capacity {
			// Find the first non-empty frequency node
			current := c.freqList
			for current != nil && len(current.items) == 0 {
				current = current.next
			}

			if current != nil {
				// Get any key from the items map
				var keyToRemove K
				for k := range current.items {
					keyToRemove = k
					break
				}

				// Remove the least frequently used item
				delete(current.items, keyToRemove)
				delete(c.cache, keyToRemove)
				delete(c.values, keyToRemove)
			}
		}

		// Add to frequency 1 node
		if c.freqList == nil || c.freqList.freq != 1 {
			c.freqList = &frequencyNode[K]{
				freq:  1,
				items: make(map[K]struct{}),
				next:  c.freqList,
			}
			if c.freqList.next != nil {
				c.freqList.next.prev = c.freqList
			}
		}
		c.freqList.items[key] = struct{}{}
		c.cache[key] = c.freqList
	}

	c.values[key] = value
}

// Remove removes a key-value pair from the cache
func (c *LFUCache[K, V]) Remove(key K) {
	if c.threadSafe {
		c.mu.Lock()
		defer c.mu.Unlock()
	}

	if node, exists := c.cache[key]; exists {
		delete(node.items, key)

		// If node becomes empty and it's not the head, remove it
		if len(node.items) == 0 && node != c.freqList {
			if node.prev != nil {
				node.prev.next = node.next
			}
			if node.next != nil {
				node.next.prev = node.prev
			}
		}

		delete(c.cache, key)
		delete(c.values, key)
	}
}

// Clear removes all items from the cache
func (c *LFUCache[K, V]) Clear() {
	if c.threadSafe {
		c.mu.Lock()
		defer c.mu.Unlock()
	}

	c.freqList = nil
	c.cache = make(map[K]*frequencyNode[K])
	c.values = make(map[K]V)
}

// Len returns the current number of items in the cache
func (c *LFUCache[K, V]) Len() int {
	if c.threadSafe {
		c.mu.RLock()
		defer c.mu.RUnlock()
	}
	return len(c.values)
}

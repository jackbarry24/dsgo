package maps

import (
	"slices"
	"sync"
)

// OrderedMap is a map that maintains keys in insertion order.
// It is not safe for concurrent use.
type OrderedMap[K comparable, V any] struct {
	keys   []K
	values []V
	index  map[K]int // Maps key to its position in the slices
}

func NewOrderedMap[K comparable, V any]() *OrderedMap[K, V] {
	return &OrderedMap[K, V]{
		keys:   make([]K, 0),
		values: make([]V, 0),
		index:  make(map[K]int),
	}
}

func (m *OrderedMap[K, V]) Get(key K) (V, bool) {
	if pos, exists := m.index[key]; exists {
		return m.values[pos], true
	}
	var zero V
	return zero, false
}

func (m *OrderedMap[K, V]) Set(key K, value V) {
	if pos, exists := m.index[key]; exists {
		m.values[pos] = value
		return
	}
	pos := len(m.keys)
	m.keys = append(m.keys, key)
	m.values = append(m.values, value)
	m.index[key] = pos
}

func (m *OrderedMap[K, V]) Delete(key K) {
	pos, exists := m.index[key]
	if !exists {
		return
	}

	// Remove from slices
	m.keys = slices.Delete(m.keys, pos, pos+1)
	m.values = slices.Delete(m.values, pos, pos+1)
	delete(m.index, key)

	// Update indices for all elements after the deleted one
	for i := pos; i < len(m.keys); i++ {
		m.index[m.keys[i]] = i
	}
}

func (m *OrderedMap[K, V]) Len() int {
	return len(m.keys)
}

func (m *OrderedMap[K, V]) IsEmpty() bool {
	return len(m.keys) == 0
}

func (m *OrderedMap[K, V]) Next(key K) (K, V, bool) {
	pos, exists := m.index[key]
	if !exists || pos+1 >= len(m.keys) {
		var zeroK K
		var zeroV V
		return zeroK, zeroV, false
	}
	return m.keys[pos+1], m.values[pos+1], true
}

func (m *OrderedMap[K, V]) Prev(key K) (K, V, bool) {
	pos, exists := m.index[key]
	if !exists || pos <= 0 {
		var zeroK K
		var zeroV V
		return zeroK, zeroV, false
	}
	return m.keys[pos-1], m.values[pos-1], true
}

// Keys returns a slice of all keys in insertion order
func (m *OrderedMap[K, V]) Keys() []K {
	keys := make([]K, len(m.keys))
	copy(keys, m.keys)
	return keys
}

// Values returns a slice of all values in insertion order
func (m *OrderedMap[K, V]) Values() []V {
	values := make([]V, len(m.values))
	copy(values, m.values)
	return values
}

// Range iterates over the map in insertion order
func (m *OrderedMap[K, V]) Range(f func(key K, value V) bool) {
	for i, key := range m.keys {
		if !f(key, m.values[i]) {
			break
		}
	}
}

// SafeOrderedMap is a thread-safe wrapper around OrderedMap.
type SafeOrderedMap[K comparable, V any] struct {
	mu    sync.RWMutex
	inner *OrderedMap[K, V]
}

func NewSafeOrderedMap[K comparable, V any]() *SafeOrderedMap[K, V] {
	return &SafeOrderedMap[K, V]{
		inner: NewOrderedMap[K, V](),
	}
}

func (m *SafeOrderedMap[K, V]) Get(key K) (V, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.inner.Get(key)
}

func (m *SafeOrderedMap[K, V]) Set(key K, value V) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.inner.Set(key, value)
}

func (m *SafeOrderedMap[K, V]) Delete(key K) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.inner.Delete(key)
}

func (m *SafeOrderedMap[K, V]) Len() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.inner.Len()
}

func (m *SafeOrderedMap[K, V]) IsEmpty() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.inner.IsEmpty()
}

func (m *SafeOrderedMap[K, V]) Next(key K) (K, V, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.inner.Next(key)
}

func (m *SafeOrderedMap[K, V]) Prev(key K) (K, V, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.inner.Prev(key)
}

func (m *SafeOrderedMap[K, V]) Keys() []K {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.inner.Keys()
}

func (m *SafeOrderedMap[K, V]) Values() []V {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.inner.Values()
}

// Range iterates over the map in insertion order
func (m *SafeOrderedMap[K, V]) Range(f func(key K, value V) bool) {
	m.mu.RLock()
	keys := make([]K, len(m.inner.keys))
	values := make([]V, len(m.inner.values))
	copy(keys, m.inner.keys)
	copy(values, m.inner.values)
	m.mu.RUnlock()

	// Iterate over the copies without holding the lock
	for i := range keys {
		if !f(keys[i], values[i]) {
			return
		}
	}
}

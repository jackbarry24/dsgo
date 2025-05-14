package maps

import "sync"

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

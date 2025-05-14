package maps

import (
	"sync"
)

type SafeSortedMap[K Ordered, V any] struct {
	mu    sync.RWMutex
	inner *SortedMap[K, V]
}

func NewSafeSortedMap[K Ordered, V any]() *SafeSortedMap[K, V] {
	return &SafeSortedMap[K, V]{
		inner: NewSortedMap[K, V](),
	}
}

func (m *SafeSortedMap[K, V]) Get(key K) (V, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.inner.Get(key)
}

func (m *SafeSortedMap[K, V]) Set(key K, value V) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.inner.Set(key, value)
}

func (m *SafeSortedMap[K, V]) Delete(key K) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.inner.Delete(key)
}

func (m *SafeSortedMap[K, V]) Len() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.inner.Len()
}

func (m *SafeSortedMap[K, V]) IsEmpty() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.inner.IsEmpty()
}

func (m *SafeSortedMap[K, V]) Next(key K) (K, V, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.inner.Next(key)
}

func (m *SafeSortedMap[K, V]) Prev(key K) (K, V, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.inner.Prev(key)
}

func (m *SafeSortedMap[K, V]) Keys() []K {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.inner.Keys()
}

func (m *SafeSortedMap[K, V]) Values() []V {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.inner.Values()
}

func (m *SafeSortedMap[K, V]) Range(f func(key K, value V) bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	m.inner.Range(f)
}

package maps

import (
	"slices"
	"sort"
	"sync"

	"dsgo/utils"
)

type SortedMap[K utils.Ordered, V any] struct {
	keys   []K
	values []V
	index  map[K]int
}

func NewSortedMap[K utils.Ordered, V any]() *SortedMap[K, V] {
	return &SortedMap[K, V]{
		keys:   make([]K, 0),
		values: make([]V, 0),
		index:  make(map[K]int),
	}
}

func (m *SortedMap[K, V]) Get(key K) (V, bool) {
	if pos, exists := m.index[key]; exists {
		return m.values[pos], true
	}
	var zero V
	return zero, false
}

func (m *SortedMap[K, V]) Set(key K, value V) {
	if pos, exists := m.index[key]; exists {
		m.values[pos] = value
		return
	}
	pos := sort.Search(len(m.keys), func(i int) bool {
		return m.keys[i] >= key
	})
	m.keys = slices.Insert(m.keys, pos, key)
	m.values = slices.Insert(m.values, pos, value)
	m.index[key] = pos
}

func (m *SortedMap[K, V]) Delete(key K) {
	pos, exists := m.index[key]
	if !exists {
		return
	}

	m.keys = slices.Delete(m.keys, pos, pos+1)
	m.values = slices.Delete(m.values, pos, pos+1)
	delete(m.index, key)

	for i := pos; i < len(m.keys); i++ {
		m.index[m.keys[i]] = i
	}
}

func (m *SortedMap[K, V]) Len() int {
	return len(m.keys)
}

func (m *SortedMap[K, V]) IsEmpty() bool {
	return len(m.keys) == 0
}

func (m *SortedMap[K, V]) Next(key K) (K, V, bool) {
	pos, exists := m.index[key]
	if !exists || pos+1 >= len(m.keys) {
		var zeroK K
		var zeroV V
		return zeroK, zeroV, false
	}
	return m.keys[pos+1], m.values[pos+1], true
}

func (m *SortedMap[K, V]) Prev(key K) (K, V, bool) {
	pos, exists := m.index[key]
	if !exists || pos <= 0 {
		var zeroK K
		var zeroV V
		return zeroK, zeroV, false
	}
	return m.keys[pos-1], m.values[pos-1], true
}

func (m *SortedMap[K, V]) Keys() []K {
	keys := make([]K, len(m.keys))
	copy(keys, m.keys)
	return keys
}

func (m *SortedMap[K, V]) Values() []V {
	values := make([]V, len(m.values))
	copy(values, m.values)
	return values
}

func (m *SortedMap[K, V]) Range(f func(key K, value V) bool) {
	for i, key := range m.keys {
		if !f(key, m.values[i]) {
			break
		}
	}
}

// SafeSortedMap is a thread-safe wrapper around SortedMap.
type SafeSortedMap[K utils.Ordered, V any] struct {
	mu    sync.RWMutex
	inner *SortedMap[K, V]
}

func NewSafeSortedMap[K utils.Ordered, V any]() *SafeSortedMap[K, V] {
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

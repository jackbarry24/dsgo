package maps

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
	m.keys = append(m.keys[:pos], m.keys[pos+1:]...)
	m.values = append(m.values[:pos], m.values[pos+1:]...)
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

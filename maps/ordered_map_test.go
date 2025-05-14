package maps

import (
	"testing"
)

func TestNewOrderedMap(t *testing.T) {
	m := NewOrderedMap[string, int]()
	if m == nil {
		t.Error("NewOrderedMap returned nil")
	}
	if m.Len() != 0 {
		t.Errorf("Expected empty map, got length %d", m.Len())
	}
	if !m.IsEmpty() {
		t.Error("Expected IsEmpty to return true for new map")
	}
}

func TestSetAndGet(t *testing.T) {
	m := NewOrderedMap[string, int]()

	// Test setting and getting a value
	m.Set("one", 1)
	val, exists := m.Get("one")
	if !exists {
		t.Error("Expected key 'one' to exist")
	}
	if val != 1 {
		t.Errorf("Expected value 1, got %d", val)
	}

	// Test getting non-existent key
	val, exists = m.Get("two")
	if exists {
		t.Error("Expected key 'two' to not exist")
	}
	if val != 0 {
		t.Errorf("Expected zero value 0, got %d", val)
	}

	// Test updating existing value
	m.Set("one", 11)
	val, _ = m.Get("one")
	if val != 11 {
		t.Errorf("Expected updated value 11, got %d", val)
	}
}

func TestDelete(t *testing.T) {
	m := NewOrderedMap[string, int]()
	m.Set("one", 1)
	m.Set("two", 2)
	m.Set("three", 3)

	// Test deleting existing key
	m.Delete("two")
	if m.Len() != 2 {
		t.Errorf("Expected length 2 after delete, got %d", m.Len())
	}
	if _, exists := m.Get("two"); exists {
		t.Error("Expected deleted key 'two' to not exist")
	}

	// Test deleting non-existent key
	m.Delete("four")
	if m.Len() != 2 {
		t.Errorf("Expected length to remain 2 after deleting non-existent key, got %d", m.Len())
	}
}

func TestNextAndPrev(t *testing.T) {
	m := NewOrderedMap[string, int]()
	m.Set("one", 1)
	m.Set("two", 2)
	m.Set("three", 3)

	// Test Next
	nextKey, nextVal, ok := m.Next("one")
	if !ok {
		t.Error("Expected Next to return true for 'one'")
	}
	if nextKey != "two" || nextVal != 2 {
		t.Errorf("Expected next key-value pair (two, 2), got (%s, %d)", nextKey, nextVal)
	}

	// Test Next on last element
	_, _, ok = m.Next("three")
	if ok {
		t.Error("Expected Next to return false for last element")
	}

	// Test Prev
	prevKey, prevVal, ok := m.Prev("two")
	if !ok {
		t.Error("Expected Prev to return true for 'two'")
	}
	if prevKey != "one" || prevVal != 1 {
		t.Errorf("Expected prev key-value pair (one, 1), got (%s, %d)", prevKey, prevVal)
	}

	// Test Prev on first element
	_, _, ok = m.Prev("one")
	if ok {
		t.Error("Expected Prev to return false for first element")
	}

	// Test Next/Prev on non-existent key
	_, _, ok = m.Next("four")
	if ok {
		t.Error("Expected Next to return false for non-existent key")
	}
	_, _, ok = m.Prev("four")
	if ok {
		t.Error("Expected Prev to return false for non-existent key")
	}
}

func TestKeysAndValues(t *testing.T) {
	m := NewOrderedMap[string, int]()
	m.Set("one", 1)
	m.Set("two", 2)
	m.Set("three", 3)

	// Test Keys
	keys := m.Keys()
	expectedKeys := []string{"one", "two", "three"}
	if len(keys) != len(expectedKeys) {
		t.Errorf("Expected %d keys, got %d", len(expectedKeys), len(keys))
	}
	for i, key := range keys {
		if key != expectedKeys[i] {
			t.Errorf("Expected key %s at position %d, got %s", expectedKeys[i], i, key)
		}
	}

	// Test Values
	values := m.Values()
	expectedValues := []int{1, 2, 3}
	if len(values) != len(expectedValues) {
		t.Errorf("Expected %d values, got %d", len(expectedValues), len(values))
	}
	for i, val := range values {
		if val != expectedValues[i] {
			t.Errorf("Expected value %d at position %d, got %d", expectedValues[i], i, val)
		}
	}
}

func TestRange(t *testing.T) {
	m := NewOrderedMap[string, int]()
	m.Set("one", 1)
	m.Set("two", 2)
	m.Set("three", 3)

	// Test normal iteration
	count := 0
	m.Range(func(key string, value int) bool {
		count++
		return true
	})
	if count != 3 {
		t.Errorf("Expected to iterate over 3 elements, got %d", count)
	}

	// Test early termination
	count = 0
	m.Range(func(key string, value int) bool {
		count++
		return count < 2
	})
	if count != 2 {
		t.Errorf("Expected to iterate over 2 elements, got %d", count)
	}
}

func TestOrderedMapWithDifferentTypes(t *testing.T) {
	// Test with int keys
	m1 := NewOrderedMap[int, string]()
	m1.Set(1, "one")
	m1.Set(2, "two")
	if val, _ := m1.Get(1); val != "one" {
		t.Errorf("Expected value 'one' for key 1, got %s", val)
	}

	// Test with struct keys
	type Key struct {
		ID   int
		Name string
	}
	m2 := NewOrderedMap[Key, int]()
	key := Key{ID: 1, Name: "test"}
	m2.Set(key, 42)
	if val, _ := m2.Get(key); val != 42 {
		t.Errorf("Expected value 42 for struct key, got %d", val)
	}
}

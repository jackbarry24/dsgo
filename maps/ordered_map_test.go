package maps

import (
	"sync"
	"testing"
)

func TestNewOrderedMap(t *testing.T) {
	m := NewOrderedMap[string, int](false)
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
	m := NewOrderedMap[string, int](false)

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
	m := NewOrderedMap[string, int](false)
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
	m := NewOrderedMap[string, int](false)
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
	m := NewOrderedMap[string, int](false)
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

// SafeOrderedMap tests

func TestSafeOrderedMap_BasicOperations(t *testing.T) {
	m := NewOrderedMap[string, int](true)

	// Test Set and Get
	m.Set("one", 1)
	m.Set("two", 2)
	m.Set("three", 3)

	if val, exists := m.Get("one"); !exists || val != 1 {
		t.Errorf("Get(\"one\") = (%v, %v), want (1, true)", val, exists)
	}

	// Test Len
	if m.Len() != 3 {
		t.Errorf("Len() = %v, want 3", m.Len())
	}

	// Test IsEmpty
	if m.IsEmpty() {
		t.Error("IsEmpty() = true, want false")
	}

	// Test Delete
	m.Delete("two")
	if val, exists := m.Get("two"); exists {
		t.Errorf("Get(\"two\") = (%v, %v), want (0, false)", val, exists)
	}

	// Test Keys and Values
	keys := m.Keys()
	expectedKeys := []string{"one", "three"}
	if len(keys) != len(expectedKeys) {
		t.Errorf("Keys() length = %v, want %v", len(keys), len(expectedKeys))
	}

	values := m.Values()
	expectedValues := []int{1, 3}
	if len(values) != len(expectedValues) {
		t.Errorf("Values() length = %v, want %v", len(values), len(expectedValues))
	}
}

func TestSafeOrderedMap_NextPrev(t *testing.T) {
	m := NewOrderedMap[string, int](true)
	m.Set("one", 1)
	m.Set("two", 2)
	m.Set("three", 3)

	// Test Next
	key, val, exists := m.Next("one")
	if !exists || key != "two" || val != 2 {
		t.Errorf("Next(\"one\") = (%v, %v, %v), want (\"two\", 2, true)", key, val, exists)
	}

	// Test Prev
	key, val, exists = m.Prev("three")
	if !exists || key != "two" || val != 2 {
		t.Errorf("Prev(\"three\") = (%v, %v, %v), want (\"two\", 2, true)", key, val, exists)
	}

	// Test non-existent keys
	key, val, exists = m.Next("three")
	if exists {
		t.Errorf("Next(\"three\") = (%v, %v, %v), want (\"\", 0, false)", key, val, exists)
	}

	key, val, exists = m.Prev("one")
	if exists {
		t.Errorf("Prev(\"one\") = (%v, %v, %v), want (\"\", 0, false)", key, val, exists)
	}
}

func TestSafeOrderedMap_Range(t *testing.T) {
	m := NewOrderedMap[string, int](true)
	m.Set("one", 1)
	m.Set("two", 2)
	m.Set("three", 3)

	// Test full range
	expected := map[string]int{
		"one":   1,
		"two":   2,
		"three": 3,
	}
	count := 0
	m.Range(func(key string, value int) bool {
		if expected[key] != value {
			t.Errorf("Range: key %v has value %v, want %v", key, value, expected[key])
		}
		count++
		return true
	})
	if count != len(expected) {
		t.Errorf("Range: processed %v items, want %v", count, len(expected))
	}

	// Test early exit
	count = 0
	m.Range(func(key string, value int) bool {
		count++
		return false // Exit after first item
	})
	if count != 1 {
		t.Errorf("Range with early exit: processed %v items, want 1", count)
	}
}

func TestSafeOrderedMap_Concurrent(t *testing.T) {
	m := NewOrderedMap[string, int](true)
	var wg sync.WaitGroup
	iterations := 1000

	// Concurrent writes
	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 0; i < iterations; i++ {
			m.Set("key1", i)
		}
	}()
	go func() {
		defer wg.Done()
		for i := 0; i < iterations; i++ {
			m.Set("key2", i)
		}
	}()

	// Concurrent reads
	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 0; i < iterations; i++ {
			m.Get("key1")
			m.Get("key2")
		}
	}()
	go func() {
		defer wg.Done()
		for i := 0; i < iterations; i++ {
			m.Keys()
			m.Values()
		}
	}()

	wg.Wait()

	// Verify final state
	if m.Len() != 2 {
		t.Errorf("Concurrent operations: Len() = %v, want 2", m.Len())
	}
}

func TestSafeOrderedMap_Empty(t *testing.T) {
	m := NewOrderedMap[string, int](true)

	if !m.IsEmpty() {
		t.Error("IsEmpty() = false, want true")
	}

	if m.Len() != 0 {
		t.Errorf("Len() = %v, want 0", m.Len())
	}

	// Test operations on empty map
	if val, exists := m.Get("nonexistent"); exists {
		t.Errorf("Get(\"nonexistent\") = (%v, %v), want (0, false)", val, exists)
	}

	key, val, exists := m.Next("nonexistent")
	if exists {
		t.Errorf("Next(\"nonexistent\") = (%v, %v, %v), want (\"\", 0, false)", key, val, exists)
	}

	key, val, exists = m.Prev("nonexistent")
	if exists {
		t.Errorf("Prev(\"nonexistent\") = (%v, %v, %v), want (\"\", 0, false)", key, val, exists)
	}

	// Test range on empty map
	count := 0
	m.Range(func(key string, value int) bool {
		count++
		return true
	})
	if count != 0 {
		t.Errorf("Range on empty map: processed %v items, want 0", count)
	}
}

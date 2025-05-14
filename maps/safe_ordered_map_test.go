package maps

import (
	"sync"
	"testing"
)

func TestSafeOrderedMap_BasicOperations(t *testing.T) {
	m := NewSafeOrderedMap[string, int]()

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
	m := NewSafeOrderedMap[string, int]()
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
	m := NewSafeOrderedMap[string, int]()
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
	m := NewSafeOrderedMap[string, int]()
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
	m := NewSafeOrderedMap[string, int]()

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

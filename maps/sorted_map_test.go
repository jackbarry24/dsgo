package maps

import (
	"sync"
	"testing"
)

func TestNewSortedMap(t *testing.T) {
	m := NewSortedMap[int, string]()
	if m == nil {
		t.Error("NewSortedMap returned nil")
	}
	if !m.IsEmpty() {
		t.Error("New map should be empty")
	}
	if m.Len() != 0 {
		t.Error("New map should have length 0")
	}
}

func TestSortedMap_SetGet(t *testing.T) {
	m := NewSortedMap[int, string]()

	// Test setting and getting values
	m.Set(1, "one")
	m.Set(2, "two")
	m.Set(3, "three")

	if val, exists := m.Get(1); !exists || val != "one" {
		t.Errorf("Get(1) = %v, %v; want 'one', true", val, exists)
	}
	if val, exists := m.Get(2); !exists || val != "two" {
		t.Errorf("Get(2) = %v, %v; want 'two', true", val, exists)
	}
	if val, exists := m.Get(3); !exists || val != "three" {
		t.Errorf("Get(3) = %v, %v; want 'three', true", val, exists)
	}

	// Test getting non-existent key
	if val, exists := m.Get(4); exists {
		t.Errorf("Get(4) = %v, %v; want '', false", val, exists)
	}

	// Test updating existing value
	m.Set(1, "ONE")
	if val, exists := m.Get(1); !exists || val != "ONE" {
		t.Errorf("Get(1) after update = %v, %v; want 'ONE', true", val, exists)
	}
}

func TestSortedMap_Delete(t *testing.T) {
	m := NewSortedMap[int, string]()
	m.Set(1, "one")
	m.Set(2, "two")
	m.Set(3, "three")

	// Test deleting existing key
	m.Delete(2)
	if val, exists := m.Get(2); exists {
		t.Errorf("Get(2) after delete = %v, %v; want '', false", val, exists)
	}
	if m.Len() != 2 {
		t.Errorf("Len() after delete = %d; want 2", m.Len())
	}

	// Test deleting non-existent key
	m.Delete(4)
	if m.Len() != 2 {
		t.Errorf("Len() after deleting non-existent key = %d; want 2", m.Len())
	}
}

func TestSortedMap_NextPrev(t *testing.T) {
	m := NewSortedMap[int, string]()
	m.Set(1, "one")
	m.Set(2, "two")
	m.Set(3, "three")

	// Test Next
	key, val, exists := m.Next(1)
	if !exists || key != 2 || val != "two" {
		t.Errorf("Next(1) = %v, %v, %v; want 2, 'two', true", key, val, exists)
	}

	key, val, exists = m.Next(3)
	if exists {
		t.Errorf("Next(3) = %v, %v, %v; want 0, '', false", key, val, exists)
	}

	// Test Prev
	key, val, exists = m.Prev(2)
	if !exists || key != 1 || val != "one" {
		t.Errorf("Prev(2) = %v, %v, %v; want 1, 'one', true", key, val, exists)
	}

	key, val, exists = m.Prev(1)
	if exists {
		t.Errorf("Prev(1) = %v, %v, %v; want 0, '', false", key, val, exists)
	}
}

func TestSortedMap_KeysValues(t *testing.T) {
	m := NewSortedMap[int, string]()
	m.Set(3, "three")
	m.Set(1, "one")
	m.Set(2, "two")

	keys := m.Keys()
	expectedKeys := []int{1, 2, 3}
	for i, k := range keys {
		if k != expectedKeys[i] {
			t.Errorf("Keys()[%d] = %v; want %v", i, k, expectedKeys[i])
		}
	}

	values := m.Values()
	expectedValues := []string{"one", "two", "three"}
	for i, v := range values {
		if v != expectedValues[i] {
			t.Errorf("Values()[%d] = %v; want %v", i, v, expectedValues[i])
		}
	}
}

func TestSortedMap_Range(t *testing.T) {
	m := NewSortedMap[int, string]()
	m.Set(1, "one")
	m.Set(2, "two")
	m.Set(3, "three")

	// Test full range
	count := 0
	m.Range(func(key int, value string) bool {
		count++
		return true
	})
	if count != 3 {
		t.Errorf("Range visited %d items; want 3", count)
	}

	// Test early exit
	count = 0
	m.Range(func(key int, value string) bool {
		count++
		return count < 2
	})
	if count != 2 {
		t.Errorf("Range with early exit visited %d items; want 2", count)
	}
}

func TestSortedMap_Ordering(t *testing.T) {
	m := NewSortedMap[string, int]()

	// Test string ordering
	m.Set("zebra", 1)
	m.Set("apple", 2)
	m.Set("banana", 3)

	keys := m.Keys()
	expectedKeys := []string{"apple", "banana", "zebra"}
	for i, k := range keys {
		if k != expectedKeys[i] {
			t.Errorf("Keys()[%d] = %v; want %v", i, k, expectedKeys[i])
		}
	}
}

// SafeSortedMap tests

func TestSafeSortedMap_BasicOperations(t *testing.T) {
	m := NewSafeSortedMap[int, string]()

	// Test empty state
	if !m.IsEmpty() {
		t.Error("Expected map to be empty")
	}
	if m.Len() != 0 {
		t.Error("Expected length to be 0")
	}

	// Test Set and Get
	m.Set(1, "one")
	m.Set(2, "two")
	m.Set(3, "three")

	if m.Len() != 3 {
		t.Errorf("Expected length to be 3, got %d", m.Len())
	}

	if val, ok := m.Get(1); !ok || val != "one" {
		t.Errorf("Expected to get 'one', got %v, %v", val, ok)
	}

	// Test Delete
	m.Delete(2)
	if val, ok := m.Get(2); ok {
		t.Errorf("Expected key 2 to be deleted, got %v", val)
	}

	// Test Next and Prev
	key, val, ok := m.Next(1)
	if !ok || key != 3 || val != "three" {
		t.Errorf("Expected Next(1) to return (3, 'three'), got (%v, %v, %v)", key, val, ok)
	}

	key, val, ok = m.Prev(3)
	if !ok || key != 1 || val != "one" {
		t.Errorf("Expected Prev(3) to return (1, 'one'), got (%v, %v, %v)", key, val, ok)
	}
}

func TestSafeSortedMap_KeysAndValues(t *testing.T) {
	m := NewSafeSortedMap[int, string]()
	m.Set(1, "one")
	m.Set(2, "two")
	m.Set(3, "three")

	keys := m.Keys()
	expectedKeys := []int{1, 2, 3}
	if len(keys) != len(expectedKeys) {
		t.Errorf("Expected %d keys, got %d", len(expectedKeys), len(keys))
	}
	for i, k := range keys {
		if k != expectedKeys[i] {
			t.Errorf("Expected key %d, got %d", expectedKeys[i], k)
		}
	}

	values := m.Values()
	expectedValues := []string{"one", "two", "three"}
	if len(values) != len(expectedValues) {
		t.Errorf("Expected %d values, got %d", len(expectedValues), len(values))
	}
	for i, v := range values {
		if v != expectedValues[i] {
			t.Errorf("Expected value %s, got %s", expectedValues[i], v)
		}
	}
}

func TestSafeSortedMap_Range(t *testing.T) {
	m := NewSafeSortedMap[int, string]()
	m.Set(1, "one")
	m.Set(2, "two")
	m.Set(3, "three")

	visited := make(map[int]bool)
	m.Range(func(key int, value string) bool {
		visited[key] = true
		return true
	})

	if len(visited) != 3 {
		t.Errorf("Expected to visit 3 items, got %d", len(visited))
	}
	for i := 1; i <= 3; i++ {
		if !visited[i] {
			t.Errorf("Expected to visit key %d", i)
		}
	}
}

func TestSafeSortedMap_Concurrent(t *testing.T) {
	m := NewSafeSortedMap[int, int]()
	var wg sync.WaitGroup
	iterations := 1000
	goroutines := 10

	// Concurrent writes
	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func(start int) {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				key := start*iterations + j
				m.Set(key, key*2)
			}
		}(i)
	}

	// Concurrent reads
	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				m.Len()
				m.IsEmpty()
				m.Keys()
				m.Values()
			}
		}()
	}

	wg.Wait()

	// Verify final state
	expectedLen := goroutines * iterations
	if m.Len() != expectedLen {
		t.Errorf("Expected length %d, got %d", expectedLen, m.Len())
	}

	// Verify all values
	m.Range(func(key int, value int) bool {
		if value != key*2 {
			t.Errorf("Expected value %d for key %d, got %d", key*2, key, value)
		}
		return true
	})
}

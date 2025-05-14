package maps

import (
	"sync"
	"testing"
)

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

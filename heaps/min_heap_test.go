package heaps

import (
	"sync"
	"testing"
)

func TestMinHeapBasicOperations(t *testing.T) {
	heap := NewMinHeap[int](func(a, b int) bool { return a < b }, false)

	// Test empty heap
	if !heap.IsEmpty() {
		t.Error("New heap should be empty")
	}
	if heap.Size() != 0 {
		t.Error("New heap should have size 0")
	}
	if _, ok := heap.Peek(); ok {
		t.Error("Peek on empty heap should return false")
	}
	if _, ok := heap.Pop(); ok {
		t.Error("Pop on empty heap should return false")
	}

	// Test Push and Peek
	heap.Push(5)
	if heap.IsEmpty() {
		t.Error("Heap should not be empty after Push")
	}
	if heap.Size() != 1 {
		t.Error("Heap size should be 1 after Push")
	}
	if val, ok := heap.Peek(); !ok || val != 5 {
		t.Error("Peek should return 5")
	}

	// Test multiple Pushes
	heap.Push(3)
	heap.Push(7)
	heap.Push(1)
	heap.Push(4)

	if heap.Size() != 5 {
		t.Error("Heap size should be 5")
	}

	// Test Pop order
	expected := []int{1, 3, 4, 5, 7}
	for _, exp := range expected {
		if val, ok := heap.Pop(); !ok || val != exp {
			t.Errorf("Expected %d, got %d", exp, val)
		}
	}

	// Test empty after all Pops
	if !heap.IsEmpty() {
		t.Error("Heap should be empty after all Pops")
	}
}

func TestMinHeapWithCustomType(t *testing.T) {
	type Person struct {
		Age  int
		Name string
	}

	heap := NewMinHeap[Person](func(a, b Person) bool {
		return a.Age < b.Age
	}, false)

	people := []Person{
		{Age: 30, Name: "Alice"},
		{Age: 20, Name: "Bob"},
		{Age: 25, Name: "Charlie"},
	}

	for _, p := range people {
		heap.Push(p)
	}

	// Test order by age
	expected := []int{20, 25, 30}
	for _, exp := range expected {
		if val, ok := heap.Pop(); !ok || val.Age != exp {
			t.Errorf("Expected age %d, got %d", exp, val.Age)
		}
	}
}

func TestSafeMinHeapBasicOperations(t *testing.T) {
	heap := NewMinHeap[int](func(a, b int) bool { return a < b }, true)

	// Test empty heap
	if !heap.IsEmpty() {
		t.Error("New heap should be empty")
	}
	if heap.Size() != 0 {
		t.Error("New heap should have size 0")
	}
	if _, ok := heap.Peek(); ok {
		t.Error("Peek on empty heap should return false")
	}
	if _, ok := heap.Pop(); ok {
		t.Error("Pop on empty heap should return false")
	}

	// Test Push and Peek
	heap.Push(5)
	if heap.IsEmpty() {
		t.Error("Heap should not be empty after Push")
	}
	if heap.Size() != 1 {
		t.Error("Heap size should be 1 after Push")
	}
	if val, ok := heap.Peek(); !ok || val != 5 {
		t.Error("Peek should return 5")
	}

	// Test multiple Pushes
	heap.Push(3)
	heap.Push(7)
	heap.Push(1)
	heap.Push(4)

	if heap.Size() != 5 {
		t.Error("Heap size should be 5")
	}

	// Test Pop order
	expected := []int{1, 3, 4, 5, 7}
	for _, exp := range expected {
		if val, ok := heap.Pop(); !ok || val != exp {
			t.Errorf("Expected %d, got %d", exp, val)
		}
	}

	// Test empty after all Pops
	if !heap.IsEmpty() {
		t.Error("Heap should be empty after all Pops")
	}
}

func TestSafeMinHeapConcurrent(t *testing.T) {
	heap := NewMinHeap[int](func(a, b int) bool { return a < b }, true)
	var wg sync.WaitGroup
	concurrentWriters := 10
	itemsPerWriter := 100

	// Concurrent writes
	wg.Add(concurrentWriters)
	for i := 0; i < concurrentWriters; i++ {
		go func(start int) {
			defer wg.Done()
			for j := 0; j < itemsPerWriter; j++ {
				heap.Push(start*itemsPerWriter + j)
			}
		}(i)
	}
	wg.Wait()

	// Verify size
	expectedSize := concurrentWriters * itemsPerWriter
	if heap.Size() != expectedSize {
		t.Errorf("Expected size %d, got %d", expectedSize, heap.Size())
	}

	// Concurrent reads and writes
	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 0; i < 1000; i++ {
			heap.Push(i)
		}
	}()
	go func() {
		defer wg.Done()
		for i := 0; i < 1000; i++ {
			heap.Pop()
		}
	}()
	wg.Wait()

	// Verify heap property
	last := -1
	for !heap.IsEmpty() {
		if val, ok := heap.Pop(); ok {
			if val < last {
				t.Errorf("Heap property violated: %d came after %d", val, last)
			}
			last = val
		}
	}
}

func TestMinHeapEdgeCases(t *testing.T) {
	heap := NewMinHeap[int](func(a, b int) bool { return a < b }, false)

	// Test with negative numbers
	heap.Push(-5)
	heap.Push(-3)
	heap.Push(-7)
	if val, ok := heap.Pop(); !ok || val != -7 {
		t.Error("Expected -7, got", val)
	}

	// Test with duplicate values
	heap = NewMinHeap[int](func(a, b int) bool { return a < b }, false)
	heap.Push(5)
	heap.Push(5)
	heap.Push(5)
	if heap.Size() != 3 {
		t.Error("Heap should handle duplicate values")
	}
	if val, ok := heap.Pop(); !ok || val != 5 {
		t.Error("Expected 5, got", val)
	}

	// Test with zero values
	heap = NewMinHeap[int](func(a, b int) bool { return a < b }, false)
	heap.Push(0)
	heap.Push(0)
	if val, ok := heap.Pop(); !ok || val != 0 {
		t.Error("Expected 0, got", val)
	}
}

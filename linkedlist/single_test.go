package linkedlist

import (
	"sync"
	"testing"
)

func TestSingleLinkedListBasic(t *testing.T) {
	// Test non-thread-safe version
	list := NewSingleLinkedList[int](false)
	testSingleBasicOperations(t, list)

	// Test thread-safe version
	list = NewSingleLinkedList[int](true)
	testSingleBasicOperations(t, list)
}

func testSingleBasicOperations(t *testing.T, list *SingleLinkedList[int]) {
	// Test empty list
	if list.Len() != 0 {
		t.Errorf("Expected length 0, got %d", list.Len())
	}

	// Test PushBack
	list.PushBack(1)
	list.PushBack(2)
	list.PushBack(3)
	if list.Len() != 3 {
		t.Errorf("Expected length 3, got %d", list.Len())
	}

	// Test Front
	if front, err := list.Front(); err != nil || front.value != 1 {
		t.Errorf("Expected front value 1, got %v", front)
	}

	// Test Back
	if back, err := list.Back(); err != nil || back.value != 3 {
		t.Errorf("Expected back value 3, got %v", back)
	}

	// Test Remove
	if err := list.Remove(2); err != nil {
		t.Errorf("Failed to remove value: %v", err)
	}
	if list.Len() != 2 {
		t.Errorf("Expected length 2 after removal, got %d", list.Len())
	}

	// Test Contains
	if !list.Contains(1) || list.Contains(2) || !list.Contains(3) {
		t.Error("Contains test failed")
	}

	// Test Clear
	list.Clear()
	if list.Len() != 0 {
		t.Errorf("Expected length 0 after clear, got %d", list.Len())
	}
}

func TestSingleLinkedListConcurrent(t *testing.T) {
	list := NewSingleLinkedList[int](true)
	var wg sync.WaitGroup
	iterations := 1000

	// Test concurrent PushBack
	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 0; i < iterations; i++ {
			list.PushBack(i)
		}
	}()
	go func() {
		defer wg.Done()
		for i := iterations; i < iterations*2; i++ {
			list.PushBack(i)
		}
	}()
	wg.Wait()

	if list.Len() != iterations*2 {
		t.Errorf("Expected length %d, got %d", iterations*2, list.Len())
	}

	// Test concurrent Remove
	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 0; i < iterations; i++ {
			list.Remove(i)
		}
	}()
	go func() {
		defer wg.Done()
		for i := iterations; i < iterations*2; i++ {
			list.Remove(i)
		}
	}()
	wg.Wait()

	if list.Len() != 0 {
		t.Errorf("Expected length 0 after concurrent removal, got %d", list.Len())
	}
}

func TestSingleLinkedListPushFront(t *testing.T) {
	list := NewSingleLinkedList[int](false)

	// Test PushFront
	list.PushFront(1)
	list.PushFront(2)
	list.PushFront(3)

	// Verify order
	if front, err := list.Front(); err != nil || front.value != 3 {
		t.Errorf("Expected front value 3, got %v", front)
	}
	if back, err := list.Back(); err != nil || back.value != 1 {
		t.Errorf("Expected back value 1, got %v", back)
	}
	if list.Len() != 3 {
		t.Errorf("Expected length 3, got %d", list.Len())
	}
}

func TestSingleLinkedListAt(t *testing.T) {
	list := NewSingleLinkedList[int](false)
	values := []int{1, 2, 3, 4, 5}

	// Add values
	for _, v := range values {
		list.PushBack(v)
	}

	// Test valid indices
	for i, v := range values {
		if val, err := list.At(i); err != nil || val != v {
			t.Errorf("At(%d): expected %d, got %v, err: %v", i, v, val, err)
		}
	}

	// Test invalid indices
	if _, err := list.At(-1); err == nil {
		t.Error("Expected error for negative index")
	}
	if _, err := list.At(len(values)); err == nil {
		t.Error("Expected error for index out of bounds")
	}
}

func TestSingleLinkedListForEach(t *testing.T) {
	list := NewSingleLinkedList[int](false)
	values := []int{1, 2, 3, 4, 5}

	// Add values
	for _, v := range values {
		list.PushBack(v)
	}

	// Test iteration
	index := 0
	list.ForEach(func(v int) {
		if v != values[index] {
			t.Errorf("Expected %d, got %d", values[index], v)
		}
		index++
	})

	if index != len(values) {
		t.Errorf("Expected %d iterations, got %d", len(values), index)
	}
}

func TestSingleLinkedListRemoveEdgeCases(t *testing.T) {
	list := NewSingleLinkedList[int](false)

	// Test removing from empty list
	if err := list.Remove(1); err != ErrEmptyList {
		t.Errorf("Expected ErrEmptyList, got %v", err)
	}

	// Test removing non-existent value
	list.PushBack(1)
	if err := list.Remove(2); err != ErrNotFound {
		t.Errorf("Expected ErrNotFound, got %v", err)
	}

	// Test removing last element
	list.PushBack(2)
	if err := list.Remove(2); err != nil {
		t.Errorf("Failed to remove last element: %v", err)
	}
	if list.Len() != 1 {
		t.Errorf("Expected length 1, got %d", list.Len())
	}
	if back, err := list.Back(); err != nil || back.value != 1 {
		t.Errorf("Expected back value 1, got %v", back)
	}
}

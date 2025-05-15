package linkedlist

import (
	"sync"
	"testing"
)

func TestDoubleLinkedListBasic(t *testing.T) {
	// Test non-thread-safe version
	list := NewDoubleLinkedList[int](false)
	testBasicOperations(t, list)

	// Test thread-safe version
	list = NewDoubleLinkedList[int](true)
	testBasicOperations(t, list)
}

func testBasicOperations(t *testing.T, list *DoubleLinkedList[int]) {
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
	if front, err := list.Front(); err != nil || front.GetValue() != 1 {
		t.Errorf("Expected front value 1, got %v", front)
	}

	// Test Back
	if back, err := list.Back(); err != nil || back.GetValue() != 3 {
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

func TestDoubleLinkedListConcurrent(t *testing.T) {
	list := NewDoubleLinkedList[int](true)
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

func TestDoubleLinkedListInsertOperations(t *testing.T) {
	list := NewDoubleLinkedList[int](false)

	// Test InsertAfter
	list.PushBack(1)
	list.PushBack(3)
	if err := list.InsertAfter(1, 2); err != nil {
		t.Errorf("Failed to insert after: %v", err)
	}
	if list.Len() != 3 {
		t.Errorf("Expected length 3, got %d", list.Len())
	}

	// Test InsertBefore
	if err := list.InsertBefore(3, 2); err != nil {
		t.Errorf("Failed to insert before: %v", err)
	}
	if list.Len() != 4 {
		t.Errorf("Expected length 4, got %d", list.Len())
	}

	// Test error cases
	if err := list.InsertAfter(5, 6); err != ErrNotFound {
		t.Errorf("Expected ErrNotFound, got %v", err)
	}
	if err := list.InsertBefore(5, 6); err != ErrNotFound {
		t.Errorf("Expected ErrNotFound, got %v", err)
	}
}

func TestDoubleLinkedListForEach(t *testing.T) {
	list := NewDoubleLinkedList[int](false)
	values := []int{1, 2, 3, 4, 5}

	for _, v := range values {
		list.PushBack(v)
	}

	// Test forward iteration
	index := 0
	list.ForEach(func(v int) {
		if v != values[index] {
			t.Errorf("Expected %d, got %d", values[index], v)
		}
		index++
	})

	// Test reverse iteration
	index = len(values) - 1
	list.ForEachReverse(func(v int) {
		if v != values[index] {
			t.Errorf("Expected %d, got %d", values[index], v)
		}
		index--
	})
}

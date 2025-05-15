package heaps

import (
	"sync"
	"testing"
)

func TestPriorityQueueBasicOperations(t *testing.T) {
	pq := NewPriorityQueue[string]()

	// Test empty queue
	if !pq.IsEmpty() {
		t.Error("New queue should be empty")
	}
	if pq.Size() != 0 {
		t.Error("New queue should have size 0")
	}
	if _, _, ok := pq.Peek(); ok {
		t.Error("Peek on empty queue should return false")
	}
	if _, _, ok := pq.Dequeue(); ok {
		t.Error("Dequeue on empty queue should return false")
	}

	// Test Enqueue and Peek
	pq.Enqueue("task1", 5)
	if pq.IsEmpty() {
		t.Error("Queue should not be empty after Enqueue")
	}
	if pq.Size() != 1 {
		t.Error("Queue size should be 1 after Enqueue")
	}
	if val, prio, ok := pq.Peek(); !ok || val != "task1" || prio != 5 {
		t.Error("Peek should return task1 with priority 5")
	}

	// Test multiple Enqueues
	pq.Enqueue("task2", 3)
	pq.Enqueue("task3", 7)
	pq.Enqueue("task4", 1)
	pq.Enqueue("task5", 4)

	if pq.Size() != 5 {
		t.Error("Queue size should be 5")
	}

	// Test Dequeue order (should be by priority)
	expected := []struct {
		value    string
		priority int
	}{
		{"task4", 1},
		{"task2", 3},
		{"task5", 4},
		{"task1", 5},
		{"task3", 7},
	}

	for _, exp := range expected {
		if val, prio, ok := pq.Dequeue(); !ok || val != exp.value || prio != exp.priority {
			t.Errorf("Expected %s with priority %d, got %s with priority %d",
				exp.value, exp.priority, val, prio)
		}
	}

	// Test empty after all Dequeues
	if !pq.IsEmpty() {
		t.Error("Queue should be empty after all Dequeues")
	}
}

func TestMaxPriorityQueue(t *testing.T) {
	pq := NewMaxPriorityQueue[string]()

	// Test multiple Enqueues
	pq.Enqueue("task1", 5)
	pq.Enqueue("task2", 3)
	pq.Enqueue("task3", 7)
	pq.Enqueue("task4", 1)
	pq.Enqueue("task5", 4)

	// Test Dequeue order (should be by priority, highest first)
	expected := []struct {
		value    string
		priority int
	}{
		{"task3", 7},
		{"task1", 5},
		{"task5", 4},
		{"task2", 3},
		{"task4", 1},
	}

	for _, exp := range expected {
		if val, prio, ok := pq.Dequeue(); !ok || val != exp.value || prio != exp.priority {
			t.Errorf("Expected %s with priority %d, got %s with priority %d",
				exp.value, exp.priority, val, prio)
		}
	}
}

func TestSafePriorityQueueConcurrent(t *testing.T) {
	pq := NewSafePriorityQueue[int]()
	var wg sync.WaitGroup
	concurrentWriters := 10
	itemsPerWriter := 100

	// Concurrent writes
	wg.Add(concurrentWriters)
	for i := 0; i < concurrentWriters; i++ {
		go func(start int) {
			defer wg.Done()
			for j := 0; j < itemsPerWriter; j++ {
				pq.Enqueue(start*itemsPerWriter+j, j)
			}
		}(i)
	}
	wg.Wait()

	// Verify size
	expectedSize := concurrentWriters * itemsPerWriter
	if pq.Size() != expectedSize {
		t.Errorf("Expected size %d, got %d", expectedSize, pq.Size())
	}

	// Concurrent reads and writes
	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 0; i < 1000; i++ {
			pq.Enqueue(i, i)
		}
	}()
	go func() {
		defer wg.Done()
		for i := 0; i < 1000; i++ {
			pq.Dequeue()
		}
	}()
	wg.Wait()

	// Verify priority order
	lastPriority := -1
	for !pq.IsEmpty() {
		if _, prio, ok := pq.Dequeue(); ok {
			if prio < lastPriority {
				t.Errorf("Priority order violated: %d came after %d", prio, lastPriority)
			}
			lastPriority = prio
		}
	}
}

func TestPriorityQueueWithCustomType(t *testing.T) {
	type Task struct {
		Name     string
		Duration int
	}

	pq := NewPriorityQueue[Task]()

	tasks := []struct {
		task     Task
		priority int
	}{
		{Task{"Short", 5}, 3},
		{Task{"Medium", 10}, 2},
		{Task{"Long", 15}, 1},
	}

	for _, t := range tasks {
		pq.Enqueue(t.task, t.priority)
	}

	// Test order by priority
	expected := []int{1, 2, 3}
	for _, exp := range expected {
		if _, prio, ok := pq.Dequeue(); !ok || prio != exp {
			t.Errorf("Expected priority %d, got %d", exp, prio)
		}
	}
}

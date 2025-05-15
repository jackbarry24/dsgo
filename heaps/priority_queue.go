package heaps

import (
	"sync"
)

// PriorityQueueItem represents an item in the priority queue with a priority value
type PriorityQueueItem[T any] struct {
	Value    T
	Priority int
}

// PriorityQueue is a generic priority queue implementation
type PriorityQueue[T any] struct {
	items []PriorityQueueItem[T]
	less  func(a, b PriorityQueueItem[T]) bool
}

// SafePriorityQueue is a thread-safe version of PriorityQueue
type SafePriorityQueue[T any] struct {
	mu    sync.RWMutex
	inner *PriorityQueue[T]
}

// NewPriorityQueue creates a new priority queue
// By default, lower priority numbers are dequeued first (like a min heap)
func NewPriorityQueue[T any]() *PriorityQueue[T] {
	return &PriorityQueue[T]{
		items: []PriorityQueueItem[T]{},
		less: func(a, b PriorityQueueItem[T]) bool {
			return a.Priority < b.Priority
		},
	}
}

// NewMaxPriorityQueue creates a new priority queue where higher priority numbers are dequeued first
func NewMaxPriorityQueue[T any]() *PriorityQueue[T] {
	return &PriorityQueue[T]{
		items: []PriorityQueueItem[T]{},
		less: func(a, b PriorityQueueItem[T]) bool {
			return a.Priority > b.Priority
		},
	}
}

// NewSafePriorityQueue creates a new thread-safe priority queue
func NewSafePriorityQueue[T any]() *SafePriorityQueue[T] {
	return &SafePriorityQueue[T]{
		mu:    sync.RWMutex{},
		inner: NewPriorityQueue[T](),
	}
}

// NewSafeMaxPriorityQueue creates a new thread-safe priority queue where higher priority numbers are dequeued first
func NewSafeMaxPriorityQueue[T any]() *SafePriorityQueue[T] {
	return &SafePriorityQueue[T]{
		mu:    sync.RWMutex{},
		inner: NewMaxPriorityQueue[T](),
	}
}

// Enqueue adds an item to the priority queue with the given priority
func (pq *PriorityQueue[T]) Enqueue(value T, priority int) {
	pq.items = append(pq.items, PriorityQueueItem[T]{Value: value, Priority: priority})
	pq.up(len(pq.items) - 1)
}

// Dequeue removes and returns the item with the highest priority
func (pq *PriorityQueue[T]) Dequeue() (T, int, bool) {
	if len(pq.items) == 0 {
		var zero T
		return zero, 0, false
	}

	item := pq.items[0]
	last := len(pq.items) - 1
	pq.items[0] = pq.items[last]
	pq.items = pq.items[:last]

	if len(pq.items) > 0 {
		pq.down(0)
	}

	return item.Value, item.Priority, true
}

// Peek returns the highest priority item without removing it
func (pq *PriorityQueue[T]) Peek() (T, int, bool) {
	if len(pq.items) == 0 {
		var zero T
		return zero, 0, false
	}
	return pq.items[0].Value, pq.items[0].Priority, true
}

// Size returns the number of items in the queue
func (pq *PriorityQueue[T]) Size() int {
	return len(pq.items)
}

// IsEmpty returns true if the queue is empty
func (pq *PriorityQueue[T]) IsEmpty() bool {
	return len(pq.items) == 0
}

// up moves an element up the heap to maintain the heap property
func (pq *PriorityQueue[T]) up(i int) {
	for {
		parent := (i - 1) / 2
		if i == parent || !pq.less(pq.items[i], pq.items[parent]) {
			break
		}
		pq.items[i], pq.items[parent] = pq.items[parent], pq.items[i]
		i = parent
	}
}

// down moves an element down the heap to maintain the heap property
func (pq *PriorityQueue[T]) down(i int) {
	for {
		left := 2*i + 1
		if left >= len(pq.items) {
			break
		}

		smallest := left
		right := left + 1

		if right < len(pq.items) && pq.less(pq.items[right], pq.items[left]) {
			smallest = right
		}

		if !pq.less(pq.items[smallest], pq.items[i]) {
			break
		}

		pq.items[i], pq.items[smallest] = pq.items[smallest], pq.items[i]
		i = smallest
	}
}

// Thread-safe wrapper methods

func (pq *SafePriorityQueue[T]) Enqueue(value T, priority int) {
	pq.mu.Lock()
	defer pq.mu.Unlock()
	pq.inner.Enqueue(value, priority)
}

func (pq *SafePriorityQueue[T]) Dequeue() (T, int, bool) {
	pq.mu.Lock()
	defer pq.mu.Unlock()
	return pq.inner.Dequeue()
}

func (pq *SafePriorityQueue[T]) Peek() (T, int, bool) {
	pq.mu.RLock()
	defer pq.mu.RUnlock()
	return pq.inner.Peek()
}

func (pq *SafePriorityQueue[T]) Size() int {
	pq.mu.RLock()
	defer pq.mu.RUnlock()
	return pq.inner.Size()
}

func (pq *SafePriorityQueue[T]) IsEmpty() bool {
	pq.mu.RLock()
	defer pq.mu.RUnlock()
	return pq.inner.IsEmpty()
}

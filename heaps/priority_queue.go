package heaps

import (
	"sync"
)

type PriorityQueueItem[T any] struct {
	Value    T
	Priority int
}

type PriorityQueue[T any] struct {
	items      []PriorityQueueItem[T]
	less       func(a, b PriorityQueueItem[T]) bool
	threadSafe bool
	mu         sync.RWMutex
}

func NewPriorityQueue[T any](threadSafe ...bool) *PriorityQueue[T] {
	isThreadSafe := true
	if len(threadSafe) > 0 {
		isThreadSafe = threadSafe[0]
	}
	return &PriorityQueue[T]{
		items: []PriorityQueueItem[T]{},
		less: func(a, b PriorityQueueItem[T]) bool {
			return a.Priority < b.Priority
		},
		threadSafe: isThreadSafe,
	}
}

func (pq *PriorityQueue[T]) Enqueue(value T, priority int) {
	if pq.threadSafe {
		pq.mu.Lock()
		defer pq.mu.Unlock()
	}
	pq.items = append(pq.items, PriorityQueueItem[T]{Value: value, Priority: priority})
	pq.up(len(pq.items) - 1)
}

func (pq *PriorityQueue[T]) Dequeue() (T, int, bool) {
	if pq.threadSafe {
		pq.mu.Lock()
		defer pq.mu.Unlock()
	}
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

func (pq *PriorityQueue[T]) Peek() (T, int, bool) {
	if pq.threadSafe {
		pq.mu.RLock()
		defer pq.mu.RUnlock()
	}
	if len(pq.items) == 0 {
		var zero T
		return zero, 0, false
	}
	return pq.items[0].Value, pq.items[0].Priority, true
}

func (pq *PriorityQueue[T]) Size() int {
	if pq.threadSafe {
		pq.mu.RLock()
		defer pq.mu.RUnlock()
	}
	return len(pq.items)
}

func (pq *PriorityQueue[T]) IsEmpty() bool {
	if pq.threadSafe {
		pq.mu.RLock()
		defer pq.mu.RUnlock()
	}
	return len(pq.items) == 0
}

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

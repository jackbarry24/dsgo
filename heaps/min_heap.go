package heaps

import (
	"sync"
)

type MinHeap[T any] struct {
	items      []T
	less       func(a, b T) bool
	threadSafe bool
	mu         sync.RWMutex
}

func NewMinHeap[T any](less func(a, b T) bool, threadSafe ...bool) *MinHeap[T] {
	isThreadSafe := true
	if len(threadSafe) > 0 {
		isThreadSafe = threadSafe[0]
	}
	return &MinHeap[T]{
		items:      []T{},
		less:       less,
		threadSafe: isThreadSafe,
	}
}

func (h *MinHeap[T]) Push(item T) {
	if h.threadSafe {
		h.mu.Lock()
		defer h.mu.Unlock()
	}
	h.items = append(h.items, item)
	h.up(len(h.items) - 1)
}

func (h *MinHeap[T]) Pop() (T, bool) {
	if h.threadSafe {
		h.mu.Lock()
		defer h.mu.Unlock()
	}
	if len(h.items) == 0 {
		var zero T
		return zero, false
	}

	item := h.items[0]
	last := len(h.items) - 1
	h.items[0] = h.items[last]
	h.items = h.items[:last]

	if len(h.items) > 0 {
		h.down(0)
	}

	return item, true
}

func (h *MinHeap[T]) Peek() (T, bool) {
	if h.threadSafe {
		h.mu.RLock()
		defer h.mu.RUnlock()
	}
	if len(h.items) == 0 {
		var zero T
		return zero, false
	}
	return h.items[0], true
}

func (h *MinHeap[T]) Size() int {
	if h.threadSafe {
		h.mu.RLock()
		defer h.mu.RUnlock()
	}
	return len(h.items)
}

func (h *MinHeap[T]) IsEmpty() bool {
	if h.threadSafe {
		h.mu.RLock()
		defer h.mu.RUnlock()
	}
	return len(h.items) == 0
}

func (h *MinHeap[T]) up(i int) {
	for {
		parent := (i - 1) / 2
		if i == parent || !h.less(h.items[i], h.items[parent]) {
			break
		}
		h.items[i], h.items[parent] = h.items[parent], h.items[i]
		i = parent
	}
}

func (h *MinHeap[T]) down(i int) {
	for {
		left := 2*i + 1
		if left >= len(h.items) {
			break
		}

		smallest := left
		right := left + 1

		if right < len(h.items) && h.less(h.items[right], h.items[left]) {
			smallest = right
		}

		if !h.less(h.items[smallest], h.items[i]) {
			break
		}

		h.items[i], h.items[smallest] = h.items[smallest], h.items[i]
		i = smallest
	}
}

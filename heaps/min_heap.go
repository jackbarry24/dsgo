package heaps

import (
	"sync"
)

type MinHeap[T any] struct {
	items []T
	less  func(a, b T) bool
}

type SafeMinHeap[T any] struct {
	mu    sync.RWMutex
	inner *MinHeap[T]
}

func NewMinHeap[T any](less func(a, b T) bool) *MinHeap[T] {
	return &MinHeap[T]{
		items: []T{},
		less:  less,
	}
}

func NewSafeMinHeap[T any](less func(a, b T) bool) *SafeMinHeap[T] {
	return &SafeMinHeap[T]{
		mu:    sync.RWMutex{},
		inner: NewMinHeap(less),
	}
}

// Push adds an element to the heap
func (h *MinHeap[T]) Push(item T) {
	h.items = append(h.items, item)
	h.up(len(h.items) - 1)
}

// Pop removes and returns the minimum element from the heap
func (h *MinHeap[T]) Pop() (T, bool) {
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

// Peek returns the minimum element without removing it
func (h *MinHeap[T]) Peek() (T, bool) {
	if len(h.items) == 0 {
		var zero T
		return zero, false
	}
	return h.items[0], true
}

// Size returns the number of elements in the heap
func (h *MinHeap[T]) Size() int {
	return len(h.items)
}

// IsEmpty returns true if the heap is empty
func (h *MinHeap[T]) IsEmpty() bool {
	return len(h.items) == 0
}

// up moves an element up the heap to maintain the heap property
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

// down moves an element down the heap to maintain the heap property
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

// Thread-safe wrapper methods

func (h *SafeMinHeap[T]) Push(item T) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.inner.Push(item)
}

func (h *SafeMinHeap[T]) Pop() (T, bool) {
	h.mu.Lock()
	defer h.mu.Unlock()
	return h.inner.Pop()
}

func (h *SafeMinHeap[T]) Peek() (T, bool) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.inner.Peek()
}

func (h *SafeMinHeap[T]) Size() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.inner.Size()
}

func (h *SafeMinHeap[T]) IsEmpty() bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.inner.IsEmpty()
}

package sets

import "sync"

// Set is a basic set implementation that is not safe for concurrent use.
type Set[T comparable] struct {
	items map[T]struct{}
}

func NewSet[T comparable]() *Set[T] {
	return &Set[T]{
		items: make(map[T]struct{}),
	}
}

func (s *Set[T]) Add(item T) {
	s.items[item] = struct{}{}
}

func (s *Set[T]) Remove(item T) {
	delete(s.items, item)
}

func (s *Set[T]) Contains(item T) bool {
	_, exists := s.items[item]
	return exists
}

func (s *Set[T]) Size() int {
	return len(s.items)
}

func (s *Set[T]) IsEmpty() bool {
	return s.Size() == 0
}

func (s *Set[T]) Clear() {
	s.items = make(map[T]struct{})
}

func (s *Set[T]) Union(other *Set[T]) *Set[T] {
	result := NewSet[T]()
	for item := range s.items {
		result.Add(item)
	}
	for item := range other.items {
		result.Add(item)
	}
	return result
}

func (s *Set[T]) Intersection(other *Set[T]) *Set[T] {
	result := NewSet[T]()
	for item := range s.items {
		if other.Contains(item) {
			result.Add(item)
		}
	}
	return result
}

func (s *Set[T]) Difference(other *Set[T]) *Set[T] {
	result := NewSet[T]()
	for item := range s.items {
		if !other.Contains(item) {
			result.Add(item)
		}
	}
	return result
}

func (s *Set[T]) Items() []T {
	items := make([]T, 0, s.Size())
	for item := range s.items {
		items = append(items, item)
	}
	return items
}

// SafeSet is a thread-safe wrapper around Set.
type SafeSet[T comparable] struct {
	mu    sync.RWMutex
	items map[T]struct{}
}

func NewSafeSet[T comparable]() *SafeSet[T] {
	return &SafeSet[T]{
		items: make(map[T]struct{}),
	}
}

func (s *SafeSet[T]) Add(item T) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.items[item] = struct{}{}
}

func (s *SafeSet[T]) Remove(item T) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.items, item)
}

func (s *SafeSet[T]) Contains(item T) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, exists := s.items[item]
	return exists
}

func (s *SafeSet[T]) Size() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.items)
}

func (s *SafeSet[T]) IsEmpty() bool {
	return s.Size() == 0
}

func (s *SafeSet[T]) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.items = make(map[T]struct{})
}

func (s *SafeSet[T]) Union(other *SafeSet[T]) *SafeSet[T] {
	s.mu.RLock()
	other.mu.RLock()
	defer s.mu.RUnlock()
	defer other.mu.RUnlock()

	result := NewSafeSet[T]()
	for item := range s.items {
		result.Add(item)
	}
	for item := range other.items {
		result.Add(item)
	}
	return result
}

func (s *SafeSet[T]) Intersection(other *SafeSet[T]) *SafeSet[T] {
	s.mu.RLock()
	other.mu.RLock()
	defer s.mu.RUnlock()
	defer other.mu.RUnlock()

	result := NewSafeSet[T]()
	for item := range s.items {
		if _, exists := other.items[item]; exists {
			result.Add(item)
		}
	}
	return result
}

func (s *SafeSet[T]) Difference(other *SafeSet[T]) *SafeSet[T] {
	s.mu.RLock()
	other.mu.RLock()
	defer s.mu.RUnlock()
	defer other.mu.RUnlock()

	result := NewSafeSet[T]()
	for item := range s.items {
		if _, exists := other.items[item]; !exists {
			result.Add(item)
		}
	}
	return result
}

func (s *SafeSet[T]) Items() []T {
	s.mu.RLock()
	defer s.mu.RUnlock()

	items := make([]T, 0, len(s.items))
	for item := range s.items {
		items = append(items, item)
	}
	return items
}

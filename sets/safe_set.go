package sets

import "sync"

type SafeSet[T comparable] struct {
	mu    sync.RWMutex
	items map[T]struct{}
}

func NewSafe[T comparable]() *SafeSet[T] {
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

	result := NewSafe[T]()
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

	result := NewSafe[T]()
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

	result := NewSafe[T]()
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

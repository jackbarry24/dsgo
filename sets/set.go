package sets

import "sync"

type Set[T comparable] struct {
	items      map[T]struct{}
	threadSafe bool
	mu         sync.RWMutex
}

func NewSet[T comparable](threadSafe ...bool) *Set[T] {
	isThreadSafe := true
	if len(threadSafe) > 0 {
		isThreadSafe = threadSafe[0]
	}
	return &Set[T]{
		items:      make(map[T]struct{}),
		threadSafe: isThreadSafe,
	}
}

func (s *Set[T]) Add(item T) {
	if s.threadSafe {
		s.mu.Lock()
		defer s.mu.Unlock()
	}
	s.items[item] = struct{}{}
}

func (s *Set[T]) Remove(item T) {
	if s.threadSafe {
		s.mu.Lock()
		defer s.mu.Unlock()
	}
	delete(s.items, item)
}

func (s *Set[T]) Contains(item T) bool {
	if s.threadSafe {
		s.mu.RLock()
		defer s.mu.RUnlock()
	}
	_, exists := s.items[item]
	return exists
}

func (s *Set[T]) Size() int {
	if s.threadSafe {
		s.mu.RLock()
		defer s.mu.RUnlock()
	}
	return len(s.items)
}

func (s *Set[T]) IsEmpty() bool {
	return s.Size() == 0
}

func (s *Set[T]) Clear() {
	if s.threadSafe {
		s.mu.Lock()
		defer s.mu.Unlock()
	}
	s.items = make(map[T]struct{})
}

func (s *Set[T]) Union(other *Set[T]) *Set[T] {
	if s.threadSafe {
		s.mu.RLock()
		defer s.mu.RUnlock()
	}
	if other.threadSafe {
		other.mu.RLock()
		defer other.mu.RUnlock()
	}

	result := NewSet[T](s.threadSafe)
	for item := range s.items {
		result.Add(item)
	}
	for item := range other.items {
		result.Add(item)
	}
	return result
}

func (s *Set[T]) Intersection(other *Set[T]) *Set[T] {
	if s.threadSafe {
		s.mu.RLock()
		defer s.mu.RUnlock()
	}
	if other.threadSafe {
		other.mu.RLock()
		defer other.mu.RUnlock()
	}

	result := NewSet[T](s.threadSafe)
	for item := range s.items {
		if other.Contains(item) {
			result.Add(item)
		}
	}
	return result
}

func (s *Set[T]) Difference(other *Set[T]) *Set[T] {
	if s.threadSafe {
		s.mu.RLock()
		defer s.mu.RUnlock()
	}
	if other.threadSafe {
		other.mu.RLock()
		defer other.mu.RUnlock()
	}

	result := NewSet[T](s.threadSafe)
	for item := range s.items {
		if !other.Contains(item) {
			result.Add(item)
		}
	}
	return result
}

func (s *Set[T]) Items() []T {
	if s.threadSafe {
		s.mu.RLock()
		defer s.mu.RUnlock()
	}
	items := make([]T, 0, s.Size())
	for item := range s.items {
		items = append(items, item)
	}
	return items
}

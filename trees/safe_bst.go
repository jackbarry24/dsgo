package trees

import (
	"dsgo/utils"
	"sync"
)

type SafeBST[K utils.Ordered, V any] struct {
	mu    sync.RWMutex
	inner *BST[K, V]
}

func NewSafeBST[K utils.Ordered, V any]() *SafeBST[K, V] {
	return &SafeBST[K, V]{
		inner: NewBST[K, V](),
	}
}

func (s *SafeBST[K, V]) Insert(key K, value V) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.inner.Insert(key, value)
}

func (s *SafeBST[K, V]) Search(key K) (V, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.inner.Search(key)
}

func (s *SafeBST[K, V]) Delete(key K) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.inner.Delete(key)
}

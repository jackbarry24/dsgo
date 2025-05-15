package linkedlist

import (
	"errors"
	"sync"
)

type Node[T comparable] struct {
	value T
	next  *Node[T]
}

type SingleLinkedList[T comparable] struct {
	head       *Node[T]
	tail       *Node[T]
	len        int
	threadSafe bool
	mu         sync.RWMutex
}

func NewSingleLinkedList[T comparable](threadSafe ...bool) *SingleLinkedList[T] {
	isThreadSafe := true
	if len(threadSafe) > 0 {
		isThreadSafe = threadSafe[0]
	}
	return &SingleLinkedList[T]{
		head:       nil,
		tail:       nil,
		threadSafe: isThreadSafe,
	}
}

func (l *SingleLinkedList[T]) PushBack(value T) {
	if l.threadSafe {
		l.mu.Lock()
		defer l.mu.Unlock()
	}

	newNode := &Node[T]{value: value}
	if l.head == nil {
		l.head = newNode
		l.tail = newNode
	} else {
		l.tail.next = newNode
		l.tail = newNode
	}
	l.len++
}

func (l *SingleLinkedList[T]) PushFront(value T) {
	if l.threadSafe {
		l.mu.Lock()
		defer l.mu.Unlock()
	}

	newNode := &Node[T]{value: value}
	if l.head == nil {
		l.head = newNode
		l.tail = newNode
	} else {
		newNode.next = l.head
		l.head = newNode
	}
	l.len++
}

func (l *SingleLinkedList[T]) Remove(value T) error {
	if l.threadSafe {
		l.mu.Lock()
		defer l.mu.Unlock()
	}

	if l.head == nil {
		return ErrEmptyList
	}

	// Special case: removing head
	if l.head.value == value {
		l.head = l.head.next
		if l.head == nil {
			l.tail = nil
		}
		l.len--
		return nil
	}

	// Search for the node to remove
	current := l.head
	for current.next != nil {
		if current.next.value == value {
			current.next = current.next.next
			if current.next == nil {
				l.tail = current
			}
			l.len--
			return nil
		}
		current = current.next
	}

	return ErrNotFound
}

func (l *SingleLinkedList[T]) Contains(value T) bool {
	if l.threadSafe {
		l.mu.RLock()
		defer l.mu.RUnlock()
	}

	current := l.head
	for current != nil {
		if current.value == value {
			return true
		}
		current = current.next
	}
	return false
}

func (l *SingleLinkedList[T]) At(index int) (T, error) {
	if l.threadSafe {
		l.mu.RLock()
		defer l.mu.RUnlock()
	}

	if index < 0 || index >= l.len {
		var zero T
		return zero, errors.New("index out of bounds")
	}

	current := l.head
	for i := 0; i < index; i++ {
		current = current.next
	}
	return current.value, nil
}

func (l *SingleLinkedList[T]) Clear() {
	if l.threadSafe {
		l.mu.Lock()
		defer l.mu.Unlock()
	}

	l.head = nil
	l.tail = nil
	l.len = 0
}

func (l *SingleLinkedList[T]) ForEach(f func(T)) {
	if l.threadSafe {
		l.mu.RLock()
		defer l.mu.RUnlock()
	}

	current := l.head
	for current != nil {
		f(current.value)
		current = current.next
	}
}

func (l *SingleLinkedList[T]) Back() (*Node[T], error) {
	if l.threadSafe {
		l.mu.RLock()
		defer l.mu.RUnlock()
	}

	if l.tail == nil {
		return nil, ErrEmptyList
	}
	return l.tail, nil
}

func (l *SingleLinkedList[T]) Front() (*Node[T], error) {
	if l.threadSafe {
		l.mu.RLock()
		defer l.mu.RUnlock()
	}

	if l.head == nil {
		return nil, ErrEmptyList
	}
	return l.head, nil
}

func (l *SingleLinkedList[T]) Len() int {
	if l.threadSafe {
		l.mu.RLock()
		defer l.mu.RUnlock()
	}
	return l.len
}

package linkedlist

import (
	"errors"
	"sync"
)

type DNode[T comparable] struct {
	value T
	prev  *DNode[T]
	next  *DNode[T]
}

// GetValue returns the value stored in the node
func (n *DNode[T]) GetValue() T {
	return n.value
}

type DoubleLinkedList[T comparable] struct {
	head       *DNode[T]
	tail       *DNode[T]
	len        int
	threadSafe bool
	mu         sync.RWMutex
}

func NewDoubleLinkedList[T comparable](threadSafe ...bool) *DoubleLinkedList[T] {
	isThreadSafe := true
	if len(threadSafe) > 0 {
		isThreadSafe = threadSafe[0]
	}
	return &DoubleLinkedList[T]{
		head:       nil,
		tail:       nil,
		threadSafe: isThreadSafe,
	}
}

func (l *DoubleLinkedList[T]) PushBack(value T) {
	if l.threadSafe {
		l.mu.Lock()
		defer l.mu.Unlock()
	}

	newNode := &DNode[T]{value: value}
	if l.head == nil {
		l.head = newNode
		l.tail = newNode
	} else {
		newNode.prev = l.tail
		l.tail.next = newNode
		l.tail = newNode
	}
	l.len++
}

func (l *DoubleLinkedList[T]) PushFront(value T) {
	if l.threadSafe {
		l.mu.Lock()
		defer l.mu.Unlock()
	}

	newNode := &DNode[T]{value: value}
	if l.head == nil {
		l.head = newNode
		l.tail = newNode
	} else {
		newNode.next = l.head
		l.head.prev = newNode
		l.head = newNode
	}
	l.len++
}

func (l *DoubleLinkedList[T]) Remove(value T) error {
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
		} else {
			l.head.prev = nil
		}
		l.len--
		return nil
	}

	// Special case: removing tail
	if l.tail.value == value {
		l.tail = l.tail.prev
		l.tail.next = nil
		l.len--
		return nil
	}

	// Search for the node to remove
	current := l.head.next
	for current != nil && current != l.tail {
		if current.value == value {
			current.prev.next = current.next
			current.next.prev = current.prev
			l.len--
			return nil
		}
		current = current.next
	}

	return ErrNotFound
}

func (l *DoubleLinkedList[T]) Contains(value T) bool {
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

func (l *DoubleLinkedList[T]) At(index int) (T, error) {
	if l.threadSafe {
		l.mu.RLock()
		defer l.mu.RUnlock()
	}

	if index < 0 || index >= l.len {
		var zero T
		return zero, errors.New("index out of bounds")
	}

	// Optimize traversal based on index position
	if index < l.len/2 {
		current := l.head
		for i := 0; i < index; i++ {
			current = current.next
		}
		return current.value, nil
	} else {
		current := l.tail
		for i := l.len - 1; i > index; i-- {
			current = current.prev
		}
		return current.value, nil
	}
}

func (l *DoubleLinkedList[T]) Clear() {
	if l.threadSafe {
		l.mu.Lock()
		defer l.mu.Unlock()
	}

	l.head = nil
	l.tail = nil
	l.len = 0
}

func (l *DoubleLinkedList[T]) ForEach(f func(T)) {
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

func (l *DoubleLinkedList[T]) ForEachReverse(f func(T)) {
	if l.threadSafe {
		l.mu.RLock()
		defer l.mu.RUnlock()
	}

	current := l.tail
	for current != nil {
		f(current.value)
		current = current.prev
	}
}

func (l *DoubleLinkedList[T]) Back() (*DNode[T], error) {
	if l.threadSafe {
		l.mu.RLock()
		defer l.mu.RUnlock()
	}

	if l.tail == nil {
		return nil, ErrEmptyList
	}
	return l.tail, nil
}

func (l *DoubleLinkedList[T]) Front() (*DNode[T], error) {
	if l.threadSafe {
		l.mu.RLock()
		defer l.mu.RUnlock()
	}

	if l.head == nil {
		return nil, ErrEmptyList
	}
	return l.head, nil
}

func (l *DoubleLinkedList[T]) Len() int {
	if l.threadSafe {
		l.mu.RLock()
		defer l.mu.RUnlock()
	}
	return l.len
}

// InsertAfter inserts a new node with the given value after the node with the target value
func (l *DoubleLinkedList[T]) InsertAfter(target, value T) error {
	if l.threadSafe {
		l.mu.Lock()
		defer l.mu.Unlock()
	}

	if l.head == nil {
		return ErrEmptyList
	}

	current := l.head
	for current != nil {
		if current.value == target {
			newNode := &DNode[T]{value: value}
			newNode.next = current.next
			newNode.prev = current
			if current.next != nil {
				current.next.prev = newNode
			} else {
				l.tail = newNode
			}
			current.next = newNode
			l.len++
			return nil
		}
		current = current.next
	}

	return ErrNotFound
}

// InsertBefore inserts a new node with the given value before the node with the target value
func (l *DoubleLinkedList[T]) InsertBefore(target, value T) error {
	if l.threadSafe {
		l.mu.Lock()
		defer l.mu.Unlock()
	}

	if l.head == nil {
		return ErrEmptyList
	}

	// Special case: inserting before head
	if l.head.value == target {
		newNode := &DNode[T]{value: value}
		newNode.next = l.head
		l.head.prev = newNode
		l.head = newNode
		l.len++
		return nil
	}

	current := l.head.next
	for current != nil {
		if current.value == target {
			newNode := &DNode[T]{value: value}
			newNode.next = current
			newNode.prev = current.prev
			current.prev.next = newNode
			current.prev = newNode
			l.len++
			return nil
		}
		current = current.next
	}

	return ErrNotFound
}

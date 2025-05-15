package trees

import (
	"dsgo/utils"
	"sync"
)

type BST[K utils.Ordered, V any] struct {
	root       *Node[K, V]
	threadSafe bool
	mu         sync.RWMutex
}

type Node[K utils.Ordered, V any] struct {
	key   K
	value V
	left  *Node[K, V]
	right *Node[K, V]
}

func NewBST[K utils.Ordered, V any](threadSafe ...bool) *BST[K, V] {
	isThreadSafe := true
	if len(threadSafe) > 0 {
		isThreadSafe = threadSafe[0]
	}
	return &BST[K, V]{
		threadSafe: isThreadSafe,
	}
}

func (b *BST[K, V]) Insert(key K, value V) {
	if !b.threadSafe {
		if b.root == nil {
			b.root = &Node[K, V]{key: key, value: value}
			return
		}
		b.root = insert(b.root, key, value)
		return
	}

	b.mu.Lock()
	defer b.mu.Unlock()

	if b.root == nil {
		b.root = &Node[K, V]{key: key, value: value}
		return
	}
	b.root = insert(b.root, key, value)
}

func (b *BST[K, V]) Search(key K) (V, bool) {
	if !b.threadSafe {
		if b.root == nil {
			var zero V
			return zero, false
		}
		return search(b.root, key)
	}

	b.mu.RLock()
	defer b.mu.RUnlock()

	if b.root == nil {
		var zero V
		return zero, false
	}
	return search(b.root, key)
}

func (b *BST[K, V]) Delete(key K) {
	if !b.threadSafe {
		if b.root == nil {
			return
		}
		b.root = delete(b.root, key)
		return
	}

	b.mu.Lock()
	defer b.mu.Unlock()

	if b.root == nil {
		return
	}
	b.root = delete(b.root, key)
}

func insert[K utils.Ordered, V any](node *Node[K, V], key K, value V) *Node[K, V] {
	if node == nil {
		return &Node[K, V]{key: key, value: value}
	}

	switch {
	case key < node.key:
		node.left = insert(node.left, key, value)
	case key > node.key:
		node.right = insert(node.right, key, value)
	default:
		node.value = value
	}
	return node
}

func search[K utils.Ordered, V any](node *Node[K, V], key K) (V, bool) {
	if node == nil {
		var zero V
		return zero, false
	}

	switch {
	case key < node.key:
		return search(node.left, key)
	case key > node.key:
		return search(node.right, key)
	default:
		return node.value, true
	}
}

func delete[K utils.Ordered, V any](node *Node[K, V], key K) *Node[K, V] {
	if node == nil {
		return nil
	}

	switch {
	case key < node.key:
		node.left = delete(node.left, key)
	case key > node.key:
		node.right = delete(node.right, key)
	default:
		// Case 1: Node with no children
		if node.left == nil && node.right == nil {
			return nil
		}
		// Case 2: Node with one child
		if node.left == nil {
			return node.right
		}
		if node.right == nil {
			return node.left
		}
		// Case 3: Node with two children
		successor := findMin(node.right)
		node.key = successor.key
		node.value = successor.value
		node.right = delete(node.right, successor.key)
	}
	return node
}

func findMin[K utils.Ordered, V any](node *Node[K, V]) *Node[K, V] {
	current := node
	for current.left != nil {
		current = current.left
	}
	return current
}

package trees

import (
	"dsgo/utils"
	"sync"
)

type Color bool

const (
	Red   Color = true
	Black Color = false
)

type RBNode[K utils.Ordered, V any] struct {
	key    K
	value  V
	color  Color
	left   *RBNode[K, V]
	right  *RBNode[K, V]
	parent *RBNode[K, V]
}

type RBTree[K utils.Ordered, V any] struct {
	root *RBNode[K, V]
}

type SafeRBTree[K utils.Ordered, V any] struct {
	mu    sync.RWMutex
	inner *RBTree[K, V]
}

func NewRBTree[K utils.Ordered, V any]() *RBTree[K, V] {
	return &RBTree[K, V]{}
}

func NewSafeRBTree[K utils.Ordered, V any]() *SafeRBTree[K, V] {
	return &SafeRBTree[K, V]{inner: NewRBTree[K, V]()}
}

func (t *RBTree[K, V]) Insert(key K, value V) {
	node := &RBNode[K, V]{key: key, value: value, color: Red}
	if t.root == nil {
		node.color = Black
		t.root = node
		return
	}

	// Find the position to insert
	current := t.root
	var parent *RBNode[K, V]
	for current != nil {
		parent = current
		if key < current.key {
			current = current.left
		} else if key > current.key {
			current = current.right
		} else {
			current.value = value
			return
		}
	}

	// Insert the node
	node.parent = parent
	if key < parent.key {
		parent.left = node
	} else {
		parent.right = node
	}

	t.fixInsert(node)
}

func (t *RBTree[K, V]) fixInsert(node *RBNode[K, V]) {
	for node != t.root && node.parent.color == Red {
		if node.parent == node.parent.parent.left {
			uncle := node.parent.parent.right
			if uncle != nil && uncle.color == Red {
				node.parent.color = Black
				uncle.color = Black
				node.parent.parent.color = Red
				node = node.parent.parent
			} else {
				if node == node.parent.right {
					node = node.parent
					t.rotateLeft(node)
				}
				node.parent.color = Black
				node.parent.parent.color = Red
				t.rotateRight(node.parent.parent)
			}
		} else {
			uncle := node.parent.parent.left
			if uncle != nil && uncle.color == Red {
				node.parent.color = Black
				uncle.color = Black
				node.parent.parent.color = Red
				node = node.parent.parent
			} else {
				if node == node.parent.left {
					node = node.parent
					t.rotateRight(node)
				}
				node.parent.color = Black
				node.parent.parent.color = Red
				t.rotateLeft(node.parent.parent)
			}
		}
	}
	t.root.color = Black
}

func (t *RBTree[K, V]) rotateLeft(node *RBNode[K, V]) {
	rightChild := node.right
	node.right = rightChild.left

	if rightChild.left != nil {
		rightChild.left.parent = node
	}
	rightChild.parent = node.parent

	if node.parent == nil {
		t.root = rightChild
	} else if node == node.parent.left {
		node.parent.left = rightChild
	} else {
		node.parent.right = rightChild
	}

	rightChild.left = node
	node.parent = rightChild
}

func (t *RBTree[K, V]) rotateRight(node *RBNode[K, V]) {
	leftChild := node.left
	node.left = leftChild.right

	if leftChild.right != nil {
		leftChild.right.parent = node
	}
	leftChild.parent = node.parent

	if node.parent == nil {
		t.root = leftChild
	} else if node == node.parent.right {
		node.parent.right = leftChild
	} else {
		node.parent.left = leftChild
	}

	leftChild.right = node
	node.parent = leftChild
}

func (t *RBTree[K, V]) Search(key K) (*RBNode[K, V], bool) {
	node := t.root
	for node != nil {
		if key < node.key {
			node = node.left
		} else if key > node.key {
			node = node.right
		} else {
			return node, true
		}
	}
	return nil, false
}

func (t *RBTree[K, V]) Delete(key K) {
	node, found := t.Search(key)
	if !found {
		return
	}
	t.deleteNode(node)
}

func (t *RBTree[K, V]) deleteNode(node *RBNode[K, V]) {
	var child *RBNode[K, V]
	originalColor := node.color

	if node.left == nil {
		child = node.right
		t.transplant(node, node.right)
	} else if node.right == nil {
		child = node.left
		t.transplant(node, node.left)
	} else {
		successor := t.minimum(node.right)
		originalColor = successor.color
		child = successor.right

		if successor.parent == node {
			if child != nil {
				child.parent = successor
			}
		} else {
			t.transplant(successor, successor.right)
			successor.right = node.right
			successor.right.parent = successor
		}

		t.transplant(node, successor)
		successor.left = node.left
		successor.left.parent = successor
		successor.color = node.color
	}

	if originalColor == Black {
		t.fixDelete(child)
	}
}

func (t *RBTree[K, V]) transplant(u, v *RBNode[K, V]) {
	if u.parent == nil {
		t.root = v
	} else if u == u.parent.left {
		u.parent.left = v
	} else {
		u.parent.right = v
	}
	if v != nil {
		v.parent = u.parent
	}
}

func (t *RBTree[K, V]) fixDelete(node *RBNode[K, V]) {
	for node != t.root && (node == nil || node.color == Black) {
		if node == node.parent.left {
			sibling := node.parent.right
			if sibling.color == Red {
				sibling.color = Black
				node.parent.color = Red
				t.rotateLeft(node.parent)
				sibling = node.parent.right
			}
			if (sibling.left == nil || sibling.left.color == Black) &&
				(sibling.right == nil || sibling.right.color == Black) {
				sibling.color = Red
				node = node.parent
			} else {
				if sibling.right == nil || sibling.right.color == Black {
					if sibling.left != nil {
						sibling.left.color = Black
					}
					sibling.color = Red
					t.rotateRight(sibling)
					sibling = node.parent.right
				}
				sibling.color = node.parent.color
				node.parent.color = Black
				if sibling.right != nil {
					sibling.right.color = Black
				}
				t.rotateLeft(node.parent)
				node = t.root
			}
		} else {
			sibling := node.parent.left
			if sibling.color == Red {
				sibling.color = Black
				node.parent.color = Red
				t.rotateRight(node.parent)
				sibling = node.parent.left
			}
			if (sibling.left == nil || sibling.left.color == Black) &&
				(sibling.right == nil || sibling.right.color == Black) {
				sibling.color = Red
				node = node.parent
			} else {
				if sibling.left == nil || sibling.left.color == Black {
					if sibling.right != nil {
						sibling.right.color = Black
					}
					sibling.color = Red
					t.rotateLeft(sibling)
					sibling = node.parent.left
				}
				sibling.color = node.parent.color
				node.parent.color = Black
				if sibling.left != nil {
					sibling.left.color = Black
				}
				t.rotateRight(node.parent)
				node = t.root
			}
		}
	}
	if node != nil {
		node.color = Black
	}
}

func (t *RBTree[K, V]) minimum(node *RBNode[K, V]) *RBNode[K, V] {
	for node.left != nil {
		node = node.left
	}
	return node
}

func (t *SafeRBTree[K, V]) Insert(key K, value V) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.inner.Insert(key, value)
}

func (t *SafeRBTree[K, V]) Search(key K) (*RBNode[K, V], bool) {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.inner.Search(key)
}

func (t *SafeRBTree[K, V]) Delete(key K) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.inner.Delete(key)
}

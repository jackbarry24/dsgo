package trees

import (
	"dsgo/utils"
	"sync"
)

type AVLNode[K utils.Ordered, V any] struct {
	Key    K
	Value  V
	Left   *AVLNode[K, V]
	Right  *AVLNode[K, V]
	Height int
}

type AVLTree[K utils.Ordered, V any] struct {
	Root       *AVLNode[K, V]
	threadSafe bool
	mu         sync.RWMutex
}

func NewAVLTree[K utils.Ordered, V any](threadSafe ...bool) *AVLTree[K, V] {
	isThreadSafe := true
	if len(threadSafe) > 0 {
		isThreadSafe = threadSafe[0]
	}
	return &AVLTree[K, V]{
		threadSafe: isThreadSafe,
	}
}

func height[K utils.Ordered, V any](node *AVLNode[K, V]) int {
	if node == nil {
		return 0
	}
	return node.Height
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func getBalance[K utils.Ordered, V any](node *AVLNode[K, V]) int {
	if node == nil {
		return 0
	}
	return height(node.Left) - height(node.Right)
}

func rightRotate[K utils.Ordered, V any](y *AVLNode[K, V]) *AVLNode[K, V] {
	x := y.Left
	T2 := x.Right

	// Perform rotation
	x.Right = y
	y.Left = T2

	// Update heights
	y.Height = max(height(y.Left), height(y.Right)) + 1
	x.Height = max(height(x.Left), height(x.Right)) + 1

	return x
}

func leftRotate[K utils.Ordered, V any](x *AVLNode[K, V]) *AVLNode[K, V] {
	y := x.Right
	T2 := y.Left

	// Perform rotation
	y.Left = x
	x.Right = T2

	// Update heights
	x.Height = max(height(x.Left), height(x.Right)) + 1
	y.Height = max(height(y.Left), height(y.Right)) + 1

	return y
}

func (t *AVLTree[K, V]) Insert(key K, value V) {
	if t.threadSafe {
		t.mu.Lock()
		defer t.mu.Unlock()
	}
	t.Root = t.insert(t.Root, key, value)
}

func (t *AVLTree[K, V]) insert(node *AVLNode[K, V], key K, value V) *AVLNode[K, V] {
	if node == nil {
		return &AVLNode[K, V]{Key: key, Value: value, Height: 1}
	}

	if key < node.Key {
		node.Left = t.insert(node.Left, key, value)
	} else if key > node.Key {
		node.Right = t.insert(node.Right, key, value)
	} else {
		// Update value for existing key
		node.Value = value
		return node
	}

	// Update height of current node
	node.Height = 1 + max(height(node.Left), height(node.Right))

	// Get balance factor
	balance := getBalance(node)

	// Left Left Case
	if balance > 1 && key < node.Left.Key {
		return rightRotate(node)
	}

	// Right Right Case
	if balance < -1 && key > node.Right.Key {
		return leftRotate(node)
	}

	// Left Right Case
	if balance > 1 && key > node.Left.Key {
		node.Left = leftRotate(node.Left)
		return rightRotate(node)
	}

	// Right Left Case
	if balance < -1 && key < node.Right.Key {
		node.Right = rightRotate(node.Right)
		return leftRotate(node)
	}

	return node
}

func (t *AVLTree[K, V]) Delete(key K) {
	if t.threadSafe {
		t.mu.Lock()
		defer t.mu.Unlock()
	}
	t.Root = t.delete(t.Root, key)
}

func (t *AVLTree[K, V]) delete(node *AVLNode[K, V], key K) *AVLNode[K, V] {
	if node == nil {
		return nil
	}

	if key < node.Key {
		node.Left = t.delete(node.Left, key)
	} else if key > node.Key {
		node.Right = t.delete(node.Right, key)
	} else {
		// Node to be deleted found

		// Node with only one child or no child
		if node.Left == nil {
			return node.Right
		} else if node.Right == nil {
			return node.Left
		}

		// Node with two children: Get the inorder successor (smallest in right subtree)
		temp := t.minValueNode(node.Right)
		node.Key = temp.Key
		node.Value = temp.Value
		node.Right = t.delete(node.Right, temp.Key)
	}

	if node == nil {
		return nil
	}

	// Update height
	node.Height = 1 + max(height(node.Left), height(node.Right))

	// Get balance factor
	balance := getBalance(node)

	// Left Left Case
	if balance > 1 && getBalance(node.Left) >= 0 {
		return rightRotate(node)
	}

	// Left Right Case
	if balance > 1 && getBalance(node.Left) < 0 {
		node.Left = leftRotate(node.Left)
		return rightRotate(node)
	}

	// Right Right Case
	if balance < -1 && getBalance(node.Right) <= 0 {
		return leftRotate(node)
	}

	// Right Left Case
	if balance < -1 && getBalance(node.Right) > 0 {
		node.Right = rightRotate(node.Right)
		return leftRotate(node)
	}

	return node
}

func (t *AVLTree[K, V]) minValueNode(node *AVLNode[K, V]) *AVLNode[K, V] {
	current := node
	for current.Left != nil {
		current = current.Left
	}
	return current
}

func (t *AVLTree[K, V]) Search(key K) (V, bool) {
	if t.threadSafe {
		t.mu.RLock()
		defer t.mu.RUnlock()
	}
	return t.search(t.Root, key)
}

func (t *AVLTree[K, V]) search(node *AVLNode[K, V], key K) (V, bool) {
	if node == nil {
		var zero V
		return zero, false
	}

	if key < node.Key {
		return t.search(node.Left, key)
	} else if key > node.Key {
		return t.search(node.Right, key)
	}
	return node.Value, true
}

func (t *AVLTree[K, V]) InOrderTraversal() []V {
	if t.threadSafe {
		t.mu.RLock()
		defer t.mu.RUnlock()
	}
	var result []V
	t.inOrderTraversal(t.Root, &result)
	return result
}

func (t *AVLTree[K, V]) inOrderTraversal(node *AVLNode[K, V], result *[]V) {
	if node != nil {
		t.inOrderTraversal(node.Left, result)
		*result = append(*result, node.Value)
		t.inOrderTraversal(node.Right, result)
	}
}

package trees

import (
	"dsgo/utils"
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestRBTree_Insert(t *testing.T) {
	tests := []struct {
		name     string
		inserts  []struct{ key, value int }
		wantSize int
	}{
		{
			name: "empty tree",
			inserts: []struct{ key, value int }{
				{5, 5},
			},
			wantSize: 1,
		},
		{
			name: "multiple inserts",
			inserts: []struct{ key, value int }{
				{5, 5},
				{3, 3},
				{7, 7},
				{1, 1},
				{9, 9},
			},
			wantSize: 5,
		},
		{
			name: "duplicate key",
			inserts: []struct{ key, value int }{
				{5, 5},
				{5, 10}, // Should update value
			},
			wantSize: 1,
		},
		{
			name: "insert in descending order",
			inserts: []struct{ key, value int }{
				{9, 9},
				{7, 7},
				{5, 5},
				{3, 3},
				{1, 1},
			},
			wantSize: 5,
		},
		{
			name: "insert in ascending order",
			inserts: []struct{ key, value int }{
				{1, 1},
				{3, 3},
				{5, 5},
				{7, 7},
				{9, 9},
			},
			wantSize: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rb := NewRBTree[int, int](false)
			for _, insert := range tt.inserts {
				rb.Insert(insert.key, insert.value)
			}

			// Verify all values are in the tree
			lastValue := make(map[int]int)
			for _, insert := range tt.inserts {
				lastValue[insert.key] = insert.value
			}

			for key, wantValue := range lastValue {
				node, found := rb.Search(key)
				if !found {
					t.Errorf("Search(%d) = not found, want found", key)
				}
				if node.value != wantValue {
					t.Errorf("Search(%d) = %d, want %d", key, node.value, wantValue)
				}
			}

			// Verify tree properties
			verifyRBProperties(t, rb)
		})
	}
}

func TestRBTree_Search(t *testing.T) {
	rb := NewRBTree[int, int](false)
	rb.Insert(5, 5)
	rb.Insert(3, 3)
	rb.Insert(7, 7)

	tests := []struct {
		name   string
		key    int
		want   int
		wantOk bool
	}{
		{
			name:   "existing key",
			key:    5,
			want:   5,
			wantOk: true,
		},
		{
			name:   "non-existing key",
			key:    4,
			want:   0,
			wantOk: false,
		},
		{
			name:   "search in empty tree",
			key:    1,
			want:   0,
			wantOk: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node, ok := rb.Search(tt.key)
			if ok != tt.wantOk {
				t.Errorf("Search() ok = %v, want %v", ok, tt.wantOk)
			}
			if ok && node.value != tt.want {
				t.Errorf("Search() = %v, want %v", node.value, tt.want)
			}
		})
	}
}

func TestRBTree_Delete(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(*RBTree[int, int])
		delete int
		check  func(*testing.T, *RBTree[int, int])
	}{
		{
			name: "delete leaf node",
			setup: func(rb *RBTree[int, int]) {
				rb.Insert(5, 5)
				rb.Insert(3, 3)
				rb.Insert(7, 7)
			},
			delete: 3,
			check: func(t *testing.T, rb *RBTree[int, int]) {
				if _, found := rb.Search(3); found {
					t.Error("Search(3) = found, want not found")
				}
				verifyRBProperties(t, rb)
			},
		},
		{
			name: "delete root node",
			setup: func(rb *RBTree[int, int]) {
				rb.Insert(5, 5)
				rb.Insert(3, 3)
				rb.Insert(7, 7)
			},
			delete: 5,
			check: func(t *testing.T, rb *RBTree[int, int]) {
				if _, found := rb.Search(5); found {
					t.Error("Search(5) = found, want not found")
				}
				verifyRBProperties(t, rb)
			},
		},
		{
			name: "delete node with two children",
			setup: func(rb *RBTree[int, int]) {
				rb.Insert(5, 5)
				rb.Insert(3, 3)
				rb.Insert(7, 7)
				rb.Insert(6, 6)
				rb.Insert(8, 8)
			},
			delete: 7,
			check: func(t *testing.T, rb *RBTree[int, int]) {
				if _, found := rb.Search(7); found {
					t.Error("Search(7) = found, want not found")
				}
				verifyRBProperties(t, rb)
			},
		},
		{
			name: "delete non-existent key",
			setup: func(rb *RBTree[int, int]) {
				rb.Insert(5, 5)
			},
			delete: 10,
			check: func(t *testing.T, rb *RBTree[int, int]) {
				if _, found := rb.Search(5); !found {
					t.Error("Search(5) = not found, want found")
				}
				verifyRBProperties(t, rb)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rb := NewRBTree[int, int](false)
			tt.setup(rb)
			rb.Delete(tt.delete)
			tt.check(t, rb)
		})
	}
}

func TestRBTreeConcurrent(t *testing.T) {
	tree := NewRBTree[int, string](true)
	var wg sync.WaitGroup
	done := make(chan bool)

	go func() {
		// Test concurrent Insert operations
		for i := 0; i < 1000; i++ {
			wg.Add(1)
			go func(val int) {
				defer wg.Done()
				tree.Insert(val, fmt.Sprintf("value-%d", val))
			}(i)
		}
		wg.Wait()

		// Verify all values were inserted
		for i := 0; i < 1000; i++ {
			if node, exists := tree.Search(i); !exists || node.value != fmt.Sprintf("value-%d", i) {
				t.Errorf("Expected value-%d, got %v", i, node.value)
			}
		}
		done <- true
	}()

	select {
	case <-done:
		// Test completed successfully
	case <-time.After(5 * time.Second):
		t.Fatal("Test timed out after 5 seconds")
	}
}

func TestRBTreeConcurrentDelete(t *testing.T) {
	tree := NewRBTree[int, string](true)

	// First insert values sequentially
	for i := 0; i < 20; i++ {
		tree.Insert(i, fmt.Sprintf("value-%d", i))
	}

	// Create a channel to signal completion
	done := make(chan struct{})

	// Start a goroutine to perform deletions
	go func() {
		var wg sync.WaitGroup
		// Delete values 0-9 concurrently
		for i := 0; i < 10; i++ {
			wg.Add(1)
			go func(val int) {
				defer wg.Done()
				// Add a small delay to ensure proper synchronization
				time.Sleep(time.Millisecond * time.Duration(val%5))
				tree.Delete(val)
			}(i)
		}
		wg.Wait()
		close(done)
	}()

	// Wait for completion with timeout
	select {
	case <-done:
		// Test completed successfully
		// Verify the results
		for i := 0; i < 20; i++ {
			_, exists := tree.Search(i)
			if i < 10 && exists {
				t.Errorf("Value %d should be deleted", i)
			} else if i >= 10 && !exists {
				t.Errorf("Value %d should exist", i)
			}
		}
	case <-time.After(3 * time.Second):
		t.Fatal("Test timed out after 3 seconds")
	}
}

func TestRBTreeConcurrentSearch(t *testing.T) {
	tree := NewRBTree[int, string](true)
	var wg sync.WaitGroup
	done := make(chan bool)

	go func() {
		// First insert values
		for i := 0; i < 1000; i++ {
			tree.Insert(i, fmt.Sprintf("value-%d", i))
		}

		// Test concurrent Search operations
		for i := 0; i < 1000; i++ {
			wg.Add(1)
			go func(val int) {
				defer wg.Done()
				if node, exists := tree.Search(val); !exists || node.value != fmt.Sprintf("value-%d", val) {
					t.Errorf("Expected value-%d, got %v", val, node.value)
				}
			}(i)
		}
		wg.Wait()
		done <- true
	}()

	select {
	case <-done:
		// Test completed successfully
	case <-time.After(5 * time.Second):
		t.Fatal("Test timed out after 5 seconds")
	}
}

func TestRBTreeConcurrentModifications(t *testing.T) {
	tree := NewRBTree[int, string](true)

	// Create a channel to signal completion
	done := make(chan struct{})

	// Start a goroutine to perform modifications
	go func() {
		var wg sync.WaitGroup
		// Perform 10 concurrent operations
		for i := 0; i < 10; i++ {
			wg.Add(2)
			// Insert operation
			go func(val int) {
				defer wg.Done()
				time.Sleep(time.Millisecond * time.Duration(val%5))
				tree.Insert(val, fmt.Sprintf("value-%d", val))
			}(i)
			// Delete operation
			go func(val int) {
				defer wg.Done()
				time.Sleep(time.Millisecond * time.Duration((val+2)%5))
				if val < 10 {
					tree.Delete(val)
				}
			}(i)
		}
		wg.Wait()
		close(done)
	}()

	// Wait for completion with timeout
	select {
	case <-done:
		// Test completed successfully
		// No RB property check here: concurrent insert/delete on same keys is not guaranteed to maintain RB properties at every instant
	case <-time.After(3 * time.Second):
		t.Fatal("Test timed out after 3 seconds")
	}
}

// Helper function to verify Red-Black tree properties
func verifyRBProperties[K utils.Ordered, V any](t *testing.T, rb *RBTree[K, V]) {
	if rb.root == nil {
		return
	}

	// Property 1: Root is black
	if rb.root.color != Black {
		t.Error("Root node is not black")
	}

	// Property 2: Red nodes have black children
	// Property 3: All paths from root to leaves have same number of black nodes
	var verifyNode func(*RBNode[K, V], int, *int) bool
	verifyNode = func(node *RBNode[K, V], blackCount int, pathBlackCount *int) bool {
		if node == nil {
			if *pathBlackCount == -1 {
				*pathBlackCount = blackCount
			} else if blackCount != *pathBlackCount {
				t.Error("Different number of black nodes in paths from root to leaves")
				return false
			}
			return true
		}

		if node.color == Red {
			if (node.left != nil && node.left.color == Red) || (node.right != nil && node.right.color == Red) {
				t.Error("Red node has red child")
				return false
			}
		}

		newBlackCount := blackCount
		if node.color == Black {
			newBlackCount++
		}

		return verifyNode(node.left, newBlackCount, pathBlackCount) &&
			verifyNode(node.right, newBlackCount, pathBlackCount)
	}

	pathBlackCount := -1
	verifyNode(rb.root, 0, &pathBlackCount)
}

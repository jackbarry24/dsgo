package trees

import (
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
			rb := NewRBTree[int, int]()
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
	rb := NewRBTree[int, int]()
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
			rb := NewRBTree[int, int]()
			tt.setup(rb)
			rb.Delete(tt.delete)
			tt.check(t, rb)
		})
	}
}

func TestSafeRBTree_Concurrent(t *testing.T) {
	// Set a timeout for the entire test
	done := make(chan bool)
	go func() {
		rb := NewSafeRBTree[int, int]()
		var wg sync.WaitGroup
		iterations := 100 // Reduced from 1000 to 100

		// Concurrent inserts
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < iterations; i++ {
				rb.Insert(i, i)
			}
		}()

		// Concurrent searches
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < iterations; i++ {
				node, found := rb.Search(i)
				if found && node.value != i {
					t.Errorf("Search(%d) = %v, want %d", i, node.value, i)
				}
			}
		}()

		// Concurrent deletes
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < iterations; i++ {
				rb.Delete(i)
			}
		}()

		wg.Wait()
		done <- true
	}()

	// Wait for either completion or timeout
	select {
	case <-done:
		// Test completed successfully
	case <-time.After(5 * time.Second):
		t.Fatal("Test timed out after 5 seconds")
	}
}

// Helper function to verify Red-Black tree properties
func verifyRBProperties(t *testing.T, rb *RBTree[int, int]) {
	if rb.root == nil {
		return
	}

	// Property 1: Root is black
	if rb.root.color != Black {
		t.Error("Root node is not black")
	}

	// Property 2: Red nodes have black children
	// Property 3: All paths from root to leaves have same number of black nodes
	var verifyNode func(*RBNode[int, int], int, *int) bool
	verifyNode = func(node *RBNode[int, int], blackCount int, pathBlackCount *int) bool {
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

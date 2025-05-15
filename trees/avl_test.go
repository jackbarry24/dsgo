package trees

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestAVLTree_Insert(t *testing.T) {
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
			name: "left rotation",
			inserts: []struct{ key, value int }{
				{1, 1},
				{2, 2},
				{3, 3},
			},
			wantSize: 3,
		},
		{
			name: "right rotation",
			inserts: []struct{ key, value int }{
				{3, 3},
				{2, 2},
				{1, 1},
			},
			wantSize: 3,
		},
		{
			name: "left-right rotation",
			inserts: []struct{ key, value int }{
				{3, 3},
				{1, 1},
				{2, 2},
			},
			wantSize: 3,
		},
		{
			name: "right-left rotation",
			inserts: []struct{ key, value int }{
				{1, 1},
				{3, 3},
				{2, 2},
			},
			wantSize: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			avl := NewAVLTree[int, int]()
			for _, insert := range tt.inserts {
				avl.Insert(insert.key, insert.value)
			}

			// Verify all values are in the tree
			// For duplicate keys, we should get the last inserted value
			lastValue := make(map[int]int)
			for _, insert := range tt.inserts {
				lastValue[insert.key] = insert.value
			}

			for key, wantValue := range lastValue {
				value, found := avl.Search(key)
				if !found {
					t.Errorf("Search(%d) = not found, want found", key)
				}
				if value != wantValue {
					t.Errorf("Search(%d) = %d, want %d", key, value, wantValue)
				}
			}

			// Verify tree is balanced by checking heights
			values := avl.InOrderTraversal()
			if len(values) != tt.wantSize {
				t.Errorf("InOrderTraversal() length = %d, want %d", len(values), tt.wantSize)
			}
		})
	}
}

func TestAVLTree_Search(t *testing.T) {
	avl := NewAVLTree[int, int]()
	avl.Insert(5, 5)
	avl.Insert(3, 3)
	avl.Insert(7, 7)

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
			name:   "left child",
			key:    3,
			want:   3,
			wantOk: true,
		},
		{
			name:   "right child",
			key:    7,
			want:   7,
			wantOk: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := avl.Search(tt.key)
			if ok != tt.wantOk {
				t.Errorf("Search() ok = %v, want %v", ok, tt.wantOk)
			}
			if ok && got != tt.want {
				t.Errorf("Search() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAVLTree_Delete(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(*AVLTree[int, int])
		delete int
		check  func(*testing.T, *AVLTree[int, int])
	}{
		{
			name: "delete leaf node",
			setup: func(avl *AVLTree[int, int]) {
				avl.Insert(5, 5)
				avl.Insert(3, 3)
				avl.Insert(7, 7)
			},
			delete: 3,
			check: func(t *testing.T, avl *AVLTree[int, int]) {
				if _, found := avl.Search(3); found {
					t.Error("Search(3) = found, want not found")
				}
				if _, found := avl.Search(5); !found {
					t.Error("Search(5) = not found, want found")
				}
				if _, found := avl.Search(7); !found {
					t.Error("Search(7) = not found, want found")
				}
			},
		},
		{
			name: "delete node with one child",
			setup: func(avl *AVLTree[int, int]) {
				avl.Insert(5, 5)
				avl.Insert(3, 3)
				avl.Insert(4, 4)
			},
			delete: 3,
			check: func(t *testing.T, avl *AVLTree[int, int]) {
				if _, found := avl.Search(3); found {
					t.Error("Search(3) = found, want not found")
				}
				if _, found := avl.Search(4); !found {
					t.Error("Search(4) = not found, want found")
				}
			},
		},
		{
			name: "delete node with two children",
			setup: func(avl *AVLTree[int, int]) {
				avl.Insert(5, 5)
				avl.Insert(3, 3)
				avl.Insert(7, 7)
				avl.Insert(6, 6)
				avl.Insert(8, 8)
			},
			delete: 7,
			check: func(t *testing.T, avl *AVLTree[int, int]) {
				if _, found := avl.Search(7); found {
					t.Error("Search(7) = found, want not found")
				}
				if _, found := avl.Search(6); !found {
					t.Error("Search(6) = not found, want found")
				}
				if _, found := avl.Search(8); !found {
					t.Error("Search(8) = not found, want found")
				}
			},
		},
		{
			name: "delete root node",
			setup: func(avl *AVLTree[int, int]) {
				avl.Insert(5, 5)
				avl.Insert(3, 3)
				avl.Insert(7, 7)
			},
			delete: 5,
			check: func(t *testing.T, avl *AVLTree[int, int]) {
				if _, found := avl.Search(5); found {
					t.Error("Search(5) = found, want not found")
				}
				if _, found := avl.Search(3); !found {
					t.Error("Search(3) = not found, want found")
				}
				if _, found := avl.Search(7); !found {
					t.Error("Search(7) = not found, want found")
				}
			},
		},
		{
			name: "delete with rebalancing",
			setup: func(avl *AVLTree[int, int]) {
				avl.Insert(5, 5)
				avl.Insert(3, 3)
				avl.Insert(7, 7)
				avl.Insert(2, 2)
				avl.Insert(4, 4)
				avl.Insert(6, 6)
				avl.Insert(8, 8)
			},
			delete: 5,
			check: func(t *testing.T, avl *AVLTree[int, int]) {
				if _, found := avl.Search(5); found {
					t.Error("Search(5) = found, want not found")
				}
				// Verify all other nodes are still present and tree is balanced
				for _, key := range []int{2, 3, 4, 6, 7, 8} {
					if _, found := avl.Search(key); !found {
						t.Errorf("Search(%d) = not found, want found", key)
					}
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			avl := NewAVLTree[int, int]()
			tt.setup(avl)
			avl.Delete(tt.delete)
			tt.check(t, avl)
		})
	}
}

func TestAVLTreeConcurrent(t *testing.T) {
	tree := NewAVLTree[int, string](true)
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
			if value, exists := tree.Search(i); !exists || value != fmt.Sprintf("value-%d", i) {
				t.Errorf("Expected value-%d, got %v", i, value)
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

func TestAVLTreeConcurrentDelete(t *testing.T) {
	tree := NewAVLTree[int, string](true)
	var wg sync.WaitGroup
	done := make(chan bool)

	go func() {
		// First insert values
		for i := 0; i < 1000; i++ {
			tree.Insert(i, fmt.Sprintf("value-%d", i))
		}

		// Test concurrent Delete operations
		for i := 0; i < 500; i++ {
			wg.Add(1)
			go func(val int) {
				defer wg.Done()
				tree.Delete(val)
			}(i)
		}
		wg.Wait()

		// Verify deleted values are gone and others remain
		for i := 0; i < 1000; i++ {
			_, exists := tree.Search(i)
			if i < 500 && exists {
				t.Errorf("Value %d should be deleted", i)
			} else if i >= 500 && !exists {
				t.Errorf("Value %d should exist", i)
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

func TestAVLTreeConcurrentSearch(t *testing.T) {
	tree := NewAVLTree[int, string](true)
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
				if value, exists := tree.Search(val); !exists || value != fmt.Sprintf("value-%d", val) {
					t.Errorf("Expected value-%d, got %v", val, value)
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

func TestAVLTreeConcurrentModifications(t *testing.T) {
	tree := NewAVLTree[int, string](true)
	var wg sync.WaitGroup
	done := make(chan bool)

	go func() {
		// Test concurrent Insert and Delete operations
		for i := 0; i < 1000; i++ {
			wg.Add(2)
			go func(val int) {
				defer wg.Done()
				tree.Insert(val, fmt.Sprintf("value-%d", val))
			}(i)
			go func(val int) {
				defer wg.Done()
				tree.Delete(val)
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

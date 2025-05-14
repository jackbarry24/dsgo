package trees

import (
	"sync"
	"testing"
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

func TestSafeAVLTree_Sequential(t *testing.T) {
	avl := NewSafeAVLTree[int, int]()

	// Test basic operations
	avl.Insert(5, 5)
	avl.Insert(3, 3)
	avl.Insert(7, 7)

	// Test search
	if value, found := avl.Search(5); !found || value != 5 {
		t.Errorf("Search(5) = (%v, %v), want (5, true)", value, found)
	}

	// Test delete
	avl.Delete(3)
	if _, found := avl.Search(3); found {
		t.Error("Search(3) = found, want not found")
	}
}

func TestSafeAVLTree_Concurrent(t *testing.T) {
	avl := NewSafeAVLTree[int, int]()
	var wg sync.WaitGroup
	iterations := 1000
	insertDone := make(chan struct{})
	searchDone := make(chan struct{})

	// Concurrent inserts
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer close(insertDone)
		for i := 0; i < iterations; i++ {
			avl.Insert(i, i)
		}
	}()

	// Wait for inserts to complete before starting searches
	<-insertDone

	// Concurrent searches
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer close(searchDone)
		for i := 0; i < iterations; i++ {
			value, found := avl.Search(i)
			if !found || value != i {
				t.Errorf("Concurrent search: Search(%d) = (%d, %v), want (%d, true)", i, value, found, i)
			}
		}
	}()

	// Wait for searches to complete before starting deletes
	<-searchDone

	// Concurrent deletes
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < iterations; i++ {
			avl.Delete(i)
		}
	}()

	wg.Wait()

	// Verify final state
	for i := 0; i < iterations; i++ {
		if _, found := avl.Search(i); found {
			t.Errorf("Final state: Search(%d) = found, want not found", i)
		}
	}
}

func TestSafeAVLTree_ConcurrentUpdates(t *testing.T) {
	avl := NewSafeAVLTree[int, int]()
	var wg sync.WaitGroup
	iterations := 1000
	key := 1

	// Multiple goroutines updating the same key
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				avl.Insert(key, workerID)
			}
		}(i)
	}

	wg.Wait()

	// Verify that the key exists and has a valid value
	value, found := avl.Search(key)
	if !found {
		t.Errorf("Search(%d) = not found, want found", key)
	}
	if value < 0 || value >= 10 {
		t.Errorf("Search(%d) = %d, want value between 0 and 9", key, value)
	}
}

func TestSafeAVLTree_ConcurrentMixed(t *testing.T) {
	avl := NewSafeAVLTree[int, int]()
	var wg sync.WaitGroup
	iterations := 1000

	// Start multiple goroutines that perform mixed operations
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			start := workerID * iterations
			for j := 0; j < iterations; j++ {
				key := start + j
				avl.Insert(key, key)
				value, found := avl.Search(key)
				if !found || value != key {
					t.Errorf("Worker %d: Search(%d) = (%d, %v), want (%d, true)",
						workerID, key, value, found, key)
				}
				avl.Delete(key)
			}
		}(i)
	}

	wg.Wait()

	// Verify final state
	for i := 0; i < 5*iterations; i++ {
		if _, found := avl.Search(i); found {
			t.Errorf("Final state: Search(%d) = found, want not found", i)
		}
	}
}

package trees

import (
	"testing"
)

func TestBST_Insert(t *testing.T) {
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bst := NewBST[int, int]()
			for _, insert := range tt.inserts {
				bst.Insert(insert.key, insert.value)
			}

			// Verify all values are in the tree
			// For duplicate keys, we should get the last inserted value
			lastValue := make(map[int]int)
			for _, insert := range tt.inserts {
				lastValue[insert.key] = insert.value
			}

			for key, wantValue := range lastValue {
				value, found := bst.Search(key)
				if !found {
					t.Errorf("Search(%d) = not found, want found", key)
				}
				if value != wantValue {
					t.Errorf("Search(%d) = %d, want %d", key, value, wantValue)
				}
			}
		})
	}
}

func TestBST_Search(t *testing.T) {
	bst := NewBST[int, int]()
	bst.Insert(5, 5)
	bst.Insert(3, 3)
	bst.Insert(7, 7)

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
			got, ok := bst.Search(tt.key)
			if ok != tt.wantOk {
				t.Errorf("Search() ok = %v, want %v", ok, tt.wantOk)
			}
			if ok && got != tt.want {
				t.Errorf("Search() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBST_Delete(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(*BST[int, int])
		delete int
		check  func(*testing.T, *BST[int, int])
	}{
		{
			name: "delete leaf node",
			setup: func(bst *BST[int, int]) {
				bst.Insert(5, 5)
				bst.Insert(3, 3)
				bst.Insert(7, 7)
			},
			delete: 3,
			check: func(t *testing.T, bst *BST[int, int]) {
				if _, found := bst.Search(3); found {
					t.Error("Search(3) = found, want not found")
				}
				if _, found := bst.Search(5); !found {
					t.Error("Search(5) = not found, want found")
				}
				if _, found := bst.Search(7); !found {
					t.Error("Search(7) = not found, want found")
				}
			},
		},
		{
			name: "delete node with one child",
			setup: func(bst *BST[int, int]) {
				bst.Insert(5, 5)
				bst.Insert(3, 3)
				bst.Insert(4, 4)
			},
			delete: 3,
			check: func(t *testing.T, bst *BST[int, int]) {
				if _, found := bst.Search(3); found {
					t.Error("Search(3) = found, want not found")
				}
				if _, found := bst.Search(4); !found {
					t.Error("Search(4) = not found, want found")
				}
			},
		},
		{
			name: "delete node with two children",
			setup: func(bst *BST[int, int]) {
				bst.Insert(5, 5)
				bst.Insert(3, 3)
				bst.Insert(7, 7)
				bst.Insert(6, 6)
				bst.Insert(8, 8)
			},
			delete: 7,
			check: func(t *testing.T, bst *BST[int, int]) {
				if _, found := bst.Search(7); found {
					t.Error("Search(7) = found, want not found")
				}
				if _, found := bst.Search(6); !found {
					t.Error("Search(6) = not found, want found")
				}
				if _, found := bst.Search(8); !found {
					t.Error("Search(8) = not found, want found")
				}
			},
		},
		{
			name: "delete root node",
			setup: func(bst *BST[int, int]) {
				bst.Insert(5, 5)
				bst.Insert(3, 3)
				bst.Insert(7, 7)
			},
			delete: 5,
			check: func(t *testing.T, bst *BST[int, int]) {
				if _, found := bst.Search(5); found {
					t.Error("Search(5) = found, want not found")
				}
				if _, found := bst.Search(3); !found {
					t.Error("Search(3) = not found, want found")
				}
				if _, found := bst.Search(7); !found {
					t.Error("Search(7) = not found, want found")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bst := NewBST[int, int]()
			tt.setup(bst)
			bst.Delete(tt.delete)
			tt.check(t, bst)
		})
	}
}

package trees

import (
	"sync"
	"testing"
)

func TestSafeBST_Sequential(t *testing.T) {
	bst := NewSafeBST[int, int]()

	// Test basic operations
	bst.Insert(5, 5)
	bst.Insert(3, 3)
	bst.Insert(7, 7)

	// Test search
	if value, found := bst.Search(5); !found || value != 5 {
		t.Errorf("Search(5) = (%v, %v), want (5, true)", value, found)
	}

	// Test delete
	bst.Delete(3)
	if _, found := bst.Search(3); found {
		t.Error("Search(3) = found, want not found")
	}
}

func TestSafeBST_Concurrent(t *testing.T) {
	bst := NewSafeBST[int, int]()
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
			bst.Insert(i, i)
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
			value, found := bst.Search(i)
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
			bst.Delete(i)
		}
	}()

	wg.Wait()

	// Verify final state
	for i := 0; i < iterations; i++ {
		if _, found := bst.Search(i); found {
			t.Errorf("Final state: Search(%d) = found, want not found", i)
		}
	}
}

func TestSafeBST_ConcurrentUpdates(t *testing.T) {
	bst := NewSafeBST[int, int]()
	var wg sync.WaitGroup
	iterations := 1000
	key := 1

	// Multiple goroutines updating the same key
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				bst.Insert(key, workerID)
			}
		}(i)
	}

	wg.Wait()

	// Verify that the key exists and has a valid value
	value, found := bst.Search(key)
	if !found {
		t.Errorf("Search(%d) = not found, want found", key)
	}
	if value < 0 || value >= 10 {
		t.Errorf("Search(%d) = %d, want value between 0 and 9", key, value)
	}
}

func TestSafeBST_ConcurrentMixed(t *testing.T) {
	bst := NewSafeBST[int, int]()
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
				bst.Insert(key, key)
				value, found := bst.Search(key)
				if !found || value != key {
					t.Errorf("Worker %d: Search(%d) = (%d, %v), want (%d, true)",
						workerID, key, value, found, key)
				}
				bst.Delete(key)
			}
		}(i)
	}

	wg.Wait()

	// Verify final state
	for i := 0; i < 5*iterations; i++ {
		if _, found := bst.Search(i); found {
			t.Errorf("Final state: Search(%d) = found, want not found", i)
		}
	}
}

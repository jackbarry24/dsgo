package sets

import (
	"sync"
	"testing"
)

func TestNew(t *testing.T) {
	s := NewSet[int]()
	if s.Size() != 0 {
		t.Errorf("New set should be empty, got size %d", s.Size())
	}
}

func TestAddAndContains(t *testing.T) {
	s := NewSet[int]()
	s.Add(1)
	s.Add(2)
	s.Add(1) // Adding duplicate

	if !s.Contains(1) {
		t.Error("Set should contain 1")
	}
	if !s.Contains(2) {
		t.Error("Set should contain 2")
	}
	if s.Size() != 2 {
		t.Errorf("Set should have size 2, got %d", s.Size())
	}
}

func TestRemove(t *testing.T) {
	s := NewSet[string]()
	s.Add("a")
	s.Add("b")
	s.Remove("a")

	if s.Contains("a") {
		t.Error("Set should not contain 'a' after removal")
	}
	if !s.Contains("b") {
		t.Error("Set should still contain 'b'")
	}
	if s.Size() != 1 {
		t.Errorf("Set should have size 1, got %d", s.Size())
	}
}

func TestUnion(t *testing.T) {
	s1 := NewSet[int]()
	s2 := NewSet[int]()
	s1.Add(1)
	s1.Add(2)
	s2.Add(2)
	s2.Add(3)

	union := s1.Union(s2)
	if union.Size() != 3 {
		t.Errorf("Union should have size 3, got %d", union.Size())
	}
	if !union.Contains(1) || !union.Contains(2) || !union.Contains(3) {
		t.Error("Union should contain all elements from both sets")
	}
}

func TestIntersection(t *testing.T) {
	s1 := NewSet[int]()
	s2 := NewSet[int]()
	s1.Add(1)
	s1.Add(2)
	s2.Add(2)
	s2.Add(3)

	intersection := s1.Intersection(s2)
	if intersection.Size() != 1 {
		t.Errorf("Intersection should have size 1, got %d", intersection.Size())
	}
	if !intersection.Contains(2) {
		t.Error("Intersection should contain 2")
	}
}

func TestDifference(t *testing.T) {
	s1 := NewSet[int]()
	s2 := NewSet[int]()
	s1.Add(1)
	s1.Add(2)
	s2.Add(2)
	s2.Add(3)

	diff := s1.Difference(s2)
	if diff.Size() != 1 {
		t.Errorf("Difference should have size 1, got %d", diff.Size())
	}
	if !diff.Contains(1) {
		t.Error("Difference should contain 1")
	}
}

func TestClear(t *testing.T) {
	s := NewSet[int]()
	s.Add(1)
	s.Add(2)
	s.Clear()

	if !s.IsEmpty() {
		t.Error("Set should be empty after Clear")
	}
	if s.Size() != 0 {
		t.Errorf("Set should have size 0 after Clear, got %d", s.Size())
	}
}

func TestItems(t *testing.T) {
	s := NewSet[int]()
	s.Add(1)
	s.Add(2)
	s.Add(3)

	items := s.Items()
	if len(items) != 3 {
		t.Errorf("Items should return 3 elements, got %d", len(items))
	}

	// Create a map to check if all items are present
	itemMap := make(map[int]bool)
	for _, item := range items {
		itemMap[item] = true
	}
	for i := 1; i <= 3; i++ {
		if !itemMap[i] {
			t.Errorf("Items should contain %d", i)
		}
	}
}

func TestSetBasicOperations(t *testing.T) {
	s := NewSet[int](true)

	// Test empty set
	if !s.IsEmpty() {
		t.Error("New set should be empty")
	}
	if s.Size() != 0 {
		t.Error("New set should have size 0")
	}

	// Test Add
	s.Add(1)
	s.Add(2)
	s.Add(3)

	if s.Size() != 3 {
		t.Errorf("Expected size 3, got %d", s.Size())
	}

	// Test Contains
	if !s.Contains(1) {
		t.Error("Set should contain 1")
	}
	if s.Contains(4) {
		t.Error("Set should not contain 4")
	}

	// Test Remove
	s.Remove(2)
	if s.Contains(2) {
		t.Error("Set should not contain 2 after removal")
	}
	if s.Size() != 2 {
		t.Errorf("Expected size 2, got %d", s.Size())
	}

	// Test Clear
	s.Clear()
	if !s.IsEmpty() {
		t.Error("Set should be empty after Clear")
	}
}

func TestSetOperations(t *testing.T) {
	s1 := NewSet[int](true)
	s2 := NewSet[int](true)

	s1.Add(1)
	s1.Add(2)
	s1.Add(3)

	s2.Add(2)
	s2.Add(3)
	s2.Add(4)

	// Test Union
	union := s1.Union(s2)
	expectedUnion := map[int]bool{1: true, 2: true, 3: true, 4: true}
	for _, item := range union.Items() {
		if !expectedUnion[item] {
			t.Errorf("Union contains unexpected item: %d", item)
		}
	}

	// Test Intersection
	intersection := s1.Intersection(s2)
	expectedIntersection := map[int]bool{2: true, 3: true}
	for _, item := range intersection.Items() {
		if !expectedIntersection[item] {
			t.Errorf("Intersection contains unexpected item: %d", item)
		}
	}

	// Test Difference
	difference := s1.Difference(s2)
	expectedDifference := map[int]bool{1: true}
	for _, item := range difference.Items() {
		if !expectedDifference[item] {
			t.Errorf("Difference contains unexpected item: %d", item)
		}
	}
}

func TestSetConcurrentOperations(t *testing.T) {
	s := NewSet[int](true)
	var wg sync.WaitGroup

	// Test concurrent Add operations
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func(val int) {
			defer wg.Done()
			s.Add(val)
		}(i)
	}
	wg.Wait()

	if s.Size() != 1000 {
		t.Errorf("Expected size 1000 after concurrent adds, got %d", s.Size())
	}

	// Test concurrent Remove operations
	for i := 0; i < 500; i++ {
		wg.Add(1)
		go func(val int) {
			defer wg.Done()
			s.Remove(val)
		}(i)
	}
	wg.Wait()

	if s.Size() != 500 {
		t.Errorf("Expected size 500 after concurrent removes, got %d", s.Size())
	}
}

func TestSetConcurrentModifications(t *testing.T) {
	s := NewSet[int](true)
	var wg sync.WaitGroup

	// Test concurrent Add and Remove operations
	for i := 0; i < 1000; i++ {
		wg.Add(2)
		go func(val int) {
			defer wg.Done()
			s.Add(val)
		}(i)
		go func(val int) {
			defer wg.Done()
			s.Remove(val)
		}(i)
	}
	wg.Wait()
}

func TestSetConcurrent(t *testing.T) {
	set := NewSet[int](true)
	var wg sync.WaitGroup

	// Test concurrent Add operations
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func(val int) {
			defer wg.Done()
			set.Add(val)
		}(i)
	}
	wg.Wait()

	if set.Size() != 1000 {
		t.Errorf("Expected size 1000 after concurrent adds, got %d", set.Size())
	}
}

func TestSetUnionConcurrent(t *testing.T) {
	set1 := NewSet[int](true)
	set2 := NewSet[int](true)
	var wg sync.WaitGroup

	// Add elements concurrently to both sets
	for i := 0; i < 1000; i++ {
		wg.Add(2)
		go func(val int) {
			defer wg.Done()
			set1.Add(val)
		}(i)
		go func(val int) {
			defer wg.Done()
			set2.Add(val + 500)
		}(i)
	}
	wg.Wait()

	union := set1.Union(set2)
	if union.Size() != 1500 {
		t.Errorf("Expected union size 1500, got %d", union.Size())
	}
}

func TestSetIntersectionConcurrent(t *testing.T) {
	set1 := NewSet[int](true)
	set2 := NewSet[int](true)
	var wg sync.WaitGroup

	// Add elements concurrently to both sets
	for i := 0; i < 1000; i++ {
		wg.Add(2)
		go func(val int) {
			defer wg.Done()
			set1.Add(val)
		}(i)
		go func(val int) {
			defer wg.Done()
			set2.Add(val)
		}(i)
	}
	wg.Wait()

	intersection := set1.Intersection(set2)
	if intersection.Size() != 1000 {
		t.Errorf("Expected intersection size 1000, got %d", intersection.Size())
	}
}

func TestSetDifferenceConcurrent(t *testing.T) {
	set1 := NewSet[int](true)
	set2 := NewSet[int](true)
	var wg sync.WaitGroup

	// Add elements concurrently to both sets
	for i := 0; i < 1000; i++ {
		wg.Add(2)
		go func(val int) {
			defer wg.Done()
			set1.Add(val)
		}(i)
		go func(val int) {
			defer wg.Done()
			set2.Add(val + 500)
		}(i)
	}
	wg.Wait()

	difference := set1.Difference(set2)
	if difference.Size() != 500 {
		t.Errorf("Expected difference size 500, got %d", difference.Size())
	}
}

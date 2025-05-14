package sets

import (
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

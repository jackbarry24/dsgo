package graphs

import (
	"testing"
)

func TestNewGraph(t *testing.T) {
	g := NewGraph[string, int]()
	if g.nodes == nil || g.edges == nil {
		t.Error("NewGraph should initialize nodes and edges maps")
	}
}

func TestNewSafeGraph(t *testing.T) {
	sg := NewSafeGraph[string, int]()
	if sg.inner == nil {
		t.Error("NewSafeGraph should initialize inner graph")
	}
}

func TestGraphBasicOperations(t *testing.T) {
	g := NewGraph[string, int]()

	// Test AddNode
	g.AddNode("A", 1)
	if !g.HasNode("A") {
		t.Error("AddNode failed to add node")
	}

	// Test AddEdge
	g.AddNode("B", 2)
	g.AddEdge("A", "B")
	if !g.HasEdge("A", "B") {
		t.Error("AddEdge failed to add edge")
	}

	// Test GetNeighbors
	neighbors := g.GetNeighbors("A")
	if len(neighbors) != 1 || neighbors[0] != "B" {
		t.Error("GetNeighbors returned incorrect neighbors")
	}

	// Test GetNodeValue
	if value, exists := g.GetNodeValue("A"); !exists || value != 1 {
		t.Error("GetNodeValue returned incorrect value")
	}

	// Test RemoveEdge
	g.RemoveEdge("A", "B")
	if g.HasEdge("A", "B") {
		t.Error("RemoveEdge failed to remove edge")
	}

	// Test RemoveNode
	g.RemoveNode("A")
	if g.HasNode("A") {
		t.Error("RemoveNode failed to remove node")
	}
}

func TestGraphTraversal(t *testing.T) {
	g := NewGraph[string, int]()

	// Create a simple graph: A -> B -> C
	//                      \-> D -> E
	g.AddNode("A", 1)
	g.AddNode("B", 2)
	g.AddNode("C", 3)
	g.AddNode("D", 4)
	g.AddNode("E", 5)
	g.AddEdge("A", "B")
	g.AddEdge("B", "C")
	g.AddEdge("A", "D")
	g.AddEdge("D", "E")

	// Test BFS
	bfsResult := g.BFS("A")
	expectedBFS := []string{"A", "B", "D", "C", "E"}
	if len(bfsResult) != len(expectedBFS) {
		t.Errorf("BFS returned wrong length: got %v, want %v", len(bfsResult), len(expectedBFS))
	}
	for i, v := range bfsResult {
		if v != expectedBFS[i] {
			t.Errorf("BFS order wrong at index %d: got %v, want %v", i, v, expectedBFS[i])
		}
	}

	// Test DFS
	dfsResult := g.DFS("A")
	expectedDFS := []string{"A", "B", "C", "D", "E"}
	if len(dfsResult) != len(expectedDFS) {
		t.Errorf("DFS returned wrong length: got %v, want %v", len(dfsResult), len(expectedDFS))
	}
	for i, v := range dfsResult {
		if v != expectedDFS[i] {
			t.Errorf("DFS order wrong at index %d: got %v, want %v", i, v, expectedDFS[i])
		}
	}
}

func TestGraphGetNodesAndEdges(t *testing.T) {
	g := NewGraph[string, int]()

	// Add some nodes and edges
	g.AddNode("A", 1)
	g.AddNode("B", 2)
	g.AddNode("C", 3)
	g.AddEdge("A", "B")
	g.AddEdge("B", "C")

	// Test GetNodes
	nodes := g.GetNodes()
	if len(nodes) != 3 {
		t.Errorf("GetNodes returned wrong number of nodes: got %v, want 3", len(nodes))
	}

	// Test GetEdges
	edges := g.GetEdges()
	if len(edges) != 2 {
		t.Errorf("GetEdges returned wrong number of edges: got %v, want 2", len(edges))
	}
}

func TestSafeGraphConcurrent(t *testing.T) {
	sg := NewSafeGraph[string, int]()

	// Test concurrent operations
	done := make(chan bool)

	// Concurrently add nodes
	for i := 0; i < 100; i++ {
		go func(n int) {
			sg.AddNode(string(rune(n)), n)
			done <- true
		}(i)
	}

	// Wait for all goroutines to complete
	for i := 0; i < 100; i++ {
		<-done
	}

	// Verify all nodes were added
	nodes := sg.GetNodes()
	if len(nodes) != 100 {
		t.Errorf("SafeGraph concurrent AddNode failed: got %v nodes, want 100", len(nodes))
	}
}

func TestGraphEdgeCases(t *testing.T) {
	g := NewGraph[string, int]()

	// Test non-existent node operations
	if g.HasNode("X") {
		t.Error("HasNode should return false for non-existent node")
	}
	if g.HasEdge("X", "Y") {
		t.Error("HasEdge should return false for non-existent edge")
	}
	if neighbors := g.GetNeighbors("X"); neighbors != nil {
		t.Error("GetNeighbors should return nil for non-existent node")
	}
	if _, exists := g.GetNodeValue("X"); exists {
		t.Error("GetNodeValue should return false for non-existent node")
	}

	// Test BFS/DFS with non-existent start node
	if bfs := g.BFS("X"); bfs != nil {
		t.Error("BFS should return nil for non-existent start node")
	}
	if dfs := g.DFS("X"); dfs != nil {
		t.Error("DFS should return nil for non-existent start node")
	}

	// Test removing non-existent node/edge
	g.RemoveNode("X")      // Should not panic
	g.RemoveEdge("X", "Y") // Should not panic
}

package graphs

import (
	"fmt"
	"sort"
	"sync"
)

type Graph[K comparable, V any] struct {
	threadSafe bool
	mu         sync.RWMutex
	nodes      map[K]V
	edges      map[K]map[K]struct{}
}

// NewGraph creates a new graph. If threadSafe is true, the graph will be safe for concurrent access.
func NewGraph[K comparable, V any](threadSafe ...bool) *Graph[K, V] {
	isThreadSafe := true
	if len(threadSafe) > 0 {
		isThreadSafe = threadSafe[0]
	}
	return &Graph[K, V]{
		threadSafe: isThreadSafe,
		nodes:      make(map[K]V),
		edges:      make(map[K]map[K]struct{}),
	}
}

// AddNode adds a node to the graph with the given key and value.
func (g *Graph[K, V]) AddNode(key K, value V) {
	if g.threadSafe {
		g.mu.Lock()
		defer g.mu.Unlock()
	}
	g.nodes[key] = value
}

// AddEdge adds a directed edge from 'from' to 'to'.
func (g *Graph[K, V]) AddEdge(from, to K) {
	if g.threadSafe {
		g.mu.Lock()
		defer g.mu.Unlock()
	}
	if _, exists := g.edges[from]; !exists {
		g.edges[from] = make(map[K]struct{})
	}
	g.edges[from][to] = struct{}{}
}

// HasNode checks if a node with the given key exists.
func (g *Graph[K, V]) HasNode(key K) bool {
	if g.threadSafe {
		g.mu.RLock()
		defer g.mu.RUnlock()
	}
	_, exists := g.nodes[key]
	return exists
}

// HasEdge checks if an edge exists from 'from' to 'to'.
func (g *Graph[K, V]) HasEdge(from, to K) bool {
	if g.threadSafe {
		g.mu.RLock()
		defer g.mu.RUnlock()
	}
	if neighbors, exists := g.edges[from]; exists {
		_, hasEdge := neighbors[to]
		return hasEdge
	}
	return false
}

// RemoveNode removes a node and all its associated edges.
func (g *Graph[K, V]) RemoveNode(key K) {
	if g.threadSafe {
		g.mu.Lock()
		defer g.mu.Unlock()
	}
	delete(g.nodes, key)
	delete(g.edges, key)
	// Remove all edges pointing to this node
	for _, neighbors := range g.edges {
		delete(neighbors, key)
	}
}

// RemoveEdge removes the edge from 'from' to 'to'.
func (g *Graph[K, V]) RemoveEdge(from, to K) {
	if g.threadSafe {
		g.mu.Lock()
		defer g.mu.Unlock()
	}
	if neighbors, exists := g.edges[from]; exists {
		delete(neighbors, to)
	}
}

// GetNeighbors returns all neighbors of the given node.
func (g *Graph[K, V]) GetNeighbors(key K) []K {
	if g.threadSafe {
		g.mu.RLock()
		defer g.mu.RUnlock()
	}
	if neighbors, exists := g.edges[key]; exists {
		result := make([]K, 0, len(neighbors))
		for neighbor := range neighbors {
			result = append(result, neighbor)
		}
		return result
	}
	return nil
}

// GetNodes returns all node keys in the graph.
func (g *Graph[K, V]) GetNodes() []K {
	if g.threadSafe {
		g.mu.RLock()
		defer g.mu.RUnlock()
	}
	nodes := make([]K, 0, len(g.nodes))
	for node := range g.nodes {
		nodes = append(nodes, node)
	}
	return nodes
}

// GetEdges returns all edges in the graph as pairs of [from, to] keys.
func (g *Graph[K, V]) GetEdges() [][2]K {
	if g.threadSafe {
		g.mu.RLock()
		defer g.mu.RUnlock()
	}
	edges := make([][2]K, 0)
	for from, neighbors := range g.edges {
		for to := range neighbors {
			edges = append(edges, [2]K{from, to})
		}
	}
	return edges
}

// GetNodeValue returns the value associated with a node key.
func (g *Graph[K, V]) GetNodeValue(key K) (V, bool) {
	if g.threadSafe {
		g.mu.RLock()
		defer g.mu.RUnlock()
	}
	value, exists := g.nodes[key]
	return value, exists
}

// BFS performs a breadth-first search starting from the given node.
func (g *Graph[K, V]) BFS(start K) []K {
	if g.threadSafe {
		g.mu.RLock()
		defer g.mu.RUnlock()
	}
	if !g.HasNode(start) {
		return nil
	}

	visited := make(map[K]bool)
	queue := []K{start}
	result := make([]K, 0)
	visited[start] = true

	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]
		result = append(result, node)

		// Get neighbors and sort them to ensure consistent order
		neighbors := g.GetNeighbors(node)
		sort.Slice(neighbors, func(i, j int) bool {
			return fmt.Sprintf("%v", neighbors[i]) < fmt.Sprintf("%v", neighbors[j])
		})

		for _, neighbor := range neighbors {
			if !visited[neighbor] {
				visited[neighbor] = true
				queue = append(queue, neighbor)
			}
		}
	}
	return result
}

// DFS performs a depth-first search starting from the given node.
func (g *Graph[K, V]) DFS(start K) []K {
	if g.threadSafe {
		g.mu.RLock()
		defer g.mu.RUnlock()
	}
	if !g.HasNode(start) {
		return nil
	}

	visited := make(map[K]bool)
	result := make([]K, 0)

	var dfs func(node K)
	dfs = func(node K) {
		visited[node] = true
		result = append(result, node)

		neighbors := g.GetNeighbors(node)
		sort.Slice(neighbors, func(i, j int) bool {
			return fmt.Sprintf("%v", neighbors[i]) < fmt.Sprintf("%v", neighbors[j])
		})
		for _, neighbor := range neighbors {
			if !visited[neighbor] {
				dfs(neighbor)
			}
		}
	}

	dfs(start)
	return result
}

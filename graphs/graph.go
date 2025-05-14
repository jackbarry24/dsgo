package graphs

import "sync"

type Graph[K comparable, V any] struct {
	nodes map[K]V
	edges map[K]map[K]struct{}
}

type SafeGraph[K comparable, V any] struct {
	mu    sync.RWMutex
	inner *Graph[K, V]
}

func NewGraph[K comparable, V any]() *Graph[K, V] {
	return &Graph[K, V]{
		nodes: make(map[K]V),
		edges: make(map[K]map[K]struct{}),
	}
}

func NewSafeGraph[K comparable, V any]() *SafeGraph[K, V] {
	return &SafeGraph[K, V]{
		mu:    sync.RWMutex{},
		inner: NewGraph[K, V](),
	}
}

func (g *Graph[K, V]) AddNode(key K, value V) {
	g.nodes[key] = value
}

func (g *Graph[K, V]) AddEdge(from, to K) {
	if _, exists := g.edges[from]; !exists {
		g.edges[from] = make(map[K]struct{})
	}
	g.edges[from][to] = struct{}{}
}

func (g *SafeGraph[K, V]) AddNode(key K, value V) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.inner.AddNode(key, value)
}

func (g *SafeGraph[K, V]) AddEdge(from, to K) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.inner.AddEdge(from, to)
}

func (g *Graph[K, V]) HasNode(key K) bool {
	_, exists := g.nodes[key]
	return exists
}

func (g *Graph[K, V]) HasEdge(from, to K) bool {
	if neighbors, exists := g.edges[from]; exists {
		_, hasEdge := neighbors[to]
		return hasEdge
	}
	return false
}

func (g *Graph[K, V]) RemoveNode(key K) {
	delete(g.nodes, key)
	delete(g.edges, key)
	// Remove all edges pointing to this node
	for _, neighbors := range g.edges {
		delete(neighbors, key)
	}
}

func (g *Graph[K, V]) RemoveEdge(from, to K) {
	if neighbors, exists := g.edges[from]; exists {
		delete(neighbors, to)
	}
}

func (g *Graph[K, V]) GetNeighbors(key K) []K {
	if neighbors, exists := g.edges[key]; exists {
		result := make([]K, 0, len(neighbors))
		for neighbor := range neighbors {
			result = append(result, neighbor)
		}
		return result
	}
	return nil
}

// Thread-safe versions of the above methods
func (g *SafeGraph[K, V]) HasNode(key K) bool {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return g.inner.HasNode(key)
}

func (g *SafeGraph[K, V]) HasEdge(from, to K) bool {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return g.inner.HasEdge(from, to)
}

func (g *SafeGraph[K, V]) RemoveNode(key K) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.inner.RemoveNode(key)
}

func (g *SafeGraph[K, V]) RemoveEdge(from, to K) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.inner.RemoveEdge(from, to)
}

func (g *SafeGraph[K, V]) GetNeighbors(key K) []K {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return g.inner.GetNeighbors(key)
}

func (g *Graph[K, V]) GetNodes() []K {
	nodes := make([]K, 0, len(g.nodes))
	for node := range g.nodes {
		nodes = append(nodes, node)
	}
	return nodes
}

func (g *Graph[K, V]) GetEdges() [][2]K {
	edges := make([][2]K, 0)
	for from, neighbors := range g.edges {
		for to := range neighbors {
			edges = append(edges, [2]K{from, to})
		}
	}
	return edges
}

func (g *Graph[K, V]) GetNodeValue(key K) (V, bool) {
	value, exists := g.nodes[key]
	return value, exists
}

func (g *Graph[K, V]) BFS(start K) []K {
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

		for _, neighbor := range g.GetNeighbors(node) {
			if !visited[neighbor] {
				visited[neighbor] = true
				queue = append(queue, neighbor)
			}
		}
	}
	return result
}

func (g *Graph[K, V]) DFS(start K) []K {
	if !g.HasNode(start) {
		return nil
	}

	visited := make(map[K]bool)
	result := make([]K, 0)

	var dfs func(node K)
	dfs = func(node K) {
		visited[node] = true
		result = append(result, node)

		for _, neighbor := range g.GetNeighbors(node) {
			if !visited[neighbor] {
				dfs(neighbor)
			}
		}
	}

	dfs(start)
	return result
}

// Thread-safe versions
func (g *SafeGraph[K, V]) GetNodes() []K {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return g.inner.GetNodes()
}

func (g *SafeGraph[K, V]) GetEdges() [][2]K {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return g.inner.GetEdges()
}

func (g *SafeGraph[K, V]) GetNodeValue(key K) (V, bool) {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return g.inner.GetNodeValue(key)
}

func (g *SafeGraph[K, V]) BFS(start K) []K {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return g.inner.BFS(start)
}

func (g *SafeGraph[K, V]) DFS(start K) []K {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return g.inner.DFS(start)
}

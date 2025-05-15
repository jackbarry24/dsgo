# DSGo - Go Data Structures Library

[![Go Report Card](https://goreportcard.com/badge/github.com/jackbarry24/dsgo)](https://goreportcard.com/report/github.com/jackbarry24/dsgo)
[![GoDoc](https://godoc.org/github.com/jackbarry24/dsgo?status.svg)](https://godoc.org/github.com/jackbarry24/dsgo)
[![CI Status](https://github.com/jackbarry24/dsgo/actions/workflows/ci.yml/badge.svg)](https://github.com/jackbarry24/dsgo/actions/workflows/ci.yml)

A comprehensive collection of data structures implemented in Go for Go.

## Data Structures

### Maps
- `OrderedMap`: A map that maintains insertion order
- `SortedMap`: A map that maintains keys in sorted order
- `SafeSortedMap`: Thread-safe version of SortedMap

### Sets
- Generic Set implementation with operations like:
  - Union
  - Intersection
  - Difference
  - Basic set operations (Add, Remove, Contains)

### Trees
- `AVLTree`: Self-balancing binary search tree
- `BST`: Binary Search Tree
- `RBTree`: Red-Black Tree implementation

### Heaps
- `MinHeap`: Binary min heap implementation
- `PriorityQueue`: Priority queue based on min heap

### Graphs
- Generic graph implementation with:
  - BFS and DFS traversal
  - Node and edge management
  - Neighbor operations

### Linked Lists
- `SingleLinkedList`: Singly linked list implementation
- `DoubleLinkedList`: Doubly linked list implementation

### Cache
- `LRUCache`: Least Recently Used (LRU) cache implementation

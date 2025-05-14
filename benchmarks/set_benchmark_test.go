package benchmarks

import (
	"testing"

	"dsgo/sets"
)

func BenchmarkSet(b *testing.B) {
	config := DefaultConfig()

	operations := []Operation[int]{
		{
			Name: "Add",
			Function: func(item int) {
				s := sets.NewSet[int]()
				s.Add(item)
			},
		},
		{
			Name: "Contains",
			Function: func(item int) {
				s := sets.NewSet[int]()
				s.Add(item)
				s.Contains(item)
			},
		},
		{
			Name: "Remove",
			Function: func(item int) {
				s := sets.NewSet[int]()
				s.Add(item)
				s.Remove(item)
			},
		},
	}

	RunBenchmarkSuite(b, config, func() interface{} {
		return sets.NewSet[int]()
	}, operations)
}

func BenchmarkSafeSet(b *testing.B) {
	config := DefaultConfig()

	operations := []Operation[int]{
		{
			Name: "Add",
			Function: func(item int) {
				s := sets.NewSafeSet[int]()
				s.Add(item)
			},
		},
		{
			Name: "Contains",
			Function: func(item int) {
				s := sets.NewSafeSet[int]()
				s.Add(item)
				s.Contains(item)
			},
		},
		{
			Name: "Remove",
			Function: func(item int) {
				s := sets.NewSafeSet[int]()
				s.Add(item)
				s.Remove(item)
			},
		},
	}

	RunBenchmarkSuite(b, config, func() any {
		return sets.NewSafeSet[int]()
	}, operations)
}

func BenchmarkSetConcurrent(b *testing.B) {
	config := DefaultConfig()

	operations := []Operation[int]{
		{
			Name: "Add",
			Function: func(item int) {
				s := sets.NewSet[int]()
				s.Add(item)
			},
		},
		{
			Name: "Contains",
			Function: func(item int) {
				s := sets.NewSet[int]()
				s.Add(item)
				s.Contains(item)
			},
		},
		{
			Name: "Remove",
			Function: func(item int) {
				s := sets.NewSet[int]()
				s.Add(item)
				s.Remove(item)
			},
		},
	}

	for _, op := range operations {
		b.Run(op.Name, func(b *testing.B) {
			RunConcurrentBenchmark(b, config, func() any {
				return sets.NewSet[int]()
			}, op)
		})
	}
}

func BenchmarkSafeSetConcurrent(b *testing.B) {
	config := DefaultConfig()

	operations := []Operation[int]{
		{
			Name: "Add",
			Function: func(item int) {
				s := sets.NewSafeSet[int]()
				s.Add(item)
			},
		},
		{
			Name: "Contains",
			Function: func(item int) {
				s := sets.NewSafeSet[int]()
				s.Add(item)
				s.Contains(item)
			},
		},
		{
			Name: "Remove",
			Function: func(item int) {
				s := sets.NewSafeSet[int]()
				s.Add(item)
				s.Remove(item)
			},
		},
	}

	for _, op := range operations {
		b.Run(op.Name, func(b *testing.B) {
			RunConcurrentBenchmark(b, config, func() any {
				return sets.NewSafeSet[int]()
			}, op)
		})
	}
}

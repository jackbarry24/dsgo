package benchmarks

import (
	"math/rand"
	"sync"
	"testing"
	"time"
)

type Operation[T comparable] struct {
	Name     string
	Function func(T)
}

type BenchmarkConfig struct {
	// Number of elements to use in the benchmark
	Size int
	// Number of concurrent operations for concurrent benchmarks
	Concurrency int
	// Random seed for reproducibility
	Seed int64
}

func DefaultConfig() BenchmarkConfig {
	return BenchmarkConfig{
		Size:        10000,
		Concurrency: 4,
		Seed:        time.Now().UnixNano(),
	}
}

func GenerateRandomData[T comparable](size int, seed int64) []T {
	r := rand.New(rand.NewSource(seed))
	data := make([]T, size)

	switch any(*new(T)).(type) {
	case int:
		for i := 0; i < size; i++ {
			data[i] = any(r.Int()).(T)
		}
	case int64:
		for i := 0; i < size; i++ {
			data[i] = any(r.Int63()).(T)
		}
	case float64:
		for i := 0; i < size; i++ {
			data[i] = any(r.Float64()).(T)
		}
	case string:
		const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
		for i := 0; i < size; i++ {
			b := make([]byte, 10)
			for j := range b {
				b[j] = charset[r.Intn(len(charset))]
			}
			data[i] = any(string(b)).(T)
		}
	default:
		panic("unsupported type for random data generation")
	}

	return data
}

func RunBenchmark[T comparable](b *testing.B, config BenchmarkConfig, setup func() any, operation Operation[T]) {
	data := GenerateRandomData[T](config.Size, config.Seed)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = setup()
		for _, item := range data {
			operation.Function(item)
		}
	}
}

func RunConcurrentBenchmark[T comparable](b *testing.B, config BenchmarkConfig, setup func() any, operation Operation[T]) {
	data := GenerateRandomData[T](config.Size, config.Seed)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = setup()
		var wg sync.WaitGroup
		for j := 0; j < config.Concurrency; j++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for _, item := range data {
					operation.Function(item)
				}
			}()
		}
		wg.Wait()
	}
}

type BenchmarkResult struct {
	Operation           string
	Iterations          int
	Duration            time.Duration
	OperationsPerSecond float64
}

func RunBenchmarkSuite[T comparable](b *testing.B, config BenchmarkConfig, setup func() any, operations []Operation[T]) []BenchmarkResult {
	results := make([]BenchmarkResult, len(operations))

	for i, op := range operations {
		b.Run(op.Name, func(b *testing.B) {
			RunBenchmark(b, config, setup, op)
		})

		results[i] = BenchmarkResult{
			Operation:           op.Name,
			Iterations:          b.N,
			Duration:            b.Elapsed(),
			OperationsPerSecond: float64(b.N) / b.Elapsed().Seconds(),
		}
	}

	return results
}

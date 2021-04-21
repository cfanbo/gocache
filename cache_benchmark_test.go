package gocache_test

import (
	"strconv"
	"sync"
	"testing"
)

func Benchmark_SetMap(b *testing.B) {
	m := make(map[string]interface{})
	var mu sync.Mutex
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mu.Lock()
		m["foo"] = i
		mu.Unlock()
	}
}
func Benchmark_Set(b *testing.B) {
	for i := 0; i < b.N; i++ {
		cache.Set("foo", i)
	}
}

func Benchmark_SetParallelMap(b *testing.B) {
	m := make(map[string]interface{})
	var mu sync.Mutex
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			mu.Lock()
			m["foo"] = 1
			mu.Unlock()
		}
	})
}
func Benchmark_ParalleSetl(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			cache.Set("foo", 1)
		}
	})
}

func Benchmark_Get(b *testing.B) {
	for i := 0; i < b.N; i++ {
		v := strconv.Itoa(i)
		cache.Set(v, v)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		v := strconv.Itoa(i)
		cache.Get(v)
	}
}

func Benchmark_Delete(b *testing.B) {
	for i := 0; i < b.N; i++ {
		cache.Delete("foo")
	}
}

func Benchmark_Exists(b *testing.B) {
	cache.Set("foo", "bar")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.Exists("foo")
	}
}

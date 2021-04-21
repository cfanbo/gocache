package gocache_test

import (
	"strconv"
	"testing"
)

func Benchmark_pcache(b *testing.B) {
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			for i := 0; i < 10; i++ {
				cache.Set(strconv.Itoa(i), i)
			}
		}
	})

}

package gocache

import (
	"sync"
	"sync/atomic"
	"time"
)

type shardMap []*CacheBucket

var SHARD_COUNT = 64

func NewShardMap() shardMap {
	m := make([]*CacheBucket, SHARD_COUNT)
	for i := 0; i < SHARD_COUNT; i++ {
		m[i] = NewCacheBucket()
	}

	return m
}

// GetShard Get a CacheBucket
func (m shardMap) GetShard(key string) *CacheBucket {
	i := m.getHashCode(key)

	return m[i%SHARD_COUNT]
}

// getHashCode
func (m shardMap) getHashCode(key string) int {
	// v := int(crc32.ChecksumIEEE([]byte(key)))
	v := int(fnv32(key))
	i := 0
	if v >= 0 {
		i = v
	} else {
		if -v >= 0 {
			i = -v
		}
	}

	return i
}

func fnv32(key string) uint32 {
	hash := uint32(2166136261)
	const prime32 = uint32(16777619)
	for i := 0; i < len(key); i++ {
		hash *= prime32
		hash ^= uint32(key[i])
	}
	return hash
}

func (m shardMap) Set(key string, value interface{}, expire ...time.Duration) {
	m.GetShard(key).Set(key, value, expire...)
}

func (m shardMap) Get(key string) (interface{}, bool) {
	return m.GetShard(key).Get(key)
}

func (m shardMap) Delete(keys ...string) {
	for _, key := range keys {
		m.GetShard(key).Delete(key)
	}
}

func (m shardMap) Exists(key string) bool {
	return m.GetShard(key).Exists(key)
}

// Count 数据库元素个数
func (m shardMap) Count() uint64 {
	var count uint64
	var wg sync.WaitGroup

	wg.Add(SHARD_COUNT)
	for i := 0; i < SHARD_COUNT; i++ {
		go func(idx int) {
			defer wg.Done()
			atomic.AddUint64(&count, m[idx].Count())
		}(i)
	}
	wg.Wait()

	return count
}

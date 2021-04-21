package gocache

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

const (
	// 定期扫描 key 的数量
	keyCount int = 1000
)

var cache *Cache
var once sync.Once

type deleteElementKeys struct {
	dbIndex uint32
	keys    []string
}

// Cache
type Cache struct {
	mu sync.Mutex
	// db index
	dbindex uint32

	// bucket 数量
	dbcount uint32

	// daemon db index
	current_db uint32

	// 每秒清扫过期键对的次数
	hz uint32

	expireDb []shardMap
	expireCh chan deleteElementKeys

	db []shardMap
}

// New
func New(opts ...Option) *Cache {
	once.Do(func() {
		options := CacheOption{
			dbcount: 1,
			dbindex: 0,
		}

		for _, fn := range opts {
			fn(&options)
		}

		if options.dbindex >= options.dbcount {
			panic(fmt.Sprintf("dbindex: %d out of range dbcount:%d", options.dbindex, options.dbcount))
		}

		cache = &Cache{
			dbcount:  options.dbcount,
			dbindex:  options.dbindex,
			expireDb: make([]shardMap, options.dbcount),
			expireCh: make(chan deleteElementKeys, 100),
			db:       make([]shardMap, options.dbcount),
		}

		// 初始化所有数据库
		for i := 0; i < int(options.dbcount); i++ {
			cache.expireDb[i] = NewShardMap()
			cache.db[i] = NewShardMap()
		}

		// 清理监控服务
		go cache.daemon()
	})

	return cache
}

// Set
func (c *Cache) Set(key string, value interface{}, expire ...time.Duration) {
	// expireDb
	if len(expire) > 0 {
		d := expire[0]
		if d > 0 {
			c.expireDb[c.dbindex].Set(key, value, expire...)
		}
	} else {
		c.expireDb[c.dbindex].Delete(key)
	}

	c.db[c.dbindex].Set(key, value, expire...)
}

// Get
func (c *Cache) Get(key string) (interface{}, bool) {
	value, ok := c.db[c.dbindex].Get(key)
	if !ok {
		c.mu.Lock()
		c.expireDb[c.dbindex].Delete(key)
		c.mu.Unlock()
	}

	return value, ok
}

// Delete
func (c *Cache) Delete(key string) {
	c.db[c.dbindex].Delete(key)

	// delete from expireDb
	c.mu.Lock()
	c.expireDb[c.dbindex].Delete(key)
	c.mu.Unlock()
}

// Exists
func (c *Cache) Exists(key string) bool {
	return c.db[c.dbindex].Exists(key)
}

// daemon 过期键值对清剿进程
func (c *Cache) daemon() {
	// 过期键清理
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go c.purgeWorker(ctx)

	// dynamic-hz
	if c.hz == 0 {
		c.hz = 10
	}

	// 每秒清扫次数
	timelimit := time.Second / time.Duration(c.hz)
	tick := time.NewTicker(timelimit)
	defer tick.Stop()

	for {
		select {
		case <-tick.C:
			c.purgeDB()
		}
	}
}

// purgeExpireKey
func (c *Cache) purgeDB() {
	c.mu.Lock()
	dbindex := c.current_db
	c.current_db++
	if c.current_db > c.dbcount-1 {
		c.current_db = 0
	}
	c.mu.Unlock()

	for _, bucket := range c.expireDb[dbindex] {
		shouldDeleteElement := make([]string, 0)

		bucket.mu.RLock()
		i := 0
		for k, v := range bucket.data {
			if i >= keyCount {
				return
			}

			if v.expired() {
				i++
				shouldDeleteElement = append(shouldDeleteElement, k)
			}
		}
		bucket.mu.RUnlock()

		if len(shouldDeleteElement) > 0 {
			elements := deleteElementKeys{
				dbIndex: dbindex,
				keys:    shouldDeleteElement,
			}

			c.expireCh <- elements
		}
	}
}

func (c *Cache) purgeWorker(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println(ctx.Err())
			return
		case elems := <-c.expireCh:
			c.expireDb[elems.dbIndex].Delete(elems.keys...)
			c.db[elems.dbIndex].Delete(elems.keys...)
		}
	}
}

// SelectDB 选择数据库
func (c *Cache) SelectDB(index uint32) error {
	if index >= c.dbcount {
		return errors.New("db number not exists")
	}
	c.dbindex = index

	return nil
}

// GetDB 当前数据库编号
func (c *Cache) GetDB() uint32 {
	return c.dbindex
}

// FlushDB 清空一个数据库
func (c *Cache) FlushDB() {
	c.db[c.dbindex] = NewShardMap()
}

// Stats
func (c *Cache) Stats() []Stats {
	statsMap := make([]Stats, c.dbcount)
	for i := 0; i < int(c.dbcount); i++ {
		statsMap[i] = Stats{
			DbIndex: uint32(i),
			Count:   c.db[i].Count(),
		}
	}

	return statsMap
}

// Stats
type Stats struct {
	DbIndex uint32
	// 元素个数
	Count uint64
}

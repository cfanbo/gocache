package gocache

import (
	"sync"
	"time"
)

// CacheBucket
type CacheBucket struct {
	mu sync.RWMutex

	data map[string]*CacheValue
}

func NewCacheBucket() *CacheBucket {
	return &CacheBucket{
		data: make(map[string]*CacheValue),
	}
}

// Set
func (c *CacheBucket) Set(key string, value interface{}, expiration ...time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	var d time.Duration = 0
	if len(expiration) > 0 {
		d = expiration[0]
	}

	data, ok := c.data[key]
	if ok {
		// 直接修改
		data.value = value
		if d == 0 {
			data.expire = time.Time{}
		} else {
			data.expire = time.Now().Add(d)
		}
	} else {
		c.data[key] = NewCacheValue(value, d)
	}

	// c.data[key] = NewCacheValue(value, d)
}

// Get 惰性删除
func (c *CacheBucket) Get(key string) (interface{}, bool) {
	c.mu.RLock()

	data, ok := c.data[key]
	// notfound
	if !ok {
		c.mu.RUnlock()
		return nil, false
	}

	// expired
	if data.expired() {
		c.mu.RUnlock()
		// 删除节省空间
		c.Delete(key)
		return nil, false
	}

	c.mu.RUnlock()

	return data.value, ok
}

// Delete
func (c *CacheBucket) Delete(keys ...string) {
	if len(keys) == 0 {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()

	for _, v := range keys {
		delete(c.data, v)
	}
}

// Exists
func (c *CacheBucket) Exists(key string) bool {
	_, ok := c.Get(key)

	return ok
}

// Count 桶内元素个数
func (c *CacheBucket) Count() uint64 {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return uint64(len(c.data))
}

// Expired
func (c *CacheBucket) Expired(key string) bool {
	c.mu.RLock()
	defer c.mu.Unlock()

	data, ok := c.data[key]
	if !ok {
		return true
	}
	return data.expired()
}

// CacheValue
type CacheValue struct {
	value  interface{}
	expire time.Time
}

// NewCacheValue
func NewCacheValue(val interface{}, expiration time.Duration) *CacheValue {
	return &CacheValue{
		value: val,
		expire: func(t time.Duration) time.Time {
			if t == 0 {
				return time.Time{}
			} else {
				return time.Now().Add(t)
			}
		}(expiration),
	}
}

// 值过期
func (v *CacheValue) expired() bool {
	if v.expire.IsZero() {
		return false
	}

	return time.Now().After(v.expire)
}

func (v *CacheValue) String() string {
	return v.value.(string)
}

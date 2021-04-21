package gocache_test

import (
	"testing"
	"time"

	"github.com/cfanbo/gocache"
)

var cache *gocache.Cache
var cb *gocache.CacheBucket

func TestMain(m *testing.M) {
	cache = gocache.New(gocache.WithDbNums(10))
	cb = gocache.NewCacheBucket()
	m.Run()
}

func TestPcSetAndGet(t *testing.T) {
	cache.Set("foo", "bar")
	val, ok := cache.Get("foo")
	if !ok {
		t.Fatal("读取失败")
	}

	if val.(string) != "bar" {
		t.Fatal("缓存读取失败")
	}
}

func TestPcDelete(t *testing.T) {
	cache.Set("foo", "bar")
	cache.Delete("foo")

	_, ok := cache.Get("foo")
	if ok {
		t.Fatal("删除出错")
	}
}

func TestPcExpire(t *testing.T) {
	m, _ := time.ParseDuration("-10m")

	// 已过期
	cache.Set("foo", "bar", m)
	_, ok := cache.Get("foo")
	if ok {
		t.Fatal("读取到了过期数据")
	}

	// 未过期
	cache.Set("foo", "bar", time.Second)
	_, ok = cache.Get("foo")
	if !ok {
		t.Fatal("数据过期出错")
	}
}

func TestPcSelectDBAndGetDB(t *testing.T) {
	cache.SelectDB(2)
	if cache.GetDB() != 2 {
		t.Error("TestPcSelectDBAndGetDB")
	}
}

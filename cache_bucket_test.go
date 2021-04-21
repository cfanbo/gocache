package gocache_test

import (
	"testing"
	"time"
)

func TestCbStoreAndRead(t *testing.T) {
	cb.Set("foo", "bar")
	val, ok := cb.Get("foo")
	if !ok {
		t.Fatal("读取失败")
	}

	if val.(string) != "bar" {
		t.Fatal("缓存读取失败")
	}
}

func TestCbDelete(t *testing.T) {
	cb.Set("foo", "bar")
	cb.Delete("foo")

	_, ok := cb.Get("foo")
	if ok {
		t.Fatal("删除出错")
	}
}

func TestCbExpire(t *testing.T) {
	m, _ := time.ParseDuration("-10s")

	// 过期时间
	cb.Set("foo", "bar", m)
	_, ok := cb.Get("foo")
	if ok {
		t.Fatal("读取到了过期数据")
	}

	// 未过期时间
	m2, _ := time.ParseDuration("10s")
	cb.Set("foo", "bar", m2)
	_, ok = cb.Get("foo")
	if !ok {
		t.Fatal("未读取到未过期的数据")
	}
}

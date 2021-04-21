package gocache_test

import (
	"testing"

	"github.com/cfanbo/gocache"
)

func TestShardMapSetAndGet(t *testing.T) {
	cMap := gocache.NewShardMap()
	cMap.Set("foo", "bar")
	value, ok := cMap.Get("foo")
	if !ok {
		t.Fatal("TestShardMapSet 存储失败")
	}

	if value.(string) != "bar" {
		t.Fatal("TestShardMapSet 读取出错")
	}
}

func TestShardMapSetEx(t *testing.T) {
	cMap := gocache.NewShardMap()
	cMap.Set("foo", "bar")
	if !cMap.Exists("foo") {
		t.Fatal("TestShardMapSetEx")
	}
}

func TestShardMapDelete(t *testing.T) {
	cMap := gocache.NewShardMap()
	cMap.Set("foo", "bar")
	cMap.Delete("foo")
	if cMap.Exists("foo") {
		t.Fatal("TestShardMapDelete 存储失败")
	}
}

func TestShardMapExist(t *testing.T) {
	cMap := gocache.NewShardMap()
	if cMap.Exists("foo") {
		t.Fatal("TestShardMapDelete 存储失败")
	}
}

func TestShardMapCount(t *testing.T) {
	cMap := gocache.NewShardMap()
	cMap.Set("foo", "a")
	cMap.Set("bar", "b")
	if cMap.Count() != 2 {
		t.Fatal("TestShardMapCount")
	}
}

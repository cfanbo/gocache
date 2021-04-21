## 高性能缓存服务

[![Go](https://github.com/cfanbo/gocache/actions/workflows/build.yml/badge.svg)](https://github.com/cfanbo/gocache/actions/workflows/build.yml) [![Build Status](https://travis-ci.com/cfanbo/gocache.svg?branch=main)](https://travis-ci.com/cfanbo/gocache) [![Go Reference](https://pkg.go.dev/badge/github.com/cfanbo/gocache.svg)](https://pkg.go.dev/github.com/cfanbo/gocache)

gocache 是一款高性能进程级缓存服务，在一定场景下可以用来代替redis，减少三方缓存依赖。



## 功能
* 支持多个数据库，类似redis，大部分命令相同
* 每个库底层自适应分区，提高数据量过多时的并发读写性能
* 支持键值对的自动过期机制

## 安装
import the package
```
import (
	"github.com/cfanbo/gocache
)
```

install the package
```
go get github.com/cfanbo/gocache
```
the package is now imported under the "gocache" namespace

## 使用方法
```
package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/cfanbo/gocache"
)

func main() {
	cache := gocache.New(gocache.WithDbNums(5), gocache.WithDefaultDb(2), gocache.WithHz(1000))
	// 普通键值对
	fmt.Println("\n写入普通内容 ---------------")
	cache.Set("foo", "bar")

	// 带过期时间的键值对
	fmt.Println("\n写入带过期时间的内容 ---------------")
	m, _ := time.ParseDuration("10m")
	cache.Set("foo", "bar", m)
	val, ok := cache.Get("foo")
	if !ok {
		fmt.Println("未找到数据")
		os.Exit(0)
	} else {
		fmt.Println(val)
	}

	oldDB := cache.GetDB()
	fmt.Println("当前库编号:", oldDB)
	// 在当前库写入10000
	for i := 0; i < 10000; i++ {
		cache.Set(strconv.Itoa(i), i+10)
	}

	fmt.Println(cache.Get("10"))
	// 切换数据库
	fmt.Println("\n切换当前数据库，并写入内容 ---------------")
	if err := cache.SelectDB(3); err != nil {
		panic(err)
	}
	fmt.Println("当前库编号:", cache.GetDB())

	for i := 0; i < 5000; i++ {
		cache.Set(strconv.Itoa(i), i)
	}
	fmt.Println(cache.Get("10"))
	for _, stat := range cache.Stats() {
		fmt.Println("dbindex:", stat.DbIndex, "count:", stat.Count)
	}

	// 清空当前数据库
	fmt.Println("\n清空当前数据库 ---------------")
	cache.FlushDB()
	for _, stat := range cache.Stats() {
		fmt.Println("dbindex:", stat.DbIndex, "count:", stat.Count)
	}
	// 切回原来数据库
	cache.SelectDB(oldDB)
	fmt.Println(cache.Get("10"))

}

```

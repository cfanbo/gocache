package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/cfanbo/gocache"
)

func main() {
	// 直接使用
	cache := gocache.New(gocache.WithDbNums(5), gocache.WithDefaultDb(2), gocache.WithHz(1000))

	fmt.Println("\n写入普通内容 ---------------")
	cache.Set("foo", "bar")

	fmt.Println("\n写入带过期时间的内容 ---------------")
	m, _ := time.ParseDuration("2s")
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
	cache.SelectDB(2)
	for i := 0; i < 10000; i++ {
		rand.Seed(int64(i + 1))
		d := time.Second * (time.Duration(rand.Intn(5)) + 1)
		cache.Set("key"+strconv.Itoa(i+10000), "value"+strconv.Itoa(i+10000)+"abc", d)
	}
	// purge
	fmt.Println("purge begin ------------")
	cache.SelectDB(1)
	for i := 0; i < 10000; i++ {
		rand.Seed(int64(i + 1))
		d := time.Second * (time.Duration(rand.Intn(7)) + 1)
		cache.Set("key"+strconv.Itoa(i), "value"+strconv.Itoa(i)+"abc", d)
	}
	fmt.Println("purge end ------------")

	time.Sleep(time.Second * 10)

}

// 文件名：rwmutex_demo.go
package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	var rwMutex sync.RWMutex
	var wg sync.WaitGroup
	data := "初始数据"

	// 1. 模拟 5 个只读的协程
	for i := 1; i <= 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			
			rwMutex.RLock() // ★ 读锁：大家都可以同时拿这个锁
			fmt.Printf("读者 %d 正在读取: %s\n", id, data)
			time.Sleep(100 * time.Millisecond) // 模拟看书耗时
			rwMutex.RUnlock()
		}(i)
	}

	// 2. 模拟 1 个写的协程
	time.Sleep(50 * time.Millisecond) // 稍微等一下，让读者先拿到锁
	wg.Add(1)
	go func() {
		defer wg.Done()
		
		rwMutex.Lock() // ★ 写锁：此时如果有读者，写锁必须等所有读者走光才能进
		fmt.Println("--- 写者：我要修改数据了，你们先别看了 ---")
		time.Sleep(200 * time.Millisecond)
		data = "修改后的新数据"
		fmt.Println("--- 写者：修改完毕 ---")
		rwMutex.Unlock()
	}()

	wg.Wait()
}
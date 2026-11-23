// 文件名：safemap_demo.go
package main

import (
	"fmt"
	"sync"
)

// 把普通的 map 和读写锁封装进一个结构体
type SafeMap struct {
	mu sync.RWMutex       // 读写锁
	data map[string]int   // 原生的 map，把它藏起来，不让外部直接碰
}

// 提供一个安全的写入方法
func (s *SafeMap) Set(key string, value int) {
	s.mu.Lock()         // 写锁，独占
	defer s.mu.Unlock()
	s.data[key] = value
}

// 提供一个安全的读取方法
func (s *SafeMap) Get(key string) int {
	s.mu.RLock()        // 读锁，可并发
	defer s.mu.RUnlock()
	return s.data[key]
}

func main() {
	// 初始化安全柜子
	safeMap := &SafeMap{
		data: make(map[string]int),
	}

	var wg sync.WaitGroup

	// 模拟 1000 个人同时往柜子里存东西
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			safeMap.Set(fmt.Sprintf("user_%d", id), id)
		}(i)
	}

	wg.Wait()
	
	// 验证数据是否正确存入
	fmt.Println("user_999 的数据是:", safeMap.Get("user_999"))
	fmt.Println("数据安全存入，没有引发并发崩溃！")
}
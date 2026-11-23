package main

import (
	"fmt"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	var mu sync.Mutex // 这就是那把锁
	
	// 假设这是大家共享的公共账户余额
	money := 100

	// 模拟 1000 个人同时想往余额里存 1 块钱
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			
			mu.Lock() // 进门，锁门！（后面的协程只能排队等）
			
			// 这段代码被锁保护，同一时间只能有一个人执行，绝对不会乱
			// （如果你把锁注释掉，最后 money 一定会小于 1100，因为发生了并发冲突）
			money = money + 1 
			
			mu.Unlock() // 出门，开锁！
		}()
	}

	wg.Wait()
	fmt.Println("1000人存钱后，最终余额:", money) // 输出绝对是 1100
}
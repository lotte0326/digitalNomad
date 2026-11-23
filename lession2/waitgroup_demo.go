// 文件名：waitgroup_demo.go
package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup

	for i := 1; i <= 3; i++ {
		wg.Add(1) // ★ 必须在这加！告诉老板还有活
		
		go func(id int, wg *sync.WaitGroup) { // 必须传指针！
			defer wg.Done() // ★ 干完活必须减一
			
			fmt.Printf("员工 %d 正在搬货...\n", id)
			time.Sleep(time.Duration(id) * 300 * time.Millisecond)
			fmt.Printf("员工 %d 搬完了！\n", id)
		}(i, &wg)
	}

	fmt.Println("老板：等你们搬完我再走...")
	wg.Wait() // 老板阻塞在这里
	fmt.Println("老板：人都齐了，关门下班！")
}
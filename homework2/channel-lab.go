package main

import (
	"fmt"
	"sync"
)

// ============================================================
// 题目一：无缓冲通道实现两个协程之间的通信
// ============================================================

// producer 生成 1 到 10 的整数，发送到通道
func producer(ch chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 1; i <= 1000; i++ {
		fmt.Printf("[生产者] 发送: %d\n", i)
		ch <- i // 发送到通道（无缓冲：会阻塞直到有人接收）
	}
	close(ch) // 生产完毕关闭通道，通知消费者
}

// consumer 从通道接收整数并打印
func consumer(ch <-chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	// range 会自动在通道关闭且数据取完后退出循环
	for v := range ch {
		fmt.Printf("  [消费者] 收到: %d\n", v)
	}
}

// runUnbufferedDemo 题目一演示
func runUnbufferedDemo() {
	fmt.Println("========== 题目一：无缓冲通道通信 ==========")

	ch := make(chan int) // 无缓冲通道
	var wg sync.WaitGroup

	wg.Add(2)
	go producer(ch, &wg)
	go consumer(ch, &wg)

	wg.Wait()
	fmt.Println("无缓冲通道演示完毕")
	fmt.Println()
}

// ============================================================
// 题目二：带缓冲通道实现生产者-消费者
// ============================================================

// bufferedProducer 向带缓冲通道发送 100 个整数
func bufferedProducer(ch chan<- int, n int, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 1; i <= n; i++ {
		ch <- i // 缓冲未满时不会阻塞
		fmt.Printf("[生产者] 发送: %d\n", i)
	}
	close(ch) // 发送完毕关闭通道
}

// bufferedConsumer 从带缓冲通道接收整数并打印
func bufferedConsumer(ch <-chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	count := 0
	for v := range ch {
		fmt.Printf("  [消费者] 收到: %d\n", v)
		count++
		if count%20 == 0 { // 每 20 个换行，输出更整齐
			fmt.Println()
		}
	}
	if count%20 != 0 {
		fmt.Println()
	}
	fmt.Printf("  [消费者] 共接收 %d 个整数\n", count)
}

// runBufferedDemo 题目二演示
func runBufferedDemo() {
	fmt.Println("========== 题目二：带缓冲通道（生产者-消费者） ==========")

	// 创建缓冲大小为 10 的通道
	const bufSize = 10
	ch := make(chan int, bufSize)
	fmt.Printf("通道缓冲大小: %d, 容量: %d, 当前长度: %d\n\n",
		bufSize, cap(ch), len(ch))

	var wg sync.WaitGroup
	wg.Add(2)
	go bufferedProducer(ch, 1000, &wg)
	go bufferedConsumer(ch, &wg)
	wg.Wait()

	fmt.Println("带缓冲通道演示完毕")
	fmt.Println()
}

// ============================================================
// 补充演示：缓冲机制的直观对比（可选，帮助理解）
// ============================================================

// runBufferMechanismDemo 演示缓冲通道"未满不阻塞"的特性
func runBufferMechanismDemo() {
	fmt.Println("========== 补充：缓冲机制直观演示 ==========")

	// 缓冲大小为 3 的通道
	ch := make(chan int, 3)

	// 连续发送 3 个，因缓冲未满，不会阻塞
	ch <- 1
	ch <- 2
	ch <- 3
	fmt.Printf("发送 3 个后: len=%d, cap=%d（未阻塞）\n", len(ch), cap(ch))

	// 第 4 个会阻塞——所以放到协程里发，主协程去接收
	done := make(chan struct{})
	go func() {
		fmt.Println("尝试发送第 4 个...（缓冲已满，阻塞，等待消费者）")
		ch <- 4
		fmt.Println("第 4 个发送成功（消费者腾出了空间）")
		close(done)
	}()

	// 接收一个，腾出空间
	v := <-ch
	fmt.Printf("主协程接收: %d, 此时 len=%d\n", v, len(ch))
	<-done

	// 清空剩余
	for len(ch) > 0 {
		<-ch
	}
	fmt.Printf("清空后: len=%d\n", len(ch))
	fmt.Println()
}

// ============================================================
// main 入口
// ============================================================
func main() {
	runUnbufferedDemo()
	runBufferedDemo()
	runBufferMechanismDemo()
}
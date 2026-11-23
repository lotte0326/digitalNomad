package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// ============================================================================
// 知识点 1：Goroutine 创建和使用
// ============================================================================
// 【底层机制】Goroutine 是 Go 运行时（Runtime）管理的轻量级线程。初始栈仅 2KB，
//             按需扩缩容。Go 调度器（GMP 模型）负责将成千上万个 Goroutine 调度到少量的 OS 线程上。
// 【语法】go 函数名(参数)
// 【防坑】主 Goroutine（main 函数）退出时，所有其他 Goroutine 会被强制终止，不会等待。
//         所以必须用 WaitGroup 或 Channel 阻塞主程序。

// 模拟抓取一个网页
func fetchPage(ctx context.Context, url string, results chan<- string, wg *sync.WaitGroup) {
	defer wg.Done() // 知识点 4：任务完成减一

	// 模拟网络延迟（随机 100-500 毫秒）
	delay := time.Duration(rand.Intn(400)+100) * time.Millisecond

	// 知识点 3：Select 语句
	// 在抓取过程中，随时监听是否被取消（超时或主动取消）
	select {
	case <-ctx.Done():
		// 如果被取消，直接放弃抓取
		fmt.Printf("  [取消] %s 抓取被中断: %v\n", url, ctx.Err())
		return
	case <-time.After(delay): // 模拟网络请求完成
		// 抓取成功，生成结果
		result := fmt.Sprintf("[成功] %s (耗时 %v)", url, delay)
		
		// 知识点 2：Channel 通信机制
		// 将结果发送到 Channel。这里使用带缓冲的 channel，防止因为没人接收而阻塞
		results <- result 
	}
}

// ============================================================================
// 知识点 6：生产者-消费者模式
// ============================================================================
// 在这里，fetchPage 是“生产者”（生产抓取结果），main 函数中的 for 循环是“消费者”。

func main() {
	fmt.Println("=== 并发爬虫启动 ===")
	
	// 知识点 5：Context 上下文控制
	// 设置一个总超时时间为 1 秒。如果 1 秒内没抓完，全部强制停止。
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel() // 最佳实践：无论正常结束还是超时，最后都要调用 cancel 释放资源，防止内存泄漏

	// 知识点 2：Channel 通信机制
	// 创建一个带缓冲的 Channel，容量 10，用来存放抓取结果
	results := make(chan string, 10)

	// 知识点 4：WaitGroup 同步
	var wg sync.WaitGroup

	// 模拟 5 个要抓取的网页
	urls := []string{
		"https://go.dev", "https://google.com", "https://github.com",
		"https://stackoverflow.com", "https://golang.org",
	}

	// 启动 Goroutine
	for _, url := range urls {
		wg.Add(1) // 每启动一个任务，计数器加 1
		// 知识点 1：使用 go 关键字启动并发任务
		go fetchPage(ctx, url, results, &wg)
	}

	// 【关键步骤】启动一个独立的 Goroutine 用来等待所有任务完成并关闭 Channel
	// 为什么不能在主协程里直接 wg.Wait() 然后 close(results)？
	// 因为如果主协程阻塞了，下面的 for range 循环就无法消费 results，导致生产者阻塞死锁。
	go func() {
		wg.Wait()      // 等待所有 fetchPage 执行完毕
		close(results) // 所有任务结束后，关闭 Channel，通知消费者没有新数据了
	}()

	fmt.Println("等待抓取结果...")
	
	// 消费结果
	// 使用 for range 持续从 Channel 读取数据，直到 Channel 被 close。
	// 注意：如果 for range 写在主协程里，它可以自然阻塞主程序，替代直接调用 wg.Wait()。
	for res := range results {
		fmt.Println("  =>", res)
	}

	fmt.Println("=== 所有任务处理完毕 ===")
}
package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// 定义订单结构体
type Order struct {
	ID    int
	Price float64
}

// ==========================================
// 知识点 1：Goroutine 创建和使用
// ==========================================
// 只要在函数调用前加上 go 关键字，就能启动一个 Goroutine（轻量级线程）。
// 启动开销极小，一个 Go 程序可以轻松启动数十万个 Goroutine。

// ==========================================
// 知识点 6：生产者-消费者模式
// ==========================================
// 生产者：不断生成订单并放入 Channel（队列）
func producer(ctx context.Context, orderChan chan<- Order, wg *sync.WaitGroup) {
	defer wg.Done() // 知识点 4：任务完成时通知 WaitGroup，等同于 wg.Add(-1)
	defer close(orderChan) // 生产者结束后，关闭 Channel 通知消费者没有新数据了

	orderID := 1
	for {
		select {
		case <-ctx.Done():
			// 知识点 5：Context 上下文控制
			// 如果外部发出了取消信号（或超时），立即停止生产
			fmt.Println("【生产者】收到取消信号，停止生产")
			return
		default:
			// 模拟生产耗时
			time.Sleep(200 * time.Millisecond)
			order := Order{
				ID:    orderID,
				Price: float64(rand.Intn(100)),
			}
			// 知识点 2：Channel 通信机制
			// 将订单发送到 Channel 中，如果 Channel 满了，这里会阻塞
			orderChan <- order
			fmt.Printf("【生产者】生成了订单 #%d\n", orderID)
			orderID++
		}
	}
}

// 消费者：从 Channel 中取出订单并处理
func consumer(ctx context.Context, id int, orderChan <-chan Order, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			// 收到取消信号，停止消费
			fmt.Printf("【消费者 %d】收到取消信号，下班了\n", id)
			return
		case order, ok := <-orderChan:
			// 知识点 3：Select 语句
			// select 会一直等待，直到某个 case 中的 Channel 操作准备就绪。
			// 这里如果 orderChan 被关闭且没有数据，ok 会变成 false。
			if !ok {
				fmt.Printf("【消费者 %d】订单通道已关闭，下班了\n", id)
				return
			}
			// 模拟处理订单耗时
			time.Sleep(300 * time.Millisecond)
			fmt.Printf("  【消费者 %d】处理了订单 #%d (金额: %.2f)\n", id, order.ID, order.Price)
		}
	}
}

func main() {
	fmt.Println("=== 订单处理系统启动 ===")

	// 知识点 5：Context 上下文控制
	// 创建一个带有取消功能的 Context，用于统一控制所有的 Goroutine 退出
	ctx, cancel := context.WithCancel(context.Background())

	// 知识点 2：Channel 通信机制
	// 创建一个带缓冲的 Channel，容量为 5。缓冲可以缓解生产和消费速度不一致的问题。
	orderChan := make(chan Order, 5)

	// 知识点 4：WaitGroup 同步
	// WaitGroup 用于等待一组 Goroutine 全部执行完毕。相当于计数器。
	var wg sync.WaitGroup

	// 启动 1 个生产者
	wg.Add(1)
	go producer(ctx, orderChan, &wg)

	// 启动 3 个消费者
	for i := 1; i <= 3; i++ {
		wg.Add(1)
		// 知识点 1：启动 Goroutine
		go consumer(ctx, i, orderChan, &wg)
	}

	// 让系统运行 2 秒钟
	time.Sleep(2 * time.Second)

	fmt.Println("\n=== 主程序准备结束，发出取消信号 ===")
	// 触发 Context 的取消信号
	// 此时，生产者通过 ctx.Done() 收到信号，停止生产并关闭 orderChan。
	// 消费者通过 ctx.Done() 收到信号，或者通过 orderChan 的关闭状态（!ok）退出。
	cancel()

	// 知识点 4：WaitGroup 同步
	// 阻塞等待上面所有的 Goroutine 都执行完毕（接收到 cancel 信号并执行到 defer wg.Done()）
	wg.Wait()
	fmt.Println("=== 所有 Goroutine 已安全退出，程序结束 ===")
}
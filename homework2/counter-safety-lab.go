package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

// ============================================================
// 题目一：使用 sync.Mutex 保护共享计数器
// ============================================================

// MutexCounter 用互斥锁保护的计数器
type MutexCounter struct {
	mu    sync.Mutex // 互斥锁，保护 value
	value int
}

// Inc 加锁递增，保证并发安全
func (c *MutexCounter) Inc() {
	c.mu.Lock()
	defer c.mu.Unlock() // defer 确保即使 panic 也能解锁
	c.value++
}

// Value 读取当前值（读也需加锁，避免读到中间状态）
func (c *MutexCounter) Value() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.value
}

// runMutexDemo 题目一演示
func runMutexDemo() {
	fmt.Println("========== 题目一：sync.Mutex 保护计数器 ==========")

	const (
		goroutines = 10
		increments = 1000
	)

	var counter MutexCounter
	var wg sync.WaitGroup

	// 启动 10 个协程，每个递增 1000 次
	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < increments; j++ {
				counter.Inc()
			}
		}()
	}

	wg.Wait()
	fmt.Printf("Mutex 计数器最终值: %d（期望: %d）\n", counter.Value(), goroutines*increments)
	if counter.Value() == goroutines*increments {
		fmt.Println("✅ 并发安全，结果正确")
	} else {
		fmt.Println("❌ 结果错误，存在数据竞争")
	}
	fmt.Println()
}

// ============================================================
// 题目二：使用 sync/atomic 实现无锁计数器
// ============================================================

// AtomicCounter 无锁计数器，使用原子操作
type AtomicCounter struct {
	value int64 // 必须用 int64，配合 atomic 函数（或 int32/uint64 等）
}

// Inc 原子递增，无锁、并发安全
func (c *AtomicCounter) Inc() {
	atomic.AddInt64(&c.value, 1)
}

// Value 原子读取当前值
func (c *AtomicCounter) Value() int64 {
	return atomic.LoadInt64(&c.value)
}

// runAtomicDemo 题目二演示
func runAtomicDemo() {
	fmt.Println("========== 题目二：sync/atomic 无锁计数器 ==========")

	const (
		goroutines = 10
		increments = 1000
	)

	var counter AtomicCounter
	var wg sync.WaitGroup

	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < increments; j++ {
				counter.Inc()
			}
		}()
	}

	wg.Wait()
	final := counter.Value()
	fmt.Printf("Atomic 计数器最终值: %d（期望: %d）\n", final, goroutines*increments)
	if final == int64(goroutines*increments) {
		fmt.Println("✅ 无锁并发安全，结果正确")
	} else {
		fmt.Println("❌ 结果错误")
	}
	fmt.Println()
}

// ============================================================
// 补充对比：不加任何保护的计数器（反例，演示数据竞争）
// ============================================================

// UnsafeCounter 未加保护的计数器，仅用于对比
type UnsafeCounter struct {
	value int
}

func (c *UnsafeCounter) Inc() {
	c.value++ // 读-改-写 三步非原子，并发下会出错
}

// runUnsafeDemo 反例演示（结果通常小于期望值，且每次不同）
func runUnsafeDemo() {
	fmt.Println("========== 补充：不加保护的反例（数据竞争） ==========")

	const (
		goroutines = 10
		increments = 1000
	)

	var counter UnsafeCounter
	var wg sync.WaitGroup

	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < increments; j++ {
				counter.Inc()
			}
		}()
	}

	wg.Wait()
	fmt.Printf("无保护计数器最终值: %d（期望: %d）\n", counter.value, goroutines*increments)
	fmt.Println("⚠️  通常小于期望值，因为 value++ 是非原子的读-改-写操作")
	fmt.Println()
}

// ============================================================
// 性能对比（可选）：Mutex vs Atomic
// ============================================================

// runBenchmarkDemo 粗略对比两种方式耗时
func runBenchmarkDemo() {
	fmt.Println("========== 补充：Mutex vs Atomic 性能粗略对比 ==========")

	const (
		goroutines = 10
		increments = 100000
	)

	var wg sync.WaitGroup

	// Mutex
	var mc MutexCounter
	start := make(chan struct{})
	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start // 同时起跑
			for j := 0; j < increments; j++ {
				mc.Inc()
			}
		}()
	}
	close(start)
	wg.Wait()
	fmt.Printf("Mutex 耗时示例（值=%d）\n", mc.Value())

	// Atomic
	var ac AtomicCounter
	var wg2 sync.WaitGroup
	start2 := make(chan struct{})
	for i := 0; i < goroutines; i++ {
		wg2.Add(1)
		go func() {
			defer wg2.Done()
			<-start2
			for j := 0; j < increments; j++ {
				ac.Inc()
			}
		}()
	}
	close(start2)
	wg2.Wait()
	fmt.Printf("Atomic 耗时示例（值=%d）\n", ac.Value())
	fmt.Println("（实际耗时受机器负载影响，可用 go test -bench 精确对比）")
	fmt.Println()
}

// ============================================================
// main 入口
// ============================================================
func main() {
	runMutexDemo()
	runAtomicDemo()
	runUnsafeDemo()
	runBenchmarkDemo()
}
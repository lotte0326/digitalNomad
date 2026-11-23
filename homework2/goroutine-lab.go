package main

import (
	"fmt"
	"sync"
	"time"
)

// ============================================================
// 题目一：使用 go 关键字启动两个协程，分别打印奇数和偶数
// ============================================================

// printOdd 打印 1 到 10 之间的奇数
func printOdd(wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 1; i <= 10; i++ {
		if i%2 != 0 {
			fmt.Printf("[奇数协程] %d\n", i)
		}
	}
}

// printEven 打印 2 到 10 之间的偶数
func printEven(wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 2; i <= 10; i++ {
		if i%2 == 0 {
			fmt.Printf("[偶数协程] %d\n", i)
		}
	}
}

// runOddEvenDemo 启动两个协程并发执行
func runOddEvenDemo() {
	fmt.Println("========== 题目一：奇偶打印协程 ==========")
	var wg sync.WaitGroup

	wg.Add(2)
	go printOdd(&wg)  // 使用 go 关键字启动协程
	go printEven(&wg) // 使用 go 关键字启动协程

	wg.Wait() // 等待两个协程执行完成
	fmt.Println("奇偶打印协程执行完毕")
	fmt.Println()
}

// ============================================================
// 题目二：任务调度器，并发执行任务并统计执行时间
// ============================================================

// Task 表示一个任务：名称 + 执行函数
type Task struct {
	Name string
	Fn   func() error
}

// TaskResult 保存任务执行结果
type TaskResult struct {
	Name     string
	Duration time.Duration
	Err      error
}

// Scheduler 任务调度器
type Scheduler struct {
	tasks []Task
}

// NewScheduler 创建调度器
func NewScheduler(tasks ...Task) *Scheduler {
	return &Scheduler{tasks: tasks}
}

// Run 并发执行所有任务，并返回每个任务的执行时间
func (s *Scheduler) Run() []TaskResult {
	var wg sync.WaitGroup
	results := make([]TaskResult, len(s.tasks))
	// 使用带缓冲的 channel 收集结果（也可直接写切片，因为每个协程写不同下标）
	resultCh := make(chan TaskResult, len(s.tasks))

	for i, task := range s.tasks {
		wg.Add(1)
		// 用局部变量捕获 task，避免闭包引用循环变量问题
		t := task
		idx := i
		go func() {
			defer wg.Done()

			start := time.Now()
			err := t.Fn()
			duration := time.Since(start)

			res := TaskResult{
				Name:     t.Name,
				Duration: duration,
				Err:      err,
			}
			results[idx] = res
			resultCh <- res
		}()
	}

	wg.Wait()
	close(resultCh)

	// 也可以从 channel 读取（这里已同时写入 results 切片）
	for range resultCh {
		// 仅为演示 channel 使用，可省略
	}

	return results
}

// runSchedulerDemo 演示任务调度器
func runSchedulerDemo() {
	fmt.Println("========== 题目二：任务调度器 ==========")

	tasks := []Task{
		{
			Name: "任务A(耗时100ms)",
			Fn: func() error {
				time.Sleep(100 * time.Millisecond)
				fmt.Println("   -> 任务A 执行完成")
				return nil
			},
		},
		{
			Name: "任务B(耗时200ms)",
			Fn: func() error {
				time.Sleep(200 * time.Millisecond)
				fmt.Println("   -> 任务B 执行完成")
				return nil
			},
		},
		{
			Name: "任务C(耗时50ms)",
			Fn: func() error {
				time.Sleep(50 * time.Millisecond)
				fmt.Println("   -> 任务C 执行完成")
				return nil
			},
		},
		{
			Name: "任务D(模拟错误)",
			Fn: func() error {
				time.Sleep(30 * time.Millisecond)
				fmt.Println("   -> 任务D 执行完成(返回错误)")
				return fmt.Errorf("任务D 出现错误")
			},
		},
	}

	scheduler := NewScheduler(tasks...)

	startTotal := time.Now()
	results := scheduler.Run()
	total := time.Since(startTotal)

	fmt.Println("\n---------- 任务执行统计 ----------")
	for _, r := range results {
		status := "成功"
		if r.Err != nil {
			status = fmt.Sprintf("失败: %v", r.Err)
		}
		fmt.Printf("任务: %-18s | 耗时: %-10v | 状态: %s\n",
			r.Name, r.Duration, status)
	}
	fmt.Printf("\n所有任务并发执行总耗时: %v\n", total)
	fmt.Println("（注：若串行执行，总耗时约为各任务耗时之和）")
}

// ============================================================
// main 入口
// ============================================================
func main() {
	runOddEvenDemo()
	runSchedulerDemo()
}
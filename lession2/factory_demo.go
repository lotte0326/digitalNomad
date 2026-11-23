package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// ============================================================================
// 场景设定：超级玩具加工厂
// 包含：订单调度中心（生产者）、3个车间工人（消费者）、质检员（WaitGroup协调）
// ============================================================================

// Task 定义订单任务
type Task struct {
	ID       int
	ToyName  string
	WorkTime time.Duration
}

// 全局常量配置
const (
	MaxTasks = 10 // 模拟接单员总共接 10 个单子就下班
)

// ★ 知识点 4 & 最佳实践：sync.Once
// 面试常考：如何保证某个操作在多并发下绝对只执行一次？（比如初始化全局配置、关闭通道）
var closeOnce sync.Once

// ============================================================================
// 知识点 1：Goroutine 创建和使用
// 知识点 6：生产者-消费者模式（生产者端）
// ============================================================================
// dispatcher 接单员：负责生产任务并送入管道
// ★ 注意参数方向：chan<- Task 表示这是一个【只写】通道，编译器会帮我们检查，防止在这个函数里误读数据
func dispatcher(ctx context.Context, taskChan chan<- Task, wg *sync.WaitGroup) {
	// ★ 知识点 defer 的核心应用：确保无论发生什么（panic或提前返回），都一定会执行
	defer wg.Done()
	
	// ★ 面试重点：谁写数据，谁负责关闭通道。消费者不能关通道，否则生产者写入会 panic。
	// 我们用 sync.Once 配合 defer 确保通道一定被关闭，且只关闭一次。
	defer closeOnce.Do(func() {
		close(taskChan)
		fmt.Println("  [调度中心] 任务通道已安全关闭，不再接收新订单。")
	})

	toyList := []string{"变形金刚", "乐高积木", "遥控赛车", "芭比娃娃"}

	for i := 1; i <= MaxTasks; i++ {
		// ★ 知识点 3：Select 语句（快速响应取消）
		// 在生成任务前，先非阻塞地检查一下老板有没有发停工信号
		select {
		case <-ctx.Done():
			fmt.Println("  [调度中心] 收到停工信号，停止接单。")
			return // 触发上面的 defer，关闭通道
		default:
			// 没有停工信号，正常生产任务
			time.Sleep(100 * time.Millisecond) // 模拟生成订单耗时
			
			task := Task{
				ID:       i,
				ToyName:  toyList[rand.Intn(len(toyList))],
				WorkTime: time.Duration(rand.Intn(300)+200) * time.Millisecond,
			}

			// ★ 知识点 2：Channel 通信机制（带缓冲）
			// 如果 taskChan 满了（容量为 3），这里会阻塞（等待工人取走任务），这叫“背压机制”。
			taskChan <- task
			fmt.Printf("  [调度中心] 派发新任务: #%d (%s)\n", task.ID, task.ToyName)
		}
	}
	fmt.Println("  [调度中心] 今日任务额度已满，正常下班。")
}

// ============================================================================
// 知识点 1：Goroutine 创建和使用
// 知识点 6：生产者-消费者模式（消费者端）
// ============================================================================
// worker 工人：负责从管道读取任务并加工
// ★ 参数方向：<-chan Task 表示这是一个【只读】通道
func worker(ctx context.Context, workerID int, taskChan <-chan Task, wg *sync.WaitGroup) {
	defer wg.Done() // 确保工人结束时，计数器 -1

	fmt.Printf("  [工人 %d] 准备就绪，等待任务...\n", workerID)

	for {
		// ★ 知识点 3：Select 语句（面试高频用法）
		// 工人处于待命状态，同时监听两个事件：1. 管道来任务了 2. 老板喊下班了
		select {
		case <-ctx.Done():
			// 如果老板紧急叫停，立刻放下工作走人
			fmt.Printf("  [工人 %d] 收到紧急停工指令，赶紧溜了。\n", workerID)
			return

		case task, isOpen := <-taskChan:
			// ★ 知识点 2：Channel 多返回值判断
			// isOpen 为 false 意味着通道已被关闭，且里面的数据已经被拿光了
			if !isOpen {
				fmt.Printf("  [工人 %d] 传送带没任务了，正常下班。\n", workerID)
				return
			}

			// 模拟加工
			fmt.Printf("    -> [工人 %d] 开始加工 #%d (%s)...\n", workerID, task.ID, task.ToyName)
			
			// 模拟耗时，并且在这个耗时过程中，也要能响应取消信号！
			// ★ 面试极容易错的坑：如果直接 time.Sleep，老板取消时，工人得硬等 Sleep 结束。
			timer := time.NewTimer(task.WorkTime)
			select {
			case <-ctx.Done():
				timer.Stop()
				fmt.Printf("    -> [工人 %d] 加工 #%d 时被紧急叫停！\n", workerID, task.ID)
				return
			case <-timer.C:
				// 正常加工完毕
				fmt.Printf("    <- [工人 %d] 完成订单 #%d，总计耗时 %v\n", workerID, task.ID, task.WorkTime)
			}
		}
	}
}

// ============================================================================
// 主函数：老板视角（掌控全局）
// ============================================================================
func main() {
	// 设置随机数种子
	rand.Seed(time.Now().UnixNano())

	fmt.Println("========== 超级玩具加工厂正式运转 ==========")

	// ★ 知识点 5：Context 上下文控制
	// 面试常考：Context 应该作为函数的第一个参数，且命名为 ctx。
	// context.WithTimeout 设置一个绝对超时时间（比如 1.5 秒）。
	// 为什么用 WithTimeout？为了演示即使任务没派完，时间一到也会被强制中断。
	ctx, cancel := context.WithTimeout(context.Background(), 1500*time.Millisecond)
	
	// ★ 最佳实践：必须 defer cancel()。即使超时了，这里也能释放 Context 内部资源，防止内存泄漏。
	// 面试官极爱问：如果不写 cancel 会怎样？答：会导致 Context 绑定的计时器和子协程无法释放（上下文泄漏）。
	defer cancel()

	// ★ 知识点 2：Channel 通信机制
	// 创建一个容量为 3 的有缓冲通道。
	// 面试常考：无缓冲通道和有缓冲通道的区别？
	// 答：无缓冲是同步通信（必须一手交钱一手交货）；有缓冲是异步通信（可以提前放几个在传送带上）。
	taskChan := make(chan Task, 3)

	// ★ 知识点 4：WaitGroup 同步
	// 类似一个倒计时计数器。
	var wg sync.WaitGroup

	// 1. 启动调度中心（生产者）
	wg.Add(1)
	// ★ 知识点 1：go 关键字启动协程
	go dispatcher(ctx, taskChan, &wg)

	// 2. 启动 3 个工人（消费者）
	for i := 1; i <= 3; i++ {
		wg.Add(1) // 启动前先 +1
		go worker(ctx, i, taskChan, &wg)
	}

	fmt.Println("老板：工厂已开工，我在这里盯着监控（Main 被阻塞）...")

	// ★ 知识点 4：WaitGroup 阻塞等待
	// wg.Wait() 会阻塞主协程，直到计数器归零。
	// 面试题：如果调度中心不小心死锁了，永远不调用 Done，会发生什么？答：程序会报 "all goroutines are asleep - deadlock!" 并崩溃。
	wg.Wait()

	// 3. 补充：另一种监听方式
	// 如果不用 wg.Wait()，我们也可以监听 ctx.Done()，但 wg.Wait() 能确保所有子协程把清理逻辑跑完。
	if ctx.Err() == context.DeadlineExceeded {
		fmt.Println("老板：到时间了，拉闸断电！(Context 超时触发)")
	} else {
		fmt.Println("老板：今天任务提前完成，顺利收工！")
	}

	fmt.Println("========== 工厂大门关闭，所有协程安全退出 ==========")
}
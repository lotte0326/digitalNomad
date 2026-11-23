package main

import "fmt"

// 演示 defer 和 recover
func safeDivide(a, b int) (result int) {
    // 1. defer 必须写在 panic 可能发生之前
    defer func() {
        // 2. recover 必须在 defer 的匿名函数中调用
        if r := recover(); r != nil {
            fmt.Println("捕获到异常:", r)
            result = 0 // 恢复时给返回值赋个默认值
        }
    }()

    if b == 0 {
        panic("除数不能为 0!") // 主动触发 panic
    }
    return a / b
}

func main() {
    // --- if / else ---
    if num := 10; num > 5 { // 带初始化语句
        fmt.Println("num 大于 5")
    } else {
        fmt.Println("num 小于等于 5")
    }

    // --- switch ---
    day := "Mon"
    switch day {
    case "Mon", "Tue":
        fmt.Println("工作日开始")
    case "Sat", "Sun":
        fmt.Println("周末")
    default:
        fmt.Println("其他日子")
    }

    // --- for 循环 ---
    for i := 0; i < 3; i++ {
        fmt.Print(i, " ")
    }
    fmt.Println()

    // --- range 遍历 ---
    fruits := []string{"apple", "banana"}
    for index, fruit := range fruits {
        fmt.Printf("索引: %d, 水果: %s\n", index, fruit)
    }

    // --- defer 后进先出 ---
    defer fmt.Println("这行最后执行 (1)")
    defer fmt.Println("这行先执行 (2)")

    // --- panic / recover ---
    fmt.Println("计算结果:", safeDivide(10, 2))
    fmt.Println("计算结果:", safeDivide(10, 0)) // 这里触发 panic，被 recover 捕获

    fmt.Println("程序正常结束")
}
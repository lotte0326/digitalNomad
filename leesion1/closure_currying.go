package main

import "fmt"

// 闭包：计数器
// counter 返回一个匿名函数，该匿名函数“记住”了 count 变量
func counter() func() int {
	count := 0
	return func() int {
		count++ // 修改外部作用域的变量
		return count
	}
}

// 闭包：生成带有特定前缀的日志打印器
func makeLogger(prefix string) func(string) {
	return func(msg string) {
		fmt.Printf("[%s] %s\n", prefix, msg)
	}
}

// 柯里化：先传第一个数，返回一个函数，再传第二个数
func curriedAdd(a int) func(int) int {
	return func(b int) int {
		return a + b
	}
}

func main() {
	// 闭包计数器：每次调用 c1，它都会基于上一次的值累加
	c1 := counter()
	fmt.Println("计数器 c1:", c1()) // 1
	fmt.Println("计数器 c1:", c1()) // 2

	// 重点：c2 是一个全新的闭包，拥有独立的内存空间
	c2 := counter()
	fmt.Println("计数器 c2:", c2()) // 1

	// 闭包生成日志器
	errLog := makeLogger("ERROR")
	infoLog := makeLogger("INFO")
	errLog("数据库连接失败")
	infoLog("用户登录成功")

	// 柯里化
	add5 := curriedAdd(5) // 相当于生成了一个“加5”的函数
	fmt.Println("5 + 10 =", add5(10))
}
package main

import "fmt"

// 基础函数：有参数，有返回值
// 语法糖：a, b int 等价于 a int, b int
func add(a, b int) int {
	return a + b
}

// 无返回值函数
func sayHello(name string) {
	fmt.Println("Hello,", name)
}

// 函数作为参数传递（高阶函数）
func calculate(a, b int, fn func(int, int) int) int {
	return fn(a, b) 
}

func main() {
	sum := add(3, 5)
	fmt.Println("两数之和:", sum)

	sayHello("Gopher")

	// 把 add 函数当成参数传进去
	result := calculate(10, 20, add)
	fmt.Println("高阶函数计算结果:", result)
}
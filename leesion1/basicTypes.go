package main

import (
	"fmt"
	"strconv"
)

func main() {
	// 1. bool：只有 true 和 false
	var isGoFun bool = true
	fmt.Println("bool 类型:", isGoFun)

	// 2. int / uint：有符号与无符号
	var age int = 25       // 常用有符号整数
	var score uint = 100   // 无符号整数，不能为负
	fmt.Printf("int: %d, uint: %d\n", age, score)

	// 3. float：浮点数（小数）
	var price float64 = 9.99 // Go 中默认的小数就是 float64
	fmt.Println("float64:", price)

	// 4. string：双引号包裹，不可变
	var name string = "Gopher"
	// 字符串拼接
	greeting := "Hello, " + name
	fmt.Println("string:", greeting)

	// 补充：字符串与数字转换（必须用 strconv 包）
	numStr := "123"
	num, _ := strconv.Atoi(numStr) // 字符串转 int
	fmt.Printf("字符串 %s 转成整数后 + 1 = %d\n", numStr, num+1)
}
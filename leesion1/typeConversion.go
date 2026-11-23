package main

import (
	"fmt"
	"strconv"
)

func main() {
	// 1. 基本类型转换（必须显式转换，没有隐式转换）
	var a int = 42
	var b float64 = float64(a) // int 转 float64
	fmt.Printf("a 的类型: %T, 值: %v\n", a, a)
	fmt.Printf("b 的类型: %T, 值: %v\n", b, b)

	// 2. 精度丢失与截断
	var pi float64 = 3.14159
	var piInt int = int(pi)        // 浮点数转整数，直接丢弃小数部分（不四舍五入）
	fmt.Println("浮点数转为整数:", piInt) // 输出 3

	// 3. 字符串与数字转换（必须使用 strconv 包）
	// 整数转字符串
	str1 := strconv.Itoa(123) // Itoa = Integer to ASCII
	fmt.Printf("整数转字符串: 类型 %T, 值 %s\n", str1, str1)

	// 字符串转整数
	str2 := "456"
	num2, err := strconv.Atoi(str2) // Atoi = ASCII to Integer
	if err == nil {                 // 检查是否有错误（例如尝试转换 "abc" 会报错）
		fmt.Printf("字符串转整数: 类型 %T, 值 %d\n", num2, num2)
	}
}

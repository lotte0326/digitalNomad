package main

import (
	"errors"
	"fmt"
)

// 多返回值：返回商和余数
func divide(a, b int) (int, int) {
	return a / b, a % b
}

// 多返回值配合 error：处理可能的错误
func safeSqrt(n float64) (float64, error) {
	if n < 0 {
		return 0, errors.New("不能对负数求平方根")
	}
	return n, nil 
}

// 可变参数：求任意数量整数的和
// nums 在函数内部是一个 []int 切片
func sumAll(nums ...int) int {
	total := 0
	for _, num := range nums {
		total += num
	}
	return total
}

func main() {
	// 1. 接收多返回值
	q, r := divide(10, 3)
	fmt.Printf("10 / 3 = 商 %d, 余数 %d\n", q, r)

	// 2. 接收结果与错误
	res, err := safeSqrt(-4)
	if err != nil {
		fmt.Println("发生错误:", err)
	} else {
		fmt.Println("计算结果:", res)
	}

	// 3. 调用可变参数函数
	fmt.Println("不传参:", sumAll())
	fmt.Println("传3个参数:", sumAll(1, 2, 3))
	
	// 把切片打散传入可变参数（注意后面的 ...）
	nums := []int{10, 20, 30}
	fmt.Println("传切片:", sumAll(nums...))
}
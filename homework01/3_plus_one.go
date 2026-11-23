// 文件名：plus_one.go
package main

import "fmt"

// 加一函数（原地修改）
func plusOne(digits []int) []int {
	// 从数组末尾（最低位）开始向前遍历
	for i := len(digits) - 1; i >= 0; i-- {
		// 如果当前位不是 9，那么加 1 后不会产生进位，直接加 1 即可
		if digits[i] != 9 {
			digits[i]++
			return digits // 加完直接返回，后续的位不受影响
		}
		// 如果当前位是 9，加 1 变成 10，当前位写 0，进位留给下一次循环
		digits[i] = 0
	}

	// 如果循环完整结束，说明原数组全是 9（比如 [9,9]）
	// 加 1 后需要多进一位，结果是 1 后面跟着 len(digits) 个 0
	result := make([]int, len(digits)+1)
	result[0] = 1 // 新数组首位设为 1，后面默认为 0
	return result
}

func main() {
	// 测试用例
	testCases := [][]int{
		{1, 2, 3},    // 预期：[1, 2, 4]
		{4, 3, 2, 1}, // 预期：[4, 3, 2, 2]
		{9},          // 预期：[1, 0]
		{9, 9, 9},    // 预期：[1, 0, 0, 0]（进位测试）
	}

	for _, digits := range testCases {
		// 【关键修复】：由于 plusOne 会修改原切片，我们需要先拷贝一份
		// 拷贝一份输入数据作为真正的“原始输入”留作打印
		inputCopy := make([]int, len(digits))
		copy(inputCopy, digits)

		// 把原切片 digits 传给函数
		result := plusOne(digits)

		// 打印时：输入用备份的 inputCopy，输出用修改后的 result
		fmt.Printf("输入: %v -> 输出: %v\n", inputCopy, result)
	}
}
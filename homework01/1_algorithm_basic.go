// 文件名：algorithm_basic.go
package main

import (
	"fmt"
)

func main() {
	fmt.Println("========== 第一题：只出现一次的数字 ==========")
	nums1 := []int{4, 1, 2, 1, 2}
	// 预期结果：4
	fmt.Printf("数组 %v 中只出现一次的数字是: %d\n", nums1, singleNumber(nums1))

	fmt.Println("\n========== 第二题：回文数 ==========")
	testNums := []int{121, -121, 10, 12321}
	for _, num := range testNums {
		fmt.Printf("数字 %d 是回文数吗？ %v\n", num, isPalindrome(num))
	}
}

// ============================================================================
// 第一题：只出现一次的数字
// ============================================================================

// 解法一：Map 哈希表法（小白最易理解，符合题目要求）
// 时间复杂度：O(N)，空间复杂度：O(N)
func singleNumber(nums []int) int {
	// 1. 创建一个 map，键是数组里的数字，值是它出现的次数
	// 语法：make(map[键的类型]值的类型)
	countMap := make(map[int]int)

	// 2. 遍历数组，记录每个数字出现的次数
	for _, num := range nums {
		// 如果这个数字在 map 里不存在，Go 默认返回 0，所以直接 +1 没问题
		// 如果存在，就取出原来的次数 +1
		countMap[num]++
	}

	// 3. 遍历 map，找出那个次数为 1 的数字
	// map 遍历的语法：for 键, 值 := range map名
	for key, count := range countMap {
		if count == 1 {
			return key // 找到了，直接返回
		}
	}

	return -1 // 理论上不会走到这里（题目保证一定有解）
}

// 【面试进阶必杀技】解法二：异或运算（XOR）—— 面试官最想看到的答案
// 原理：任何数和自己异或都是 0（a ^ a = 0），任何数和 0 异或都是它自己（a ^ 0 = a）。
// 异或运算满足交换律和结合律，所以把数组中所有数字异或一遍，成对的抵消为0，剩下的就是单身狗。
// 时间复杂度：O(N)，空间复杂度：O(1) —— 完美！
/*
func singleNumberAdvanced(nums []int) int {
	result := 0
	for _, num := range nums {
		result ^= num // ^= 是异或赋值运算符
	}
	return result
}
*/

// ============================================================================
// 第二题：回文数
// ============================================================================

// 解法一：数字反转法（不用转字符串，纯数学计算，面试加分项）
// 时间复杂度：O(log10(N))，空间复杂度：O(1)
func isPalindrome(x int) bool {
	// 【防坑点 1】：负数一定不是回文数（因为负号无法对称）
	// 【防坑点 2】：如果数字的最后一位是 0，但数字本身不是 0（比如 10），那它肯定不是回文数。
	if x < 0 || (x%10 == 0 && x != 0) {
		return false
	}

	// 假设 x = 121，我们要把它反转一半，对比前一半和后一半是否相等。
	// 这样不仅省事，还能防止数字太大导致 int 溢出。
	revertedNumber := 0

	// 当原数字大于反转后的数字时，说明还没走到一半，继续反转
	for x > revertedNumber {
		// 获取 x 的最后一位数字（例如 121 % 10 = 1）
		revertedNumber = revertedNumber*10 + x%10
		// x 砍掉最后一位（例如 121 / 10 = 12）
		x /= 10
	}

	// 当循环结束时，有两种情况：
	// 1. 数字长度是偶数（如 1221）：x 变成了 12，revertedNumber 也变成了 12。直接相等。
	// 2. 数字长度是奇数（如 121）：x 变成了 1，revertedNumber 变成了 12。我们需要用 revertedNumber / 10 去掉中间那个数字（12/10=1），再和 x 比较。
	return x == revertedNumber || x == revertedNumber/10
}

// 解法二：字符串反转法（小白最容易想到，借助内置函数，但有性能开销）
/*
func isPalindromeString(x int) bool {
	// 负数直接返回 false
	if x < 0 {
		return false
	}
	// 把数字转成字符串
	str := strconv.Itoa(x)
	// 双指针技巧：一个从头往后走，一个从尾往前走，对比字符是否相等
	left, right := 0, len(str)-1
	for left < right {
		if str[left] != str[right] {
			return false
		}
		left++
		right--
	}
	return true
}
*/
// 文件名：two_sum_demo.go
package main

import "fmt"

func main() {
	fmt.Println("========== 第一题：两数之和 ==========")

	// 测试用例 1
	nums1 := []int{2, 7, 11, 15}
	target1 := 9
	fmt.Printf("输入: nums = %v, target = %d\n", nums1, target1)
	fmt.Printf("输出: %v\n\n", twoSum(nums1, target1))

	// 测试用例 2
	nums2 := []int{3, 2, 4}
	target2 := 6
	fmt.Printf("输入: nums = %v, target = %d\n", nums2, target2)
	fmt.Printf("输出: %v\n\n", twoSum(nums2, target2))

	// 测试用例 3（有重复元素）
	nums3 := []int{3, 3}
	target3 := 6
	fmt.Printf("输入: nums = %v, target = %d\n", nums3, target3)
	fmt.Printf("输出: %v\n", twoSum(nums3, target3))
}

// ============================================================================
// 🧩 最优解法：Map 哈希表法（时间复杂度 O(N)，空间复杂度 O(N)）
// ============================================================================
// 【核心思路】：把“找两个数”变成“找一个人”。
// 我们在遍历数组时，问自己一句：“我现在的数字是 x，那么能和我凑成 target 的 y 在哪里？”
// 答案显然是 y = target - x。
// 我们只需要拿个本子（Map）记录下来：我见过哪些数字，它们分别在哪个位置。
// 只要本子里有 y，就直接大功告成！
// 
// 【为什么用 Map？】
// 因为 Map 的查找速度是 O(1)（瞬间找到）。数组查找是 O(N)（要挨个看）。
func twoSum(nums []int, target int) []int {
	// 1. 准备一个“记忆本”（Map）
	// Key 存“数字本身”，Value 存“这个数字在数组中的下标”
	// 语法：make(map[键的类型]值的类型)
	seen := make(map[int]int)

	// 2. 遍历数组，一边看，一边记
	for i, x := range nums {
		// 3. 计算理想中的“另一半” y 应该是多少
		y := target - x

		// 4. 去记忆本里查一下，以前见过 y 吗？
		// 【高频考点】：这里用到了 Go 独有的 "comma-ok" 语法
		// index 是 y 的下标，ok 是一个布尔值（true/false）
		if index, ok := seen[y]; ok {
			// 如果见过，说明找到了！直接返回 y 的下标和当前 x 的下标
			return []int{index, i}
		}

		// 5. 如果没见过 y，就把当前的 x 和它的下标记下来，留给后面的人用
		seen[x] = i
	}

	// 题目保证有且仅有一个答案，所以理论上走不到这里
	return nil
}

// ============================================================================
// ❌ 小白避坑：暴力解法（时间复杂度 O(N^2)），仅供理解，面试不推荐
// ============================================================================
// 如果用两层 for 循环，就像全班同学两两握手。效率极低，LeetCode 提交会超时。
/*
func twoSumBruteForce(nums []int, target int) []int {
	for i := 0; i < len(nums); i++ {
		for j := i + 1; j < len(nums); j++ {
			if nums[i]+nums[j] == target {
				return []int{i, j}
			}
		}
	}
	return nil
}
*/
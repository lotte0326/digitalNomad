// 文件名：slice_algorithm_demo.go
package main

import (
	"fmt"
	"sort"
)

func main() {
	fmt.Println("==================================================")
	fmt.Println("        🚀 切片算法大练兵 🚀        ")
	fmt.Println("==================================================")

	// 运行第一题：删除有序数组中的重复项
	fmt.Println("\n【第一题：删除有序数组中的重复项】")
	// 【防坑注意】因为题目要求原地修改，会破坏原切片，所以测试前先准备好数据
	nums1 := []int{1, 1, 2}
	fmt.Printf("原始数组: %v\n", nums1)
	newLen1 := removeDuplicates(nums1)
	// 注意：只打印新长度范围内的有效数据
	fmt.Printf("去重后长度: %d, 数组前 %d 个元素为: %v\n", newLen1, newLen1, nums1[:newLen1])

	nums2 := []int{0, 0, 1, 1, 1, 2, 2, 3, 3, 4}
	fmt.Printf("\n原始数组: %v\n", nums2)
	newLen2 := removeDuplicates(nums2)
	fmt.Printf("去重后长度: %d, 数组前 %d 个元素为: %v\n", newLen2, newLen2, nums2[:newLen2])

	// 运行第二题：合并区间
	fmt.Println("\n【第二题：合并区间】")
	intervals := [][]int{{1, 3}, {2, 6}, {8, 10}, {15, 18}}
	fmt.Printf("原始区间: %v\n", intervals)
	merged := merge(intervals)
	fmt.Printf("合并后区间: %v\n", merged)
}

// ============================================================================
// 🧩 第一题：删除有序数组中的重复项 (Remove Duplicates from Sorted Array)
// ============================================================================
// 思路：双指针法（快慢指针）。
// 因为数组是有序的，重复的元素一定挨在一起。
// 慢指针 i 指向“已经去重后的最后一个元素”的位置。
// 快指针 j 负责往前探路，找新的不重复元素。

// 时间复杂度：O(N)，空间复杂度：O(1)（符合题目要求）
func removeDuplicates(nums []int) int {
	// 【防坑 1】：如果数组为空，直接返回 0
	if len(nums) == 0 {
		return 0
	}

	// 慢指针 i，初始指向第一个元素（它肯定不重复）
	i := 0

	// 快指针 j，从第二个元素开始往前遍历
	for j := 1; j < len(nums); j++ {
		// 如果快指针找到的元素，和慢指针指向的元素不一样，说明找到了新的不重复元素
		if nums[j] != nums[i] {
			// 慢指针往前移动一步，腾出位置
			i++
			// 把新找到的元素，放到慢指针的位置上（原地覆盖旧的重复元素）
			nums[i] = nums[j]
		}
		// 如果一样，快指针继续往前找，慢指针不动
	}

	// i 是最后一个不重复元素的索引，所以新长度是 i + 1
	return i + 1
}

// ============================================================================
// 🧩 第二题：合并区间 (Merge Intervals)
// ============================================================================
// 思路：先排序，再合并。
// 1. 按照区间的起始位置（start）从小到大排序。这样重叠的区间一定会在相邻的位置。
// 2. 准备一个结果切片 merged。
// 3. 遍历排序后的区间：
//    - 如果 merged 为空，或者当前区间的 start 大于 merged 最后一个区间的 end，说明没有重叠，直接加入结果。
//    - 否则说明有重叠，需要合并：把 merged 最后一个区间的 end 更新为“两者 end 的较大值”。

func merge(intervals [][]int) [][]int {
	// 【防坑 1】：如果区间数组为空，直接返回空
	if len(intervals) == 0 {
		return [][]int{}
	}

	// 1. 对区间按照起始位置（第一个元素）进行升序排序
	// 注意：这里用到了 Go 的 sort.Slice 函数和匿名函数
	sort.Slice(intervals, func(i, j int) bool {
		return intervals[i][0] < intervals[j][0]
	})

	// 2. 准备结果切片，先把第一个区间放进去作为基准
	merged := [][]int{intervals[0]}

	// 3. 从第二个区间开始遍历
	for i := 1; i < len(intervals); i++ {
		current := intervals[i] // 当前遍历到的区间
		// 获取结果切片中最后一个区间
		last := merged[len(merged)-1]

		// 【核心判断】：当前区间的起始位置，是否小于或等于最后一个区间的结束位置？
		// 如果是（比如 [1,3] 和 [2,6]），说明发生了重叠
		if current[0] <= last[1] {
			// 合并操作：取两个区间结束位置的较大值，更新最后一个区间的结束位置
			if current[1] > last[1] {
				last[1] = current[1]
			}
			// 注意：last 是切片的元素引用，修改 last[1] 就等于修改了 merged 最后一个元素的第 2 个值
		} else {
			// 没有重叠（比如 [1,3] 和 [8,10]），直接把当前区间追加到结果中
			merged = append(merged, current)
		}
	}

	return merged
}
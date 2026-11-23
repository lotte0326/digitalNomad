package main

import "fmt"

func main() {
	// --- 1. array（数组）---
	// 定长：[3]int 表示只能装 3 个 int。长度一旦定义不能变。
	var arr [3]int = [3]int{1, 2, 3}
	fmt.Println("数组:", arr, "长度为:", len(arr))

	// --- 2. slice（切片）---
	// 变长：最常用。可以通过 make 创建，也可以用字面量
	var s1 []int = []int{10, 20, 30} // 创建并初始化
	fmt.Println("切片:", s1, "长度:", len(s1), "容量:", cap(s1))

	// 使用 append 动态追加元素（容量不够时自动扩容）
	s1 = append(s1, 40)
	fmt.Println("追加元素后的切片:", s1)

	// 切片操作（左闭右开：包含起始索引，不包含结束索引）
	s2 := s1[1:3] // 取索引 1 和 2 的元素
	fmt.Println("切片截取 s1[1:3]:", s2)

	// --- 3. map（映射/字典）---
	// 必须用 make 初始化才能赋值
	// key 是 string 类型，value 是 int 类型
	scores := make(map[string]int)
	scores["Alice"] = 90
	scores["Bob"] = 85

	// 获取值时，可以返回两个参数：value 和 是否存在（ok）
	aliceScore, exists := scores["Alice"]
	fmt.Println("Alice 的分数:", aliceScore, "是否存在:", exists)

	// 删除某个 key
	delete(scores, "Bob")
	fmt.Println("map 内容:", scores)
}
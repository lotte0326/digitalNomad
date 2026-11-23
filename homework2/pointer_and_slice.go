// 文件名：pointer_and_slice.go
//
// 知识点：指针（pointer）与切片（slice）
// 对应 Go 进阶课程 · Lesson 03 作业
//
//	题目 1：定义一个函数，接收一个整数指针，在函数内部将该指针指向的值增加 10，
//	        然后在 main 中调用并输出修改后的值。
//	        考察点：指针的使用、值传递与引用传递的区别。
//
//	题目 2：实现一个函数，接收一个整数切片的指针，将切片中的每个元素乘以 2。
//	        考察点：指针、切片操作。
//
// 运行：
//
//	cd /Users/tangming/Documents/DigitalNomad/code/digitalnomad
//	go run ./homework2
//
// 核心结论（背下来）：
//  1. Go 只有值传递。传指针 = 传「地址的拷贝」，但顺着地址改的是原来那块内存，所以外面看得见。
//  2. 切片是描述符 {ptr, len, cap}，按值拷贝；拷贝和原件指向同一个底层数组。
//  3. 所以：改元素 s[i] = x → 外面看得见（不传指针也行）；
//     改切片头 s = append(...) / s = s[1:] → 外面看不见（要传 *[]int 或用返回值）。
//  4. nil 指针解引用会 panic，函数入口先判 if p == nil { return }。
package main

import "fmt"

// ============================================================
// 题目 1：整数指针 +10
// ============================================================

// addTen 接收 *int（指向 int 的指针），把指针指向的值加 10。
//
// 注意：Go 里只有「值传递」，形参 a 拿到的只是「地址」这个整数值的一份拷贝。
// 但 *a = *a + 10 会顺着这个地址去改「原变量本身」，所以 main 能看见变化。
func addTen(a *int) {
	if a == nil { // 防御：nil 指针解引用会 panic
		return
	}
	*a += 10 // *a 读值，*a = ... 写回原变量所在的内存；等价于 *a = *a + 10
}

// addTenWrong 是「值传递」的反例：拿到的是一份 int 拷贝，
// 函数内加的 10 只改这份拷贝，调用方完全不受影响。
func addTenWrong(a int) {
	a += 10
}

// changePointer 演示：指针变量本身也是按值传递的。
// 函数里 p = &x 只改了形参这份拷贝，调用方的 p 仍然指向原来的地址。
func changePointer(p *int) {
	x := 100
	p = &x
	fmt.Printf("    [changePointer 内] 形参 p 现在指向 x = %d\n", *p)
}

// ============================================================
// 题目 2：整数切片指针，每个元素 ×2
// ============================================================

// doubleSlice 接收「切片指针」*[]int，把切片里每个元素乘以 2。
//
// 括号陷阱：*s[i] 会被解析成 *(s[i])，对指针做下标，编译报错；
// 必须写 (*s)[i]。更简洁的写法是 for i := range *s { (*s)[i] *= 2 }。
func doubleSlice(s *[]int) {
	if s == nil { // 连切片指针本身都是 nil
		return
	}
	if *s == nil { // 指针指向一个 nil 切片
		return
	}
	for i := range *s {
		(*s)[i] *= 2
	}
}

// doubleSliceNew 演示 *[]int 的真正威力：构造新切片再整块替换回去。
// 长度变化 / 底层数组重新分配也能反映到调用方 —— 这是只用 []int 做不到的。
func doubleSliceNew(s *[]int) {
	if s == nil || *s == nil {
		return
	}
	ns := make([]int, 0, len(*s)) // 预分配容量，避免多次扩容
	for _, v := range *s {
		ns = append(ns, v*2)
	}
	*s = ns
}

// doubleInPlace 对比写法：只传 []int，用下标原地修改。
// 能改元素（因为在改共享的底层数组），但不能改长度。
func doubleInPlace(s []int) {
	for i := range s {
		s[i] *= 2
	}
}

// appendOne 对比写法：只传 []int 时，append 的结果只活在本函数内。
func appendOne(s []int) {
	s = append(s, 999) // 只改了这份拷贝的「切片头」
	fmt.Printf("    [appendOne 内] 形参 = %v\n", s)
}

func main() {
	fmt.Println("========== 题目 1：整数指针 +10 ==========")

	num := 5
	fmt.Printf("调用前                   num = %d\n", num)

	addTen(&num) // 传 &num，类型就是 *int
	fmt.Printf("调用 addTen(&num) 之后   num = %d\n", num)

	addTenWrong(num) // 反例：传值
	fmt.Printf("调用 addTenWrong 之后    num = %d（没变，因为传的是拷贝）\n", num)

	p := &num
	changePointer(p)
	fmt.Printf("changePointer 之后       p 仍指向 num，num = %d\n", num)

	var nilP *int
	addTen(nilP) // nil 保护生效，不会 panic
	fmt.Println("addTen(nil) 安全返回，没有 panic")

	fmt.Println()
	fmt.Println("========== 题目 2：切片指针，元素 ×2 ==========")

	nums := []int{1, 2, 3, 4, 5}
	fmt.Printf("调用前                       nums = %v\n", nums)
	doubleSlice(&nums)
	fmt.Printf("doubleSlice(&nums) 之后      nums = %v\n", nums)

	nums2 := []int{10, 20, 30}
	doubleSliceNew(&nums2)
	fmt.Printf("doubleSliceNew(&nums2) 之后  nums2 = %v（整块替换，长度可变）\n", nums2)

	fmt.Println()
	fmt.Println("---- 对比：只传 []int 会怎样 ----")

	a := []int{1, 2, 3}
	doubleInPlace(a)
	fmt.Printf("doubleInPlace(a) 之后        a = %v（改元素：看得见）\n", a)

	b := []int{1, 2, 3}
	appendOne(b)
	fmt.Printf("appendOne(b) 之后            b = %v（改长度：看不见！）\n", b)

	doubleSlice(nil) // 同样有 nil 保护
	fmt.Println("doubleSlice(nil) 安全返回，没有 panic")
}

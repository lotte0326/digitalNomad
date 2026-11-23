package main

import "fmt"

// 定义一个结构体用于演示
type Person struct {
	Name string
	Age  int
}

// 传递指针：可以修改原变量
func growUp(p *Person) {
	p.Age++ // Go 语言支持语法糖，p.Age 等价于 (*p).Age
}

// 传递值：只会修改副本，原变量不变
func growUpFail(p Person) {
	p.Age++
}

func main() {
	// 1. 基础用法：获取地址和修改值
	num := 100
	ptr := &num // & 获取变量的内存地址

	fmt.Printf("num 的值: %d, 内存地址: %p\n", num, ptr)
	
	*ptr = 200 // * 解引用，通过指针修改原变量的值
	fmt.Println("通过指针修改后 num 的值:", num)

	// 2. 实际应用：在函数中修改结构体
	person := Person{Name: "Tom", Age: 20}
	
	// 传指针（传递的是内存地址，非常轻量）
	growUp(&person)
	fmt.Println("传指针后，Tom 的年龄:", person.Age) // 输出 21

	// 传值（会完整拷贝一份结构体，无法修改原数据）
	growUpFail(person)
	fmt.Println("传值后，Tom 的年龄:", person.Age) // 依旧是 21
}
package main

import (
	"fmt"
	"math"
)

// ============================================================
// 题目一：Shape 接口的定义与实现
// ============================================================

// Shape 接口：包含 Area() 和 Perimeter() 两个方法
type Shape interface {
	Area() float64
	Perimeter() float64
}

// Rectangle 矩形结构体
type Rectangle struct {
	Width  float64
	Height float64
}

// Area 实现 Shape 接口的 Area 方法
func (r Rectangle) Area() float64 {
	return r.Width * r.Height
}

// Perimeter 实现 Shape 接口的 Perimeter 方法
func (r Rectangle) Perimeter() float64 {
	return 2 * (r.Width + r.Height)
}

// Circle 圆形结构体
type Circle struct {
	Radius float64
}

// Area 实现 Shape 接口的 Area 方法
func (c Circle) Area() float64 {
	return math.Pi * c.Radius * c.Radius
}

// Perimeter 实现 Shape 接口的 Perimeter 方法（周长 = 2πr）
func (c Circle) Perimeter() float64 {
	return 2 * math.Pi * c.Radius
}

// printShapeInfo 接收 Shape 接口，体现多态：任何实现 Shape 的类型都能传入
func printShapeInfo(s Shape) {
	fmt.Printf("  类型: %T\n", s)
	fmt.Printf("  面积: %.2f\n", s.Area())
	fmt.Printf("  周长: %.2f\n", s.Perimeter())
}

// runShapeDemo 题目一演示
func runShapeDemo() {
	fmt.Println("========== 题目一：Shape 接口 ==========")

	rect := Rectangle{Width: 4, Height: 3}
	circle := Circle{Radius: 2}

	fmt.Println("[矩形]")
	printShapeInfo(rect)

	fmt.Println("[圆形]")
	printShapeInfo(circle)

	// 直接调用方法（题目要求）
	fmt.Printf("矩形面积=%.2f, 矩形周长=%.2f\n", rect.Area(), rect.Perimeter())
	fmt.Printf("圆形面积=%.2f, 圆形周长=%.2f\n", circle.Area(), circle.Perimeter())

	// 使用接口切片演示多态
	shapes := []Shape{rect, circle}
	fmt.Println("\n[通过接口切片遍历]")
	for i, s := range shapes {
		fmt.Printf("  Shape[%d] 面积=%.2f 周长=%.2f\n", i, s.Area(), s.Perimeter())
	}
	fmt.Println()
}

// ============================================================
// 题目二：组合的使用与方法接收者
// ============================================================

// Person 人员结构体
type Person struct {
	Name string
	Age  int
}

// Employee 员工结构体：组合 Person（匿名嵌入）并添加 EmployeeID
type Employee struct {
	Person     // 匿名嵌入（组合）
	EmployeeID string
}

// PrintInfo 为 Employee 实现 PrintInfo 方法（值接收者）
func (e Employee) PrintInfo() {
	fmt.Printf("  姓名: %s\n", e.Name)       // 可直接访问嵌入字段 Name
	fmt.Printf("  年龄: %d\n", e.Age)        // 可直接访问嵌入字段 Age
	fmt.Printf("  工号: %s\n", e.EmployeeID)
}

// 演示指针接收者的方法（可选，展示两种接收者）
func (e *Employee) SetAge(age int) {
	e.Age = age
}

// runEmployeeDemo 题目二演示
func runEmployeeDemo() {
	fmt.Println("========== 题目二：组合与方法接收者 ==========")

	emp := Employee{
		Person: Person{
			Name: "张三",
			Age:  28,
		},
		EmployeeID: "E10086",
	}

	fmt.Println("[员工信息]")
	emp.PrintInfo()

	// 组合带来的"字段提升"：可像访问自己字段一样访问 Person 的字段
	fmt.Printf("\n字段提升演示: emp.Name=%s, emp.Person.Name=%s\n", emp.Name, emp.Person.Name)

	// 使用指针接收者方法修改
	emp.SetAge(29)
	fmt.Println("\n[修改年龄后]")
	emp.PrintInfo()

	// 组合的方法提升：Person 若有方法，也可被 Employee 直接调用（此处无额外方法，仅说明机制）
	fmt.Println()
}

// ============================================================
// main 入口
// ============================================================
func main() {
	runShapeDemo()
	runEmployeeDemo()
}
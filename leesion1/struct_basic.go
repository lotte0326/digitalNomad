package main

import "fmt"

// 标准定义
type Person struct {
	Name string
	Age  int
}

// 结构体嵌套
type Address struct {
	City string
}
type User struct {
	Username string
	Addr     Address 
}

// 匿名结构体切片
var tempData = []struct {
	ID   int
	Name string
}{
	{1, "Alice"},
	{2, "Bob"},
}

func main() {
	// 实例化
	p1 := Person{"Tom", 18}
	p2 := Person{Name: "Jerry", Age: 20}
	p3 := new(Person)
	p3.Name = "Spike"

	fmt.Println("p1:", p1)
	fmt.Println("p2:", p2)
	fmt.Println("p3 指针指向的内容:", *p3)

	// 访问嵌套结构体
	u := User{
		Username: "gopher_01",
		Addr:     Address{City: "Beijing"},
	}
	fmt.Println("用户所在城市:", u.Addr.City)

	// 使用匿名结构体
	for _, item := range tempData {
		fmt.Printf("ID: %d, Name: %s\n", item.ID, item.Name)
	}
}
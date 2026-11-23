package main

import (
	"errors"
	"fmt"
)

// ==========================================
// 知识点 1：接口定义和实现
// ==========================================
// 接口（Interface）是一组方法签名的集合。
// 在 Go 中，接口是隐式实现的：只要一个类型实现了接口里的所有方法，它就自动实现了该接口。
// 这里定义一个名为 Payment 的接口，包含两个方法。
type Payment interface {
	Pay(amount float64) (string, error) // 支付方法
	Name() string                       // 获取支付方式名称
}

// 定义一个微信支付结构体
type WeChatPay struct {
	UserID string
}

// 让 WeChatPay 实现 Payment 接口（注意方法名和参数必须完全一致）
func (w WeChatPay) Pay(amount float64) (string, error) {
	if amount <= 0 {
		return "", errors.New("金额必须大于0")
	}
	return fmt.Sprintf("微信用户 %s 成功支付了 %.2f 元", w.UserID, amount), nil
}

func (w WeChatPay) Name() string {
	return "微信支付"
}

// 定义一个支付宝结构体
type AliPay struct {
	Account string
}

// 让 AliPay 也实现 Payment 接口
func (a AliPay) Pay(amount float64) (string, error) {
	if amount > 50000 {
		return "", errors.New("单笔支付不能超过 5 万元")
	}
	return fmt.Sprintf("支付宝账号 %s 成功支付了 %.2f 元", a.Account, amount), nil
}

func (a AliPay) Name() string {
	return "支付宝"
}

// ==========================================
// 知识点 2：空接口和类型断言
// ==========================================
// 空接口 interface{}（在 Go 1.18 后可写为 any）不包含任何方法，
// 因此所有的类型都实现了空接口，它常用来接收任意类型的参数。

// 打印任意类型数据的函数
func PrintAny(v interface{}) {
	// 类型断言：判断 v 是否是 string 类型
	// 语法：val, ok := v.(目标类型)
	strVal, ok := v.(string)
	if ok {
		fmt.Println("这是字符串，且长度为:", len(strVal), "内容:", strVal)
	} else {
		fmt.Printf("这不是字符串，它的实际类型是: %T, 值是: %v\n", v, v)
	}

	// 使用 type switch 进行多种类型的判断（非常常用的模式）
	switch val := v.(type) {
	case int:
		fmt.Println("  分支判断：这是 int 类型，值为", val)
	case float64:
		fmt.Println("  分支判断：这是 float64 类型，值为", val)
	case Payment:
		fmt.Println("  分支判断：这是一个支付对象，名称为", val.Name())
	default:
		fmt.Printf("  分支判断：未知类型 %T\n", val)
	}
}

// ==========================================
// 知识点 3：接口组合
// ==========================================
// Go 允许把一个或多个接口嵌入到另一个接口中，形成新的接口。
type Reader interface {
	Read() string
}

type Writer interface {
	Write(s string)
}

// 组合接口：ReadWriter 包含了 Reader 和 Writer 的所有方法
type ReadWriter interface {
	Reader
	Writer
}

// 定义一个实现组合接口的结构体
type File struct {
	content string
}

func (f *File) Read() string {
	return f.content
}

func (f *File) Write(s string) {
	f.content = s
}

// ==========================================
// 知识点 4：多态应用
// ==========================================
// 多态意味着：相同的接口调用，因为具体实现不同，表现出不同的行为。
// 下面的函数接收 Payment 接口类型，不管传入的是微信还是支付宝，都能正常调用。
func Checkout(p Payment, amount float64) {
	fmt.Printf("\n[准备使用 %s 结算...]\n", p.Name())
	msg, err := p.Pay(amount)
	if err != nil {
		fmt.Println("支付失败:", err)
	} else {
		fmt.Println("支付成功:", msg)
	}
}

// ==========================================
// 知识点 5：接口最佳实践
// ==========================================
// 最佳实践 1：接口隔离原则（接口要尽量小）。
// 不要定义一个包含十几个方法的超级大接口，而是拆分成多个小接口（如上面的 Reader、Writer）。
// 最佳实践 2：“接受接口，返回结构体”。
// 编写函数时，参数尽量是接口（为了灵活，如 Checkout），返回值尽量是具体结构体（为了调用方便）。

// 下面是一个典型的小接口定义，只包含一个方法
type Stringer interface {
	String() string
}

// 最佳实践 3：编译期断言。
// 确保某个结构体一定实现了某个接口。如果没有实现，编译时会报错，防止上线后才发现问题。
var _ Payment = (*WeChatPay)(nil) // 检查 *WeChatPay 是否实现了 Payment

func main() {
	fmt.Println("=== 1 & 4. 接口定义与多态应用 ===")
	wx := WeChatPay{UserID: "wx_12345"}
	ali := AliPay{Account: "ali_98765"}

	// 多态体现：Checkout 接收的是 Payment 接口，但传入不同的结构体，执行结果不同
	Checkout(wx, 100.50)
	Checkout(ali, 200.00)
	Checkout(ali, 60000.00) // 触发支付宝的报错逻辑

	fmt.Println("\n=== 2. 空接口与类型断言 ===")
	PrintAny("Hello Go")
	PrintAny(42)
	PrintAny(3.14)
	PrintAny(wx) // 传入实现了 Payment 的结构体

	fmt.Println("\n=== 3. 接口组合 ===")
	var rw ReadWriter = &File{content: "初始内容"}
	rw.Write("写入新数据：Go 语言接口组合")
	fmt.Println("文件读取内容:", rw.Read())

	fmt.Println("\n=== 5. 接口最佳实践（Stringer） ===")
	// 注意：这里只是为了演示 Stringer 接口。
	// 在实际开发中，实现了 String() 方法的对象，用 fmt.Println 打印时会自动调用该方法的返回值。
	fmt.Println("最佳实践：接口尽量小，这里是编译期断言演示。")
}
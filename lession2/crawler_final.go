// 文件名：crawler_final.go
package main

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

// =====================================================================
// 🎬 场景设定：我们要写一个微型爬虫，抓取网页，计算哈希，存成 JSON 文件。
// =====================================================================

// 1. JSON 编解码：定义数据结构，准备存入文件
type WebPage struct {
	URL        string    `json:"url"`         // 网址
	Title      string    `json:"title"`       // 网页标题（模拟抓取）
	Content    string    `json:"content"`     // 网页内容
	HashMD5    string    `json:"hash_md5"`    // MD5 摘要
	HashSHA256 string    `json:"hash_sha256"` // SHA256 摘要
	FetchTime  string    `json:"fetch_time"`  // 抓取时间
}

func main() {
	fmt.Println("=== 简易网页爬虫启动 ===")

	// =====================================================================
	// 知识点 1：网络编程基础 & Web 爬虫实战
	// =====================================================================
	// 我们用 net/http 的 Get 方法去请求一个真实的网址
	targetURL := "https://httpbin.org/html" 
	fmt.Printf("正在抓取网页: %s ...\n", targetURL)

	// 创建一个带超时的 http 客户端（强烈建议：生产环境必须设超时，否则可能永久卡死）
	client := &http.Client{Timeout: 5 * time.Second}
	
	resp, err := client.Get(targetURL)
	if err != nil {
		fmt.Println("抓取失败:", err)
		return
	}
	// ★ 最佳实践：网络请求拿到响应体后，必须关闭，防止连接泄漏
	defer resp.Body.Close() 

	// 简单读取响应内容（真实爬虫会用 goquery 等库解析 HTML）
	// 这里我们只读取一部分来模拟
	buffer := make([]byte, 512)
	n, _ := resp.Body.Read(buffer)
	content := string(buffer[:n])
	
	// 模拟提取网页标题（真实场景用正则或 HTML 解析器）
	title := "HTTPBIN 示例页面"

	fmt.Println("网页抓取成功！")

	// =====================================================================
	// 知识点 2：时间处理
	// =====================================================================
	// 获取当前时间，并格式化输出
	now := time.Now()
	// Go 的时间格式化非常特别，必须用 "2006-01-02 15:04:05" 这个魔法数字
	timeStr := now.Format("2006-01-02 15:04:05")
	
	// 时间计算：如果想知道 1 小时后是几点
	oneHourLater := now.Add(1 * time.Hour)
	fmt.Printf("【时间处理】抓取时间: %s, 预计一小时后: %s\n", timeStr, oneHourLater.Format("15:04:05"))

	// =====================================================================
	// 知识点 3：加密哈希（crypto包）
	// =====================================================================
	// 哈希函数常用于验证文件完整性（内容有一点变化，哈希值完全不同）
	// 或者生成唯一指纹（如秒传功能）
	
	// 计算 MD5 (128位)
	md5Hash := md5.Sum([]byte(content))
	// 把字节数组转成十六进制字符串，方便阅读
	md5Str := hex.EncodeToString(md5Hash[:])
	
	// 计算 SHA256 (256位，更安全)
	sha256Hash := sha256.Sum256([]byte(content))
	sha256Str := hex.EncodeToString(sha256Hash[:])
	
	fmt.Printf("【加密哈希】MD5: %s\n", md5Str)
	fmt.Printf("【加密哈希】SHA256: %s\n", sha256Str)

	// =====================================================================
	// 知识点 4：JSON 编解码（序列化）
	// =====================================================================
	// 将刚才抓取的数据和计算结果打包成一个结构体
	pageData := WebPage{
		URL:        targetURL,
		Title:      title,
		Content:    content,
		HashMD5:    md5Str,
		HashSHA256: sha256Str,
		FetchTime:  timeStr,
	}

	// 使用 json.Marshal 将结构体序列化为 JSON 字节数组
	// 注意：如果结构体字段名是大写，默认会被转换成小写字母开头的 JSON 键；
	// 我们在结构体里用了 `json:"url"` 这样的 tag，可以自定义键名。
	jsonBytes, err := json.MarshalIndent(pageData, "", "  ") // MarshalIndent 带缩进，更美观
	if err != nil {
		fmt.Println("JSON 编码失败:", err)
		return
	}
	fmt.Println("【JSON 编解码】序列化结果 (预览前 100 字节):")
	fmt.Println(string(jsonBytes[:100]) + "...")

	// =====================================================================
	// 知识点 5：文件操作（写入磁盘）
	// =====================================================================
	// 将序列化后的 JSON 数据保存到本地文件
	fileName := "webpage_data.json"
	
	// 1. 创建文件（如果存在则覆盖）
	file, err := os.Create(fileName)
	if err != nil {
		fmt.Println("创建文件失败:", err)
		return
	}
	// ★ 最佳实践：打开文件后立刻 defer 关闭
	defer file.Close()

	// 2. 把 JSON 数据写入文件
	_, err = file.Write(jsonBytes)
	if err != nil {
		fmt.Println("写入文件失败:", err)
		return
	}
	fmt.Printf("【文件操作】数据已成功保存到 %s\n", fileName)
	
	// 3. 演示文件读取（可选）
	readBytes, err := os.ReadFile(fileName)
	if err == nil {
		fmt.Println("【文件操作】读取文件成功，前 50 字节内容如下：")
		fmt.Println(string(readBytes[:50]) + "...")
	}

	// =====================================================================
	// 知识点 6：补充（JSON 解码）
	// =====================================================================
	// 如果要把刚才存的文件读回来，变成 Go 结构体
	var parsedPage WebPage
	err = json.Unmarshal(readBytes, &parsedPage)
	if err == nil {
		fmt.Printf("【JSON 解码】反序列化后的标题: %s\n", parsedPage.Title)
	}

	fmt.Println("\n=== 爬虫任务结束 ===")
}
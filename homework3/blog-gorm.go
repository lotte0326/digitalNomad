package main

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// ============================================================
// 模型定义
// ============================================================

type User struct {
	gorm.Model
	Name  string
	Posts []Post
}

type Post struct {
	gorm.Model
	Title    string
	Content  string
	UserID   uint
	Comments []Comment
}

type Comment struct {
	gorm.Model
	Content string
	PostID  uint
}

// ============================================================
// 显式创建 SQLite 数据库文件（幂等）
// ============================================================

// ensureDatabaseFile 确保数据库文件存在
// 幂等：文件已存在则直接返回，不覆盖、不清空
func ensureDatabaseFile(dbPath string) error {
	// O_CREATE：不存在则创建；O_RDWR：读写打开
	// 若文件已存在，O_CREATE 不会截断（无 O_TRUNC），内容保留
	f, err := os.OpenFile(dbPath, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return fmt.Errorf("创建数据库文件失败: %w", err)
	}
	defer f.Close()

	// 判断是新建还是已存在，仅用于输出提示
	info, err := f.Stat()
	if err == nil && info.Size() == 0 {
		fmt.Printf("✅ 已创建新的数据库文件: %s\n", dbPath)
	} else {
		fmt.Printf("✅ 数据库文件已存在，直接复用: %s\n", dbPath)
	}
	return nil
}

// ============================================================
// 建表（幂等）
// ============================================================

func initTables(db *gorm.DB) error {
	if err := db.AutoMigrate(&User{}, &Post{}, &Comment{}); err != nil {
		return fmt.Errorf("迁移失败: %w", err)
	}
	fmt.Println("✅ 数据表 users / posts / comments 已就绪（幂等）")
	return nil
}

// ============================================================
// 主流程
// ============================================================

func main() {
	const dbPath = "blog.db"

	// 1. 显式创建数据库文件（幂等）
	if err := ensureDatabaseFile(dbPath); err != nil {
		log.Fatal(err)
	}

	// 2. 连接数据库（文件已存在）
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		log.Fatal("连接数据库失败:", err)
	}

	// 3. 建表（幂等）
	if err := initTables(db); err != nil {
		log.Fatal(err)
	}

	// 4. 验证表
	var tables []string
	db.Raw("SELECT name FROM sqlite_master WHERE type='table' AND name NOT LIKE 'sqlite_%'").
		Scan(&tables)
	fmt.Println("当前数据库中的表:", tables)

	// 5. 幂等写入示例数据
	user := User{Name: "张三"}
	if err := db.FirstOrCreate(&user, User{Name: "张三"}).Error; err != nil {
		log.Fatal("写入示例数据失败:", err)
	}
	var count int64
	db.Model(&User{}).Count(&count)
	fmt.Printf("users 表记录数: %d\n", count)
}
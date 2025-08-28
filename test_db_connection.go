package main

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v3"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Config 配置结构
type Config struct {
	Settings struct {
		Database struct {
			Driver string `yaml:"driver"`
			Source string `yaml:"source"`
		} `yaml:"database"`
	} `yaml:"settings"`
}

func main() {
	fmt.Println("=== MySQL 数据库连接测试工具 ===")

	// 读取配置文件
	configData, err := os.ReadFile("config/settings.yml")
	if err != nil {
		log.Fatalf("读取配置文件失败: %v", err)
	}

	// 解析配置
	var config Config
	err = yaml.Unmarshal(configData, &config)
	if err != nil {
		log.Fatalf("解析配置文件失败: %v", err)
	}

	fmt.Printf("数据库类型: %s\n", config.Settings.Database.Driver)
	fmt.Printf("连接字符串: %s\n", maskPassword(config.Settings.Database.Source))

	// 尝试连接数据库
	fmt.Println("\n正在测试数据库连接...")
	db, err := gorm.Open(mysql.Open(config.Settings.Database.Source), &gorm.Config{})
	if err != nil {
		fmt.Printf("❌ 数据库连接失败: %v\n", err)
		fmt.Println("\n解决建议:")
		fmt.Println("1. 检查MySQL服务是否运行")
		fmt.Println("2. 验证用户名和密码是否正确")
		fmt.Println("3. 确认数据库名是否存在")
		fmt.Println("4. 检查防火墙设置")
		os.Exit(1)
	}

	// 测试数据库连接
	sqlDB, err := db.DB()
	if err != nil {
		fmt.Printf("❌ 获取数据库实例失败: %v\n", err)
		os.Exit(1)
	}

	err = sqlDB.Ping()
	if err != nil {
		fmt.Printf("❌ 数据库Ping测试失败: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("✅ 数据库连接测试成功！")
	fmt.Println("您可以继续运行数据库迁移命令:")
	fmt.Println("go run main.go migrate -c config/settings.yml")
}

// maskPassword 隐藏密码信息
func maskPassword(source string) string {
	if len(source) <= 8 {
		return "***"
	}
	return source[:4] + "***" + source[len(source)-4:]
}

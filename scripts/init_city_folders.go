package main

import (
	"fmt"
	"log"
	"rentPro/rentpro-admin/common/database"
	"rentPro/rentpro-admin/common/initialize"
	"rentPro/rentpro-admin/common/utils"
)

func main() {
	fmt.Println("=== 楼盘管理城市文件夹初始化工具 ===")

	// 1. 初始化数据库连接
	fmt.Println("🔗 初始化数据库连接...")
	database.Setup()

	// 2. 初始化七牛云服务
	fmt.Println("☁️  初始化七牛云服务...")
	err := initialize.InitQiniu("dev")
	if err != nil {
		log.Fatalf("❌ 七牛云服务初始化失败: %v", err)
	}

	// 3. 初始化图片管理器
	fmt.Println("🖼️  初始化图片管理器...")
	err = utils.InitImageManager()
	if err != nil {
		log.Fatalf("❌ 图片管理器初始化失败: %v", err)
	}

	// 4. 获取图片管理器实例
	imageManager := utils.GetImageManager()
	if imageManager == nil {
		log.Fatal("❌ 获取图片管理器实例失败")
	}

	// 5. 初始化城市文件夹结构
	fmt.Println("🏙️  开始初始化城市文件夹结构...")
	err = imageManager.InitializeCityFolders()
	if err != nil {
		log.Fatalf("❌ 城市文件夹初始化失败: %v", err)
	}

	fmt.Println("\n🎉 楼盘管理城市文件夹初始化完成！")
	fmt.Println("📁 文件夹结构:")
	fmt.Println("   楼盘管理/")
	fmt.Println("   ├── 北京/")
	fmt.Println("   ├── 上海/")
	fmt.Println("   ├── 广州/")
	fmt.Println("   ├── 深圳/")
	fmt.Println("   └── ...")
	fmt.Println("=== 初始化完成 ===")
}

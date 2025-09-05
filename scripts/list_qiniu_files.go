package main

import (
	"fmt"
	"log"

	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
)

func main() {
	fmt.Println("=== 七牛云文件夹结构查看工具 ===")

	// 直接使用配置信息（从配置文件中获取）
	accessKey := "YRVDLN0uWnICy6OtDyWFy2OgJlUGQI2tMxGU13-z"
	secretKey := "6VAEYJIads4_EI5zjGhD0zSmy_d2IEdVmQCo3DLd"
	bucket := "rentpro-floor-plans"

	fmt.Printf("🔍 查看七牛云存储空间: %s\n", bucket)

	// 创建认证对象
	mac := qbox.NewMac(accessKey, secretKey)

	// 配置存储区域（华北）
	cfg := storage.Config{
		Zone:          &storage.ZoneHuabei,
		UseHTTPS:      true,
		UseCdnDomains: false,
	}

	bucketManager := storage.NewBucketManager(mac, &cfg)

	// 列出所有文件
	limit := 1000
	prefix := ""
	delimiter := ""
	marker := ""

	totalFiles := 0
	batchCount := 0

	fmt.Println("📋 扫描存储空间中的文件...")

	for {
		entries, _, nextMarker, hasNext, err := bucketManager.ListFiles(bucket, prefix, delimiter, marker, limit)
		if err != nil {
			log.Printf("❌ 列出文件失败: %v", err)
			break
		}

		if len(entries) == 0 {
			fmt.Println("📭 存储空间中没有文件")
			break
		}

		batchCount++
		totalFiles += len(entries)
		fmt.Printf("\n📁 第 %d 批：发现 %d 个文件\n", batchCount, len(entries))

		// 显示文件列表
		for _, entry := range entries {
			fmt.Printf("📄  %s (大小: %d bytes, 修改时间: %d)\n", entry.Key, entry.Fsize, entry.PutTime)
		}

		if !hasNext {
			break
		}
		marker = nextMarker
	}

	fmt.Printf("\n📊 统计信息:\n")
	fmt.Printf("   - 总共扫描批次: %d\n", batchCount)
	fmt.Printf("   - 总共文件数量: %d 个\n", totalFiles)

	if totalFiles == 0 {
		fmt.Println("📭 存储空间是空的")
	} else {
		fmt.Println("✅ 文件夹结构查看完成")
	}

	fmt.Println("=== 查看结束 ===")
}

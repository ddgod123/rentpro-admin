package main

import (
	"fmt"
	"log"

	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
)

func main() {
	fmt.Println("=== 七牛云存储清理工具 ===")

	// 直接使用配置信息（从配置文件中获取）
	accessKey := "YRVDLN0uWnICy6OtDyWFy2OgJlUGQI2tMxGU13-z"
	secretKey := "6VAEYJIads4_EI5zjGhD0zSmy_d2IEdVmQCo3DLd"
	bucket := "rentpro-floor-plans"

	fmt.Printf("🔍 开始清理七牛云存储空间: %s\n", bucket)

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

	totalDeleted := 0
	totalFiles := 0
	batchCount := 0

	fmt.Println("📋 正在扫描存储空间中的文件...")

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
		fmt.Printf("\n📁 第 %d 批：发现 %d 个文件，准备删除...\n", batchCount, len(entries))

		// 逐个删除文件（避免批量删除的复杂性）
		successCount := 0
		for _, entry := range entries {
			fmt.Printf("🗑️  删除: %s", entry.Key)

			err := bucketManager.Delete(bucket, entry.Key)
			if err != nil {
				fmt.Printf(" ❌ 失败: %v\n", err)
			} else {
				fmt.Printf(" ✅ 成功\n")
				successCount++
			}
		}

		totalDeleted += successCount
		fmt.Printf("✅ 第 %d 批删除完成: 成功 %d 个，失败 %d 个\n", batchCount, successCount, len(entries)-successCount)

		if !hasNext {
			break
		}
		marker = nextMarker
	}

	fmt.Printf("\n🎉 清理完成！\n")
	fmt.Printf("📊 统计信息:\n")
	fmt.Printf("   - 总共扫描批次: %d\n", batchCount)
	fmt.Printf("   - 总共发现文件: %d 个\n", totalFiles)
	fmt.Printf("   - 成功删除文件: %d 个\n", totalDeleted)
	fmt.Printf("   - 删除失败文件: %d 个\n", totalFiles-totalDeleted)

	if totalDeleted == totalFiles && totalFiles > 0 {
		fmt.Println("🎊 所有文件已成功清理！")
	} else if totalFiles == 0 {
		fmt.Println("📭 存储空间本来就是空的")
	}

	fmt.Println("=== 清理结束 ===")
}

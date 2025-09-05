package main

import (
	"fmt"
	"log"
	"rentPro/rentpro-admin/common/config"

	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
)

func main() {
	fmt.Println("=== 七牛云存储清理工具 ===")

	// 加载七牛云配置
	qiniuConfig := config.GetQiniuConfig()
	if qiniuConfig == nil {
		log.Fatal("❌ 七牛云配置加载失败")
	}

	fmt.Printf("🔍 开始清理七牛云存储空间: %s\n", qiniuConfig.Bucket)

	// 创建认证对象
	mac := qbox.NewMac(qiniuConfig.AccessKey, qiniuConfig.SecretKey)

	// 配置存储区域
	cfg := storage.Config{
		UseHTTPS:      qiniuConfig.UseHTTPS,
		UseCdnDomains: qiniuConfig.UseCdnDomains,
	}

	// 设置存储区域
	switch qiniuConfig.Zone {
	case "huadong":
		cfg.Zone = &storage.ZoneHuadong
	case "huabei":
		cfg.Zone = &storage.ZoneHuabei
	case "huanan":
		cfg.Zone = &storage.ZoneHuanan
	case "beimei":
		cfg.Zone = &storage.ZoneBeimei
	case "xinjiapo":
		cfg.Zone = &storage.ZoneXinjiapo
	default:
		cfg.Zone = &storage.ZoneHuabei // 默认华北
	}

	bucketManager := storage.NewBucketManager(mac, &cfg)

	// 列出所有文件
	limit := 1000
	prefix := ""
	delimiter := ""
	marker := ""

	totalDeleted := 0
	totalFiles := 0

	fmt.Println("📋 正在扫描存储空间中的文件...")

	for {
		entries, _, nextMarker, hasNext, err := bucketManager.ListFiles(qiniuConfig.Bucket, prefix, delimiter, marker, limit)
		if err != nil {
			log.Printf("❌ 列出文件失败: %v", err)
			break
		}

		if len(entries) == 0 {
			break
		}

		totalFiles += len(entries)
		fmt.Printf("📁 发现 %d 个文件，准备删除...\n", len(entries))

		// 批量删除文件
		deleteOps := make([]string, 0, len(entries))
		for _, entry := range entries {
			deleteOps = append(deleteOps, storage.URIDelete(qiniuConfig.Bucket, entry.Key))
			fmt.Printf("🗑️  %s\n", entry.Key)
		}

		// 执行批量删除
		rets, err := bucketManager.Batch(deleteOps)
		if err != nil {
			log.Printf("❌ 批量删除失败: %v", err)
		} else {
			successCount := 0
			for i, ret := range rets {
				if ret.Code == 200 {
					successCount++
				} else {
					fmt.Printf("⚠️  删除失败 %s: 状态码 %d\n", entries[i].Key, ret.Code)
				}
			}
			totalDeleted += successCount
			fmt.Printf("✅ 本批次删除成功: %d 个文件\n", successCount)
		}

		if !hasNext {
			break
		}
		marker = nextMarker
	}

	fmt.Printf("\n🎉 清理完成！\n")
	fmt.Printf("📊 统计信息:\n")
	fmt.Printf("   - 总共发现: %d 个文件\n", totalFiles)
	fmt.Printf("   - 成功删除: %d 个文件\n", totalDeleted)
	fmt.Printf("   - 删除失败: %d 个文件\n", totalFiles-totalDeleted)
	fmt.Println("=== 清理结束 ===")
}

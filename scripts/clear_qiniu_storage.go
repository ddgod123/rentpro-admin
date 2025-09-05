package main

import (
	"fmt"
	"log"
	"rentPro/rentpro-admin/common/global"
	"rentPro/rentpro-admin/common/initialize"

	"github.com/qiniu/go-sdk/v7/storage"
)

func main() {
	fmt.Println("=== 七牛云存储清理工具 ===")

	// 初始化七牛云服务
	initialize.InitializeQiniu()

	if global.QiniuService == nil {
		log.Fatal("❌ 七牛云服务初始化失败")
	}

	fmt.Printf("🔍 开始清理七牛云存储空间: %s\n", global.QiniuConfig.Bucket)

	// 获取七牛云配置
	mac := global.QiniuService.GetMac()
	cfg := storage.Config{
		Zone:          global.QiniuService.GetZone(),
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

	for {
		entries, _, nextMarker, hasNext, err := bucketManager.ListFiles(global.QiniuConfig.Bucket, prefix, delimiter, marker, limit)
		if err != nil {
			log.Printf("❌ 列出文件失败: %v", err)
			break
		}

		if len(entries) == 0 {
			break
		}

		// 批量删除文件
		deleteOps := make([]string, 0, len(entries))
		for _, entry := range entries {
			deleteOps = append(deleteOps, storage.URIDelete(global.QiniuConfig.Bucket, entry.Key))
			fmt.Printf("🗑️  准备删除: %s\n", entry.Key)
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
					fmt.Printf("⚠️  删除失败 %s: %s\n", entries[i].Key, ret.Error)
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

	fmt.Printf("\n🎉 清理完成！总共删除了 %d 个文件\n", totalDeleted)
	fmt.Println("=== 清理结束 ===")
}

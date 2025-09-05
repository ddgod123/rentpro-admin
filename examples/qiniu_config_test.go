package main
package main

import (
	"fmt"
	"log"

	"rentPro/rentpro-admin/common/config"
	"rentPro/rentpro-admin/common/initialize"
)

func main() {
	fmt.Println("开始测试七牛云配置...")

	// 测试从 settings.yml 加载配置
	fmt.Println("1. 测试从 settings.yml 加载七牛云配置...")
	qiniuConfig := config.GetQiniuConfig()
	if qiniuConfig != nil {
		fmt.Printf("   ✅ 成功加载配置\n")
		fmt.Printf("   存储空间: %s\n", qiniuConfig.Bucket)
		fmt.Printf("   访问域名: %s\n", qiniuConfig.Domain)
		fmt.Printf("   存储区域: %s\n", qiniuConfig.Zone)
		fmt.Printf("   使用HTTPS: %v\n", qiniuConfig.UseHTTPS)
		fmt.Printf("   使用CDN: %v\n", qiniuConfig.UseCdnDomains)
		fmt.Printf("   最大文件大小: %d bytes\n", qiniuConfig.Upload.MaxFileSize)
		fmt.Printf("   支持的文件类型: %v\n", qiniuConfig.Upload.AllowedTypes)
	} else {
		fmt.Printf("   ❌ 未能从 settings.yml 加载七牛云配置\n")
	}

	// 测试从 qiniu.yml 加载配置
	fmt.Println("\n2. 测试从 qiniu.yml 加载七牛云配置...")
	err := initialize.InitQiniu("development")
	if err != nil {
		log.Printf("   ❌ 初始化七牛云配置失败: %v", err)
	} else {
		fmt.Printf("   ✅ 成功从 qiniu.yml 初始化配置\n")
		qiniuConfig = config.GetQiniuConfig()
		if qiniuConfig != nil {
			fmt.Printf("   存储空间: %s\n", qiniuConfig.Bucket)
			fmt.Printf("   访问域名: %s\n", qiniuConfig.Domain)
			fmt.Printf("   存储区域: %s\n", qiniuConfig.Zone)
			fmt.Printf("   使用HTTPS: %v\n", qiniuConfig.UseHTTPS)
			fmt.Printf("   使用CDN: %v\n", qiniuConfig.UseCdnDomains)
			fmt.Printf("   最大文件大小: %d bytes\n", qiniuConfig.Upload.MaxFileSize)
			fmt.Printf("   支持的文件类型: %v\n", qiniuConfig.Upload.AllowedTypes)
		}
	}

	fmt.Println("\n测试完成!")
}
package initialize

import (
	"fmt"
	"log"
	"path/filepath"

	"rentPro/rentpro-admin/common/config"
	"rentPro/rentpro-admin/common/utils"
)

// InitQiniu 初始化七牛云服务
func InitQiniu(env string) error {
	log.Println("开始初始化七牛云服务...")

	// 配置文件路径
	configPath := filepath.Join("config", "qiniu.yml")

	// 初始化七牛云配置
	err := config.InitQiniuConfig(configPath, env)
	if err != nil {
		return fmt.Errorf("初始化七牛云配置失败: %v", err)
	}

	// 初始化七牛云服务
	err = utils.InitQiniuService()
	if err != nil {
		return fmt.Errorf("初始化七牛云服务失败: %v", err)
	}

	// 验证服务可用性
	qiniuService := utils.GetQiniuService()
	if qiniuService == nil {
		return fmt.Errorf("七牛云服务初始化失败")
	}

	log.Println("✅ 七牛云服务初始化成功！")
	
	// 输出配置信息（脱敏）
	qiniuConfig := config.GetQiniuConfig()
	if qiniuConfig != nil {
		log.Printf("存储空间: %s", qiniuConfig.Bucket)
		log.Printf("访问域名: %s", qiniuConfig.Domain)
		log.Printf("存储区域: %s", qiniuConfig.Zone)
		log.Printf("使用HTTPS: %v", qiniuConfig.UseHTTPS)
		log.Printf("使用CDN: %v", qiniuConfig.UseCdnDomains)
		log.Printf("最大文件大小: %d bytes", qiniuConfig.Upload.MaxFileSize)
		log.Printf("支持的文件类型: %v", qiniuConfig.Upload.AllowedTypes)
	}

	return nil
}

// CheckQiniuHealth 检查七牛云服务健康状态
func CheckQiniuHealth() error {
	qiniuService := utils.GetQiniuService()
	if qiniuService == nil {
		return fmt.Errorf("七牛云服务未初始化")
	}

	// 尝试列出文件来验证连接
	_, err := qiniuService.ListFiles("", 1)
	if err != nil {
		return fmt.Errorf("七牛云服务连接失败: %v", err)
	}

	return nil
}

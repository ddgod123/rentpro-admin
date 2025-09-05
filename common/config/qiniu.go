package config

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

// QiniuConfig 七牛云配置结构
type QiniuConfig struct {
	AccessKey     string                `yaml:"access_key"`      // Access Key
	SecretKey     string                `yaml:"secret_key"`      // Secret Key
	Bucket        string                `yaml:"bucket"`          // 存储空间名称
	Domain        string                `yaml:"domain"`          // 访问域名
	Zone          string                `yaml:"zone"`            // 存储区域
	UseHTTPS      bool                  `yaml:"use_https"`       // 是否使用HTTPS
	UseCdnDomains bool                  `yaml:"use_cdn_domains"` // 是否使用CDN域名
	Upload        QiniuUploadConfig     `yaml:"upload"`          // 上传配置
	ImageStyles   map[string]ImageStyle `yaml:"image_styles"`    // 图片样式配置
}

// QiniuUploadConfig 上传配置
type QiniuUploadConfig struct {
	MaxFileSize  int64    `yaml:"max_file_size"` // 最大文件大小
	AllowedTypes []string `yaml:"allowed_types"` // 允许的文件类型
	UploadDir    string   `yaml:"upload_dir"`    // 上传目录前缀
}

// ImageStyle 图片样式配置
type ImageStyle struct {
	Name        string `yaml:"name"`        // 样式名称
	Process     string `yaml:"process"`     // 处理参数
	Description string `yaml:"description"` // 样式描述
}

// EnvironmentConfig 环境配置
type EnvironmentConfig struct {
	Qiniu QiniuConfig `yaml:"qiniu"`
}

// QiniuConfigManager 七牛云配置管理器
type QiniuConfigManager struct {
	config     *QiniuConfig
	configPath string
	env        string
}

// NewQiniuConfigManager 创建配置管理器
func NewQiniuConfigManager(configPath string, env string) (*QiniuConfigManager, error) {
	manager := &QiniuConfigManager{
		configPath: configPath,
		env:        env,
	}

	err := manager.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("加载七牛云配置失败: %v", err)
	}

	return manager, nil
}

// LoadConfig 加载配置文件
func (m *QiniuConfigManager) LoadConfig() error {
	// 读取配置文件
	data, err := os.ReadFile(m.configPath)
	if err != nil {
		return fmt.Errorf("读取配置文件失败: %v", err)
	}

	// 解析YAML
	var fullConfig map[string]interface{}
	err = yaml.Unmarshal(data, &fullConfig)
	if err != nil {
		return fmt.Errorf("解析YAML配置失败: %v", err)
	}

	// 根据环境选择配置
	var targetConfig map[string]interface{}

	if m.env == "development" || m.env == "production" {
		// 使用环境特定配置
		if envConfig, exists := fullConfig[m.env]; exists {
			if envMap, ok := envConfig.(map[string]interface{}); ok {
				if qiniuConfig, exists := envMap["qiniu"]; exists {
					targetConfig = qiniuConfig.(map[string]interface{})
				}
			}
		}
	} else {
		// 使用默认配置
		if qiniuConfig, exists := fullConfig["qiniu"]; exists {
			targetConfig = qiniuConfig.(map[string]interface{})
		}
	}

	if targetConfig == nil {
		return fmt.Errorf("未找到七牛云配置")
	}

	// 转换为YAML并解析到结构体
	configData, err := yaml.Marshal(targetConfig)
	if err != nil {
		return fmt.Errorf("序列化配置失败: %v", err)
	}

	var config QiniuConfig
	err = yaml.Unmarshal(configData, &config)
	if err != nil {
		return fmt.Errorf("解析七牛云配置失败: %v", err)
	}

	// 处理环境变量替换
	config.AccessKey = m.expandEnvVar(config.AccessKey)
	config.SecretKey = m.expandEnvVar(config.SecretKey)
	config.Domain = m.expandEnvVar(config.Domain)

	m.config = &config
	return nil
}

// expandEnvVar 展开环境变量
func (m *QiniuConfigManager) expandEnvVar(value string) string {
	if strings.HasPrefix(value, "${") && strings.HasSuffix(value, "}") {
		envName := value[2 : len(value)-1]
		if envValue := os.Getenv(envName); envValue != "" {
			return envValue
		}
	}
	return value
}

// GetConfig 获取配置
func (m *QiniuConfigManager) GetConfig() *QiniuConfig {
	return m.config
}

// ValidateConfig 验证配置
func (m *QiniuConfigManager) ValidateConfig() error {
	if m.config == nil {
		return fmt.Errorf("配置未加载")
	}

	if m.config.AccessKey == "" || m.config.AccessKey == "your_access_key_here" {
		return fmt.Errorf("Access Key 未配置")
	}

	if m.config.SecretKey == "" || m.config.SecretKey == "your_secret_key_here" {
		return fmt.Errorf("Secret Key 未配置")
	}

	if m.config.Bucket == "" {
		return fmt.Errorf("存储空间名称未配置")
	}

	if m.config.Domain == "" || m.config.Domain == "your-domain.com" {
		return fmt.Errorf("访问域名未配置")
	}

	return nil
}

// GetImageStyleURL 获取图片样式URL
func (m *QiniuConfigManager) GetImageStyleURL(baseURL string, styleName string) string {
	if style, exists := m.config.ImageStyles[styleName]; exists {
		return fmt.Sprintf("%s-%s", baseURL, style.Name)
	}
	return baseURL
}

// GetUploadKey 生成上传Key
func (m *QiniuConfigManager) GetUploadKey(filename string) string {
	if m.config.Upload.UploadDir != "" {
		return fmt.Sprintf("%s/%s", m.config.Upload.UploadDir, filename)
	}
	return filename
}

// IsAllowedFileType 检查文件类型是否允许
func (m *QiniuConfigManager) IsAllowedFileType(contentType string) bool {
	for _, allowedType := range m.config.Upload.AllowedTypes {
		if allowedType == contentType {
			return true
		}
	}
	return false
}

// GetMaxFileSize 获取最大文件大小
func (m *QiniuConfigManager) GetMaxFileSize() int64 {
	return m.config.Upload.MaxFileSize
}

// GetPublicURL 生成公共访问URL
func (m *QiniuConfigManager) GetPublicURL(key string) string {
	protocol := "http"
	if m.config.UseHTTPS {
		protocol = "https"
	}
	return fmt.Sprintf("%s://%s/%s", protocol, m.config.Domain, key)
}

// 全局配置实例
var QiniuConfigInstance *QiniuConfigManager

// InitQiniuConfig 初始化七牛云配置
func InitQiniuConfig(configPath string, env string) error {
	var err error
	QiniuConfigInstance, err = NewQiniuConfigManager(configPath, env)
	if err != nil {
		return err
	}

	// 验证配置
	err = QiniuConfigInstance.ValidateConfig()
	if err != nil {
		return fmt.Errorf("七牛云配置验证失败: %v", err)
	}

	return nil
}

// GetQiniuConfig 获取全局七牛云配置
func GetQiniuConfig() *QiniuConfig {
	if QiniuConfigInstance == nil {
		return nil
	}
	return QiniuConfigInstance.GetConfig()
}

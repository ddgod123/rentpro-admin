// Package config 提供配置信息相关的命令行功能
// 用于显示和验证 rentpro-admin 系统的配置信息
package config

import (
	"fmt"
	"os"
	"path/filepath"
	"rentPro/rentpro-admin/common/global"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

// 配置数据结构
type Config struct {
	Settings struct {
		Application struct {
			Mode         string `yaml:"mode"`
			Host         string `yaml:"host"`
			Name         string `yaml:"name"`
			Port         int    `yaml:"port"`
			ReadTimeout  int    `yaml:"readtimeout"`
			WriteTimeout int    `yaml:"writertimeout"`
			EnabledDP    bool   `yaml:"enabledp"`
		} `yaml:"application"`
		Logger struct {
			Path      string `yaml:"path"`
			Stdout    string `yaml:"stdout"`
			Level     string `yaml:"level"`
			EnabledDB bool   `yaml:"enableddb"`
		} `yaml:"logger"`
		JWT struct {
			Secret  string `yaml:"secret"`
			Timeout int    `yaml:"timeout"`
		} `yaml:"jwt"`
		Database struct {
			Driver string `yaml:"driver"`
			Source string `yaml:"source"`
		} `yaml:"database"`
	} `yaml:"settings"`
}

var (
	// configYml 配置文件路径
	configYml string

	// showVersion 显示版本信息
	showVersion bool

	// StartCmd 定义了 config 子命令
	// 用于显示和验证系统配置信息
	StartCmd = &cobra.Command{
		Use:     "config",
		Short:   "获取应用程序配置信息",
		Long:    `显示 rentpro-admin 系统的完整配置信息，包括应用设置、数据库配置、日志配置等`,
		Example: "rentpro-admin config -c config/settings.yml",

		// PreRun 在实际命令执行前运行
		PreRun: func(cmd *cobra.Command, args []string) {
			// 如果指定了版本标志，直接显示版本信息
			if showVersion {
				fmt.Printf("rentpro-admin config version: %s\n", global.Version)
				return
			}
		},

		// RunE 是命令的主要执行函数
		RunE: func(cmd *cobra.Command, args []string) error {
			// 如果只是查看版本，不执行配置显示
			if showVersion {
				return nil
			}

			return run()
		},
	}
)

// init 初始化命令标志
func init() {
	// 添加版本标志支持
	StartCmd.Flags().BoolVarP(&showVersion, "version", "v", false, "显示版本信息")

	// 添加配置文件标志
	StartCmd.PersistentFlags().StringVarP(&configYml, "config", "c", "config/settings.yml", "指定配置文件路径")
}

// run 执行配置信息显示的核心逻辑
// 读取指定的 YAML 配置文件并显示内容
func run() error {
	fmt.Printf("=== rentpro-admin 配置信息 v%s ===\n", global.Version)
	fmt.Printf("配置文件: %s\n", configYml)

	// 检查配置文件是否存在
	if !fileExists(configYml) {
		return fmt.Errorf("配置文件不存在: %s", configYml)
	}

	// 读取配置文件
	configData, err := os.ReadFile(configYml)
	if err != nil {
		return fmt.Errorf("读取配置文件失败: %v", err)
	}

	// 解析 YAML 配置
	var config Config
	err = yaml.Unmarshal(configData, &config)
	if err != nil {
		return fmt.Errorf("解析配置文件失败: %v", err)
	}

	// 显示配置信息
	fmt.Println("\n=== 应用程序配置 ===")
	fmt.Printf("运行模式: %s\n", config.Settings.Application.Mode)
	fmt.Printf("绑定地址: %s:%d\n", config.Settings.Application.Host, config.Settings.Application.Port)
	fmt.Printf("应用名称: %s\n", config.Settings.Application.Name)
	fmt.Printf("数据权限: %t\n", config.Settings.Application.EnabledDP)

	fmt.Println("\n=== 日志配置 ===")
	fmt.Printf("日志路径: %s\n", config.Settings.Logger.Path)
	fmt.Printf("日志级别: %s\n", config.Settings.Logger.Level)
	fmt.Printf("数据库日志: %t\n", config.Settings.Logger.EnabledDB)

	fmt.Println("\n=== JWT 配置 ===")
	fmt.Printf("密钥: %s\n", maskSensitiveInfo(config.Settings.JWT.Secret))
	fmt.Printf("过期时间: %d秒\n", config.Settings.JWT.Timeout)

	fmt.Println("\n=== 数据库配置 ===")
	fmt.Printf("数据库类型: %s\n", config.Settings.Database.Driver)
	fmt.Printf("连接字符串: %s\n", maskSensitiveInfo(config.Settings.Database.Source))

	fmt.Println("\n✅ 配置信息显示完成！")

	return nil
}

// fileExists 检查文件是否存在
func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

// maskSensitiveInfo 隐藏敏感信息，用于安全显示
func maskSensitiveInfo(info string) string {
	if len(info) <= 8 {
		return "***"
	}
	return info[:4] + "***" + info[len(info)-4:]
}

// GetConfigVersion 获取配置工具版本信息
// 便于测试和版本管理
func GetConfigVersion() string {
	return global.Version
}

// ValidateConfigFile 验证配置文件是否有效
// 用于配置文件验证和调试
func ValidateConfigFile(configPath string) error {
	if !fileExists(configPath) {
		return fmt.Errorf("配置文件不存在: %s", configPath)
	}

	// 检查文件扩展名
	ext := filepath.Ext(configPath)
	if ext != ".yml" && ext != ".yaml" {
		return fmt.Errorf("不支持的配置文件格式: %s，只支持 .yml 或 .yaml", ext)
	}

	// 尝试解析配置文件
	configData, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("读取配置文件失败: %v", err)
	}

	var config Config
	err = yaml.Unmarshal(configData, &config)
	if err != nil {
		return fmt.Errorf("解析配置文件失败: %v", err)
	}

	return nil
}

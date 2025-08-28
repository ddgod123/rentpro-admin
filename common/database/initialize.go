// Package database 提供数据库初始化和配置管理功能
// 用于 rentpro-admin 租赁管理系统的数据库连接和 ORM 配置
package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"rentPro/rentpro-admin/common/global"

	"gopkg.in/yaml.v3"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

// DatabaseConfig 数据库配置结构体
// 用于解析 settings.yml 中的数据库配置信息
type DatabaseConfig struct {
	Settings struct {
		Database struct {
			// Driver 数据库驱动类型（mysql, sqlite3, postgres 等）
			Driver string `yaml:"driver"`
			// Source 数据库连接字符串
			Source string `yaml:"source"`
		} `yaml:"database"`
	} `yaml:"settings"`
}

// DB 全局数据库实例
// 用于在整个应用中共享数据库连接
var DB *gorm.DB

// Setup 配置和初始化数据库连接
// 使用项目内部实现，不依赖外部框架
func Setup() {
	log.Printf("开始初始化数据库连接...")

	// 读取和解析配置文件
	config, err := loadDatabaseConfig("config/settings.yml")
	if err != nil {
		log.Fatalf("加载数据库配置失败: %v", err)
	}

	// 验证数据库配置
	if err := validateDatabaseConfig(config); err != nil {
		log.Fatalf("数据库配置验证失败: %v", err)
	}

	// 设置全局驱动类型
	global.Driver = config.Settings.Database.Driver

	// 创建数据库连接
	db, err := createDatabaseConnection(config)
	if err != nil {
		log.Fatalf("创建数据库连接失败: %v", err)
	}

	// 设置全局数据库实例
	DB = db

	// 测试数据库连接
	if err := testDatabaseConnection(db); err != nil {
		log.Fatalf("数据库连接测试失败: %v", err)
	}

	log.Printf("✅ 数据库连接初始化成功！驱动: %s", config.Settings.Database.Driver)
}

// loadDatabaseConfig 加载和解析数据库配置文件
func loadDatabaseConfig(configPath string) (*DatabaseConfig, error) {
	// 检查文件是否存在
	if !fileExists(configPath) {
		return nil, fmt.Errorf("配置文件不存在: %s", configPath)
	}

	// 读取配置文件
	configData, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %v", err)
	}

	// 解析 YAML 配置
	var config DatabaseConfig
	err = yaml.Unmarshal(configData, &config)
	if err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %v", err)
	}

	return &config, nil
}

// validateDatabaseConfig 验证数据库配置的有效性
func validateDatabaseConfig(config *DatabaseConfig) error {
	if config.Settings.Database.Driver == "" {
		return fmt.Errorf("数据库驱动类型不能为空")
	}

	if config.Settings.Database.Source == "" {
		return fmt.Errorf("数据库连接字符串不能为空")
	}

	// 验证支持的数据库类型
	supportedDrivers := []string{"mysql", "sqlite3", "postgres"}
	isSupported := false
	for _, driver := range supportedDrivers {
		if config.Settings.Database.Driver == driver {
			isSupported = true
			break
		}
	}

	if !isSupported {
		return fmt.Errorf("不支持的数据库类型: %s，支持的类型: %v",
			config.Settings.Database.Driver, supportedDrivers)
	}

	return nil
}

// createDatabaseConnection 创建数据库连接
func createDatabaseConnection(config *DatabaseConfig) (*gorm.DB, error) {
	log.Printf("正在连接数据库: %s", config.Settings.Database.Driver)
	log.Printf("连接字符串: %s", maskSensitiveInfo(config.Settings.Database.Source))

	// 配置 GORM 日志记录器
	gormLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Info,
			Colorful:      true,
		},
	)

	// 配置 GORM 设置
	gormConfig := &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		Logger: gormLogger,
	}

	// 根据数据库类型创建连接
	var db *gorm.DB
	var err error

	switch config.Settings.Database.Driver {
	case "mysql":
		db, err = gorm.Open(mysql.Open(config.Settings.Database.Source), gormConfig)
	default:
		return nil, fmt.Errorf("不支持的数据库类型: %s", config.Settings.Database.Driver)
	}

	if err != nil {
		return nil, fmt.Errorf("数据库连接失败: %v", err)
	}

	// 配置连接池参数
	if sqlDB, err := db.DB(); err == nil {
		sqlDB.SetMaxIdleConns(10)
		sqlDB.SetMaxOpenConns(100)
		sqlDB.SetConnMaxLifetime(time.Hour)
	} else {
		log.Printf("警告: 无法配置数据库连接池: %v", err)
	}

	return db, nil
}

// testDatabaseConnection 测试数据库连接是否正常
func testDatabaseConnection(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("获取数据库实例失败: %v", err)
	}

	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("数据库连接测试失败: %v", err)
	}

	log.Printf("✅ 数据库连接测试成功")
	return nil
}

// fileExists 检查文件是否存在
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// maskSensitiveInfo 隐藏敏感信息，用于安全显示
func maskSensitiveInfo(info string) string {
	if len(info) <= 8 {
		return "***"
	}
	return info[:4] + "***" + info[len(info)-4:]
}

// Package migrate 提供数据库迁移相关的命令行功能
// 用于管理 rentpro-admin 系统的数据库结构迁移
package migrate

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"rentPro/rentpro-admin/cmd/migrate/migration"
	_ "rentPro/rentpro-admin/cmd/migrate/migration/version-local"
	"rentPro/rentpro-admin/common/database"
	"rentPro/rentpro-admin/common/global"
	"rentPro/rentpro-admin/common/models/base"
)

// 配置数据结构，用于解析 settings.yml
type Config struct {
	Settings struct {
		Database struct {
			Driver string `yaml:"driver"`
			Source string `yaml:"source"`
		} `yaml:"database"`
	} `yaml:"settings"`
}

var (
	configYml   string
	generate    bool
	goAdmin     bool
	host        string
	showVersion bool

	// StartCmd 定义了 migrate 子命令
	// 用于执行数据库迁移操作，支持以下功能：
	// 1. 数据库表结构迁移和初始化
	// 2. 生成新的迁移文件
	// 3. 支持本地迁移管理
	// 命令注册：通过 rootCmd.AddCommand(migrate.StartCmd) 注册到根命令
	// 使用方式：
	//   - rentpro-admin migrate -c config/settings.yml  : 执行数据库迁移
	//   - rentpro-admin migrate -g                      : 生成迁移文件
	//   - rentpro-admin migrate -v                      : 显示版本信息
	// 版本信息来源：common/global/adm.go 中的 Version 常量
	StartCmd = &cobra.Command{
		Use:     "migrate",
		Short:   "数据库迁移工具",
		Long:    `rentpro-admin 数据库迁移工具，用于管理数据库结构的版本控制和自动迁移`,
		Example: "rentpro-admin migrate -c config/settings.yml",
		PreRun: func(cmd *cobra.Command, args []string) {
			if showVersion {
				fmt.Printf("rentpro-admin migrate version: %s\n", global.Version)
				return
			}
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if showVersion {
				return nil
			}
			return run()
		},
	}
)

// init 初始化命令标志
// 遵循CLI命令注释规范：包含命令注册机制说明、支持的使用方式、版本信息来源
func init() {
	// 添加版本标志支持，支持 -v 和 --version 两种形式
	StartCmd.PersistentFlags().BoolVarP(&showVersion, "version", "v", false, "显示版本信息")

	// 配置文件路径标志
	StartCmd.PersistentFlags().StringVarP(&configYml, "config", "c", "config/settings.yml", "指定配置文件路径")

	// 生成迁移文件标志
	StartCmd.PersistentFlags().BoolVarP(&generate, "generate", "g", false, "生成迁移文件")

	// go-admin 模式标志
	StartCmd.PersistentFlags().BoolVarP(&goAdmin, "goAdmin", "a", false, "生成 go-admin 框架迁移文件")

	// 多租户主机选择标志
	StartCmd.PersistentFlags().StringVarP(&host, "domain", "d", "*", "选择租户主机域名")
}

// run 执行数据库迁移的核心逻辑
func run() error {
	fmt.Printf("=== rentpro-admin 数据库迁移工具 v%s ===\n", global.Version)

	if !generate {
		fmt.Println("开始初始化数据库...")

		// 验证配置文件参数
		if configYml == "" {
			return fmt.Errorf("请指定配置文件路径，使用 -c 参数")
		}

		fmt.Printf("配置文件: %s\n", configYml)

		// 读取和解析配置文件
		config, err := loadConfig(configYml)
		if err != nil {
			return fmt.Errorf("加载配置文件失败: %v", err)
		}

		// 执行数据库初始化和迁移
		err = initDB(config)
		if err != nil {
			return fmt.Errorf("数据库初始化失败: %v", err)
		}

		fmt.Println("✅ 数据库初始化完成！")
	} else {
		fmt.Println("生成迁移文件...")
		err := genFile()
		if err != nil {
			return fmt.Errorf("生成迁移文件失败: %v", err)
		}
		fmt.Println("✅ 迁移文件生成完成！")
	}

	return nil
}

// loadConfig 加载和解析配置文件
func loadConfig(configPath string) (*Config, error) {
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
	var config Config
	err = yaml.Unmarshal(configData, &config)
	if err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %v", err)
	}

	return &config, nil
}

// initDB 执行数据库初始化和迁移
func initDB(config *Config) error {
	fmt.Printf("数据库类型: %s\n", config.Settings.Database.Driver)
	fmt.Printf("数据库连接: %s\n", maskSensitiveInfo(config.Settings.Database.Source))

	// 1. 初始化数据库连接
	fmt.Println("初始化数据库连接...")
	// 直接使用 go-admin 框架的数据库初始化
	database.Setup()

	// 2. 执行数据库迁移
	fmt.Println("数据库迁移开始")
	err := migrateModel()
	if err != nil {
		return fmt.Errorf("数据库迁移失败: %v", err)
	}

	fmt.Println("数据库基础数据初始化成功")
	return nil
}

// migrateModel 执行模型迁移
func migrateModel() error {
	// 获取数据库实例
	if database.DB == nil {
		return fmt.Errorf("数据库连接未初始化")
	}

	db := database.DB

	// 设置 MySQL 表选项
	if global.Driver == "mysql" {
		db = db.Set("gorm:table_options", "ENGINE=InnoDB CHARSET=utf8mb4")
	}

	fmt.Println("执行数据库迁移...")
	fmt.Println("- 创建迁移记录表...")

	// 自动迁移 Migration 模型
	err := db.Debug().AutoMigrate(&base.Migration{})
	if err != nil {
		return fmt.Errorf("迁移 Migration 模型失败: %v", err)
	}

	fmt.Println("- 执行业务表迁移...")

	// 设置迁移管理器的数据库连接
	migration.Migrate.SetDb(db.Debug())

	// 执行所有注册的迁移
	migration.Migrate.Migrate()

	fmt.Println("✅ 数据库迁移执行完成")
	return nil
}

// genFile 生成迁移文件
func genFile() error {
	// 创建迁移文件目录
	migrationDir := "cmd/migrate/migration/version-local"
	if goAdmin {
		migrationDir = "cmd/migrate/migration/version"
	}

	if !fileExists(migrationDir) {
		err := os.MkdirAll(migrationDir, 0755)
		if err != nil {
			return fmt.Errorf("创建迁移目录失败: %v", err)
		}
	}

	// 生成时间戳
	timestamp := strconv.FormatInt(time.Now().UnixNano()/1e6, 10)
	filename := filepath.Join(migrationDir, timestamp+"_migrate.go")

	// 生成迁移文件内容
	packageName := "version_local"
	if goAdmin {
		packageName = "version"
	}

	content := generateMigrationTemplate(packageName, timestamp)

	// 写入文件
	err := os.WriteFile(filename, []byte(content), 0644)
	if err != nil {
		return fmt.Errorf("写入迁移文件失败: %v", err)
	}

	fmt.Printf("迁移文件已生成: %s\n", filename)
	return nil
}

// generateMigrationTemplate 生成迁移文件模板
func generateMigrationTemplate(packageName, timestamp string) string {
	return fmt.Sprintf(`package %s

import (
	"rentPro/rentpro-admin/cmd/migrate/migration"
	"rentPro/rentpro-admin/common/models/base"
	"gorm.io/gorm"
)

func init() {
	migration.Migrate.SetVersion("%s", migrate_%s)
}

// migrate_%s 迁移函数
func migrate_%s(db *gorm.DB, version string) error {
	// TODO: 在这里实现具体的迁移逻辑
	// 例如：
	// err := db.AutoMigrate(&models.YourModel{})
	// if err != nil {
	//     return err
	// }

	// 记录迁移完成
	return db.Create(&base.Migration{
		Version: version,
		Name:    "迁移描述", // TODO: 修改为具体的迁移描述
		Status:  "completed",
	}).Error
}
`, packageName, timestamp, timestamp, timestamp, timestamp)
}

// fileExists 检查文件或目录是否存在
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

// GetMigrateVersion 获取迁移工具版本信息
func GetMigrateVersion() string {
	return global.Version
}

// Package migration 提供数据库迁移管理功能
// 用于 rentpro-admin 租赁管理系统的版本化数据库迁移
//
// 核心功能：
// 1. 版本化迁移管理：按时间戳排序执行迁移脚本
// 2. 迁移状态跟踪：记录已执行的迁移，避免重复执行
// 3. 并发安全：使用互斥锁保证迁移注册的线程安全
// 4. 错误处理：提供完整的错误处理和日志记录
//
// 使用方式：
//  1. 通过 init() 函数自动注册迁移：
//     migration.Migrate.SetVersion("版本号", 迁移函数)
//  2. 在迁移命令中调用：
//     migration.Migrate.SetDb(数据库连接)
//     migration.Migrate.Migrate()
//
// 设计模式：
// - 单例模式：全局唯一的 Migrate 实例
// - 注册器模式：通过 SetVersion 注册迁移函数
// - 策略模式：每个版本对应不同的迁移策略
package migration

import (
	"log"
	"path/filepath"
	"sort"
	"sync"

	"gorm.io/gorm"
)

// Migrate 全局迁移管理器实例
// 这是一个全局单例，用于管理所有的数据库迁移操作
// 初始化时创建了一个空的 version map，用于存储所有注册的迁移函数
// 通过各个包的 init() 函数调用 SetVersion() 方法进行迁移注册
var Migrate = &Migration{
	version: make(map[string]func(db *gorm.DB, version string) error),
}

// Migration 数据库迁移管理器结构体
// 负责管理数据库迁移的整个生命周期，包括注册、执行和状态管理
type Migration struct {
	// db GORM 数据库连接实例
	// 用于执行迁移操作和查询迁移状态
	db *gorm.DB

	// version 版本号到迁移函数的映射
	// key: 版本号（通常为时间戳格式，如 "1700000000001"）
	// value: 迁移执行函数，接收数据库连接和版本号参数
	version map[string]func(db *gorm.DB, version string) error

	// mutex 互斥锁，用于保证并发安全
	// 在注册迁移函数时防止竞态条件，确保 map 操作的原子性
	mutex sync.Mutex
}

// GetDb 获取当前的数据库连接实例
// 返回值: *gorm.DB - GORM 数据库连接实例，如果未设置则返回 nil
// 使用场景：在迁移函数中获取数据库连接进行操作
func (e *Migration) GetDb() *gorm.DB {
	return e.db
}

// SetDb 设置数据库连接实例
// 参数: db *gorm.DB - GORM 数据库连接实例
// 使用场景：在执行迁移之前设置数据库连接
// 注意：必须在调用 Migrate() 方法之前设置数据库连接
func (e *Migration) SetDb(db *gorm.DB) {
	e.db = db
}

// SetVersion 注册一个版本化迁移函数
// 这个方法是线程安全的，使用互斥锁保证并发安全
//
// 参数:
//
//	k string - 版本号，通常使用时间戳格式（如 "1700000000001"）
//	         版本号必须唯一，并且能按字典序排列
//	f func(db *gorm.DB, version string) error - 迁移执行函数
//	         接收数据库连接和版本号，返回错误信息
//
// 使用场景：
//
//	在各个迁移包的 init() 函数中调用，自动注册迁移
//
// 示例:
//
//	migration.Migrate.SetVersion("1700000000001", createUsersTable)
func (e *Migration) SetVersion(k string, f func(db *gorm.DB, version string) error) {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	e.version[k] = f
}

// Migrate 执行所有已注册的数据库迁移
// 这是整个迁移系统的核心方法，负责管理和执行所有迁移操作
//
// 执行流程：
// 1. 收集所有已注册的版本号
// 2. 按版本号进行字典序排列（确保按时间顺序执行）
// 3. 逐个检查每个版本是否已经执行过
// 4. 对于未执行的版本，调用对应的迁移函数
// 5. 在执行过程中记录日志和处理错误
//
// 数据库表依赖：
//
//	依赖 sys_migration 表来记录已执行的迁移
//	表结构包含 version 字段用于检查迁移状态
//
// 错误处理：
//
//	如果任何迁移执行失败，会调用 log.Fatalln 终止程序
//	这确保了数据库状态的一致性和安全性
//
// 注意事项：
//
//	调用前必须先调用 SetDb() 设置数据库连接
func (e *Migration) Migrate() {
	// 收集所有已注册的版本号
	versions := make([]string, 0)
	for k := range e.version {
		versions = append(versions, k)
	}

	// 确保版本号按字典序排列（时间戳格式保证了正确的时间顺序）
	if !sort.StringsAreSorted(versions) {
		sort.Strings(versions)
	}

	var err error
	var count int64

	// 逐个执行迁移
	for _, v := range versions {
		// 检查该版本是否已经执行过
		err = e.db.Table("sys_migration").Where("version = ?", v).Count(&count).Error
		if err != nil {
			log.Fatalln(err)
		}

		// 如果已经执行过，跳过该版本
		if count > 0 {
			log.Println(count)
			count = 0
			continue
		}

		// 执行迁移函数
		err = (e.version[v])(e.db.Debug(), v)
		if err != nil {
			log.Fatalln(err)
		}
	}
}

// GetFilename 从文件路径中提取时间戳版本号
// 这是一个工具函数，用于从迁移文件路径中解析出版本号
//
// 参数:
//
//	s string - 文件路径，如 "/path/to/1700000000001_create_table.go"
//
// 返回值:
//
//	string - 提取的版本号，取文件名的前13个字符
//	        如 "1700000000001" （时间戳格式）
//
// 使用场景:
//
//	在生成迁移文件时根据文件名确定版本号
//	或者在动态扫描迁移文件时解析版本信息
//
// 注意:
//
//	函数假设文件名格式为: {13位时间戳}_{description}.go
//	如果文件名不符合这个格式，可能会导致错误的版本号
func GetFilename(s string) string {
	s = filepath.Base(s)
	return s[:13]
}

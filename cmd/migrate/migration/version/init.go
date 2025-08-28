package version
// Package version 包含 go-admin 框架的标准迁移文件
// 这些迁移文件通常来自 go-admin 框架本身，用于创建系统基础表结构
package version

import (
	"rentPro/rentpro-admin/cmd/migrate/migration"
)

func init() {
	// 这里会自动注册 go-admin 框架的迁移文件
	// 具体的迁移文件会通过 _ import 的方式自动加载
	
	// 示例：注册一个基础的系统表迁移
	// migration.Migrate.SetVersion("1599190683659", createSystemTables)
	
	// TODO: 添加 go-admin 框架的标准迁移文件
}
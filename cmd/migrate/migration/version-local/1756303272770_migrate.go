package version_local

import (
	"rentPro/rentpro-admin/cmd/migrate/migration"
	"rentPro/rentpro-admin/common/models/base"
	"rentPro/rentpro-admin/common/models/system"

	"gorm.io/gorm"
)

func init() {
	migration.Migrate.SetVersion("1756303272770", migrate_1756303272770)
}

// migrate_1756303272770 权限系统数据库表创建和初始化
func migrate_1756303272770(db *gorm.DB, version string) error {
	// 创建权限系统相关表
	err := db.AutoMigrate(
		&system.SysUser{},
		&system.SysRole{},
		&system.SysMenu{},
		&system.SysDept{},
		&system.SysPost{},
		&base.Migration{},
	)
	if err != nil {
		return err
	}

	// 初始化权限系统数据
	err = base.InitDefaultData(db)
	if err != nil {
		return err
	}

	// 记录迁移完成
	return db.Create(&base.Migration{
		Version: version,
		Name:    "创建权限系统表并初始化数据",
		Status:  "completed",
	}).Error
}

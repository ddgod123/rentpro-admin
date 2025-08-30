package version_local

import (
	"rentPro/rentpro-admin/cmd/migrate/migration"
	commonModels "rentPro/rentpro-admin/common/models/base"

	"gorm.io/gorm"
)

func init() {
	migration.Migrate.SetVersion("1756303700000", initAuthDataFromSQL)
}

// initAuthDataFromSQL 使用新的SQL文件方式初始化权限数据
func initAuthDataFromSQL(db *gorm.DB, version string) error {
	// 使用新的SQL文件加载器初始化数据
	if err := commonModels.InitDefaultData(db); err != nil {
		return err
	}

	// 记录迁移完成
	return db.Create(&commonModels.Migration{
		Version: version,
		Name:    "使用SQL文件方式初始化权限系统数据",
		Status:  "completed",
	}).Error
}

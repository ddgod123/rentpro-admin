package version_local

import (
	"rentPro/rentpro-admin/cmd/migrate/migration"
	"rentPro/rentpro-admin/cmd/migrate/migration/models"
	commonModels "rentPro/rentpro-admin/common/models"

	"gorm.io/gorm"
)

func init() {
	migration.Migrate.SetVersion("1700000000001", createRentalTables)
}

func createRentalTables(db *gorm.DB, version string) error {
	err := db.AutoMigrate(&models.SysBuildings{})
	if err != nil {
		return err
	}

	// 记录迁移完成
	return db.Create(&commonModels.Migration{
		Version: version,
		Name:    "创建租赁管理系统基础表结构",
		Status:  "completed",
	}).Error
}

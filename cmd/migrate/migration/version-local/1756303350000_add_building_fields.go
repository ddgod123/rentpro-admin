package version_local

import (
	"rentPro/rentpro-admin/cmd/migrate/migration"
	commonModels "rentPro/rentpro-admin/common/models"

	"gorm.io/gorm"
)

func init() {
	migration.Migrate.SetVersion("1756303350000", addBuildingIdAndTimestamps)
}

// addBuildingIdAndTimestamps 为sys_buildings表添加ID和时间戳字段
func addBuildingIdAndTimestamps(db *gorm.DB, version string) error {
	// 添加ID字段（主键，自增）
	if err := db.Exec("ALTER TABLE sys_buildings ADD COLUMN id bigint(20) NOT NULL AUTO_INCREMENT PRIMARY KEY FIRST").Error; err != nil {
		return err
	}

	// 添加时间戳字段
	if err := db.Exec("ALTER TABLE sys_buildings ADD COLUMN created_at datetime DEFAULT NULL").Error; err != nil {
		return err
	}

	if err := db.Exec("ALTER TABLE sys_buildings ADD COLUMN updated_at datetime DEFAULT NULL").Error; err != nil {
		return err
	}

	if err := db.Exec("ALTER TABLE sys_buildings ADD COLUMN deleted_at datetime DEFAULT NULL").Error; err != nil {
		return err
	}

	// 为deleted_at字段添加索引
	if err := db.Exec("ALTER TABLE sys_buildings ADD INDEX idx_deleted_at (deleted_at)").Error; err != nil {
		return err
	}

	// 记录迁移完成
	return db.Create(&commonModels.Migration{
		Version: version,
		Name:    "为sys_buildings表添加ID和时间戳字段",
		Status:  "completed",
	}).Error
}

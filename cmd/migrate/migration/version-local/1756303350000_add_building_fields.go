package version_local

import (
	"rentPro/rentpro-admin/cmd/migrate/migration"
	commonModels "rentPro/rentpro-admin/common/models/base"

	"gorm.io/gorm"
)

func init() {
	migration.Migrate.SetVersion("1756303350000", addBuildingIdAndTimestamps)
}

// addBuildingIdAndTimestamps 为sys_buildings表添加时间戳字段
func addBuildingIdAndTimestamps(db *gorm.DB, version string) error {
	// 检查并添加时间戳字段（如果不存在）
	// 添加created_at字段
	var count int64
	db.Raw("SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'sys_buildings' AND COLUMN_NAME = 'created_at'").Scan(&count)
	if count == 0 {
		if err := db.Exec("ALTER TABLE sys_buildings ADD COLUMN created_at datetime(3) DEFAULT NULL").Error; err != nil {
			return err
		}
	}

	// 添加updated_at字段
	db.Raw("SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'sys_buildings' AND COLUMN_NAME = 'updated_at'").Scan(&count)
	if count == 0 {
		if err := db.Exec("ALTER TABLE sys_buildings ADD COLUMN updated_at datetime(3) DEFAULT NULL").Error; err != nil {
			return err
		}
	}

	// 添加deleted_at字段
	db.Raw("SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'sys_buildings' AND COLUMN_NAME = 'deleted_at'").Scan(&count)
	if count == 0 {
		if err := db.Exec("ALTER TABLE sys_buildings ADD COLUMN deleted_at datetime(3) DEFAULT NULL").Error; err != nil {
			return err
		}
	}

	// 检查并添加deleted_at索引（如果不存在）
	db.Raw("SELECT COUNT(*) FROM INFORMATION_SCHEMA.STATISTICS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'sys_buildings' AND INDEX_NAME = 'idx_sys_buildings_deleted_at'").Scan(&count)
	if count == 0 {
		if err := db.Exec("ALTER TABLE sys_buildings ADD INDEX idx_sys_buildings_deleted_at (deleted_at)").Error; err != nil {
			return err
		}
	}

	// 记录迁移完成
	return db.Create(&commonModels.Migration{
		Version: version,
		Name:    "为sys_buildings表添加时间戳字段",
		Status:  "completed",
	}).Error
}

package version_local

import (
	"rentPro/rentpro-admin/cmd/migrate/migration"
	"rentPro/rentpro-admin/cmd/migrate/migration/models"
	commonModels "rentPro/rentpro-admin/common/models"

	"gorm.io/gorm"
)

func init() {
	migration.Migrate.SetVersion("1756303600000", createRentalBusinessTables)
}

// createRentalBusinessTables 创建租赁业务相关表结构
func createRentalBusinessTables(db *gorm.DB, version string) error {
	// 暂时禁用外键约束检查
	db.Exec("SET FOREIGN_KEY_CHECKS = 0")
	defer db.Exec("SET FOREIGN_KEY_CHECKS = 1")

	// 按依赖关系顺序创建表
	// 1. 先创建经纪人表（独立表）
	err := db.AutoMigrate(&models.SysAgent{})
	if err != nil {
		return err
	}

	// 2. 创建房源表（依赖楼盘表和经纪人表）
	err = db.AutoMigrate(&models.SysHouse{})
	if err != nil {
		return err
	}

	// 3. 创建租客表（独立表）
	err = db.AutoMigrate(&models.SysTenant{})
	if err != nil {
		return err
	}

	// 4. 创建租赁记录表（依赖房源表、租客表、经纪人表）
	err = db.AutoMigrate(&models.SysRentalRecord{})
	if err != nil {
		return err
	}

	// 手动添加外键约束（如果需要）
	// 由于类型不匹配，我们暂时不添加外键约束
	// 可以在应用层面维护数据一致性

	// 记录迁移完成
	return db.Create(&commonModels.Migration{
		Version: version,
		Name:    "创建租赁业务相关表结构(经纪人/房源/租客/租赁记录)",
		Status:  "completed",
	}).Error
}

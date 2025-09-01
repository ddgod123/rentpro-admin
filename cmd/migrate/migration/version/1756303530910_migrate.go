package version

import (
	"rentPro/rentpro-admin/cmd/migrate/migration"
	"rentPro/rentpro-admin/common/models/base"
	"rentPro/rentpro-admin/common/models/rental"
	"rentPro/rentpro-admin/common/models/system"

	"gorm.io/gorm"
)

func init() {
	migration.Migrate.SetVersion("1756303530910", migrate_1756303530910)
}

// migrate_1756303530910 迁移函数
// 创建所有系统表和业务表
func migrate_1756303530910(db *gorm.DB, version string) error {
	// 系统表迁移
	systemModels := []interface{}{
		&system.SysUser{},
		&system.SysRole{},
		&system.SysMenu{},
		&system.SysDept{},
		&system.SysPost{},
	}

	// 租赁业务表迁移
	rentalModels := []interface{}{
		&rental.SysBuildings{},
		&rental.SysHouseType{},
		&rental.District{},
		&rental.BusinessArea{},
	}

	// 基础表迁移
	baseModels := []interface{}{
		&base.Migration{},
	}

	// 合并所有模型
	allModels := append(systemModels, rentalModels...)
	allModels = append(allModels, baseModels...)

	// 执行自动迁移
	for _, model := range allModels {
		if err := db.AutoMigrate(model); err != nil {
			return err
		}
	}

	// 记录迁移完成
	return db.Create(&base.Migration{
		Version: version,
		Name:    "创建所有系统表和业务表",
		Status:  "completed",
	}).Error
}

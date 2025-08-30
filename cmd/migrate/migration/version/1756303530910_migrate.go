package version

import (
	"rentPro/rentpro-admin/cmd/migrate/migration"
	"rentPro/rentpro-admin/common/models/base"
	"gorm.io/gorm"
)

func init() {
	migration.Migrate.SetVersion("1756303530910", migrate_1756303530910)
}

// migrate_1756303530910 迁移函数
func migrate_1756303530910(db *gorm.DB, version string) error {
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

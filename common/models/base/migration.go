package base

import (
	"time"
)

// Migration 数据库迁移记录模型
// 用于跟踪已执行的数据库迁移版本
type Migration struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	Version   string    `gorm:"size:191;not null;unique;index:idx_version" json:"version" comment:"迁移版本号"`
	Name      string    `gorm:"size:255" json:"name" comment:"迁移名称"`
	Status    string    `gorm:"size:20;default:'completed'" json:"status" comment:"迁移状态(pending,running,completed,failed)"`
	CreatedAt time.Time `json:"created_at" comment:"创建时间"`
	UpdatedAt time.Time `json:"updated_at" comment:"更新时间"`
}

// TableName 设置表名
func (Migration) TableName() string {
	return "sys_migration"
}

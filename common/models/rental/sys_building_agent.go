package rental

import (
	"time"
)

// SysBuildingAgent 楼盘房源人关联表模型
type SysBuildingAgent struct {
	// 主键
	ID uint `json:"id" gorm:"primaryKey;autoIncrement" comment:"主键ID"`

	// 关联字段
	BuildingID uint `json:"building_id" gorm:"not null;index:idx_building_id" comment:"楼盘ID"`
	AgentID    uint `json:"agent_id" gorm:"not null;index:idx_agent_id" comment:"房源人ID"`

	// 关联信息
	Role        string `json:"role" gorm:"size:50;default:'agent'" comment:"角色(agent:普通房源人, primary:主要负责人, backup:备用联系人)"`
	IsActive    bool   `json:"is_active" gorm:"default:true;index:idx_is_active" comment:"是否激活"`
	Priority    int    `json:"priority" gorm:"default:1" comment:"优先级(1-10, 数字越小优先级越高)"`
	Description string `json:"description" gorm:"size:200" comment:"备注说明"`

	// 管理信息
	CreatedBy string `json:"created_by" gorm:"size:50" comment:"创建人"`
	UpdatedBy string `json:"updated_by" gorm:"size:50" comment:"更新人"`

	// 时间戳
	CreatedAt *time.Time `json:"created_at" gorm:"autoCreateTime" comment:"创建时间"`
	UpdatedAt *time.Time `json:"updated_at" gorm:"autoUpdateTime" comment:"更新时间"`
	DeletedAt *time.Time `json:"deleted_at" gorm:"index" comment:"删除时间"`

	// 关联对象(用于GORM关联查询)
	Building *SysBuildings `json:"building,omitempty" gorm:"foreignKey:BuildingID" comment:"关联楼盘"`
	Agent    *SysAgent     `json:"agent,omitempty" gorm:"foreignKey:AgentID" comment:"关联房源人"`
}

// TableName 设置表名
func (SysBuildingAgent) TableName() string {
	return "sys_building_agents"
}

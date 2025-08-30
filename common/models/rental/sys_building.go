package rental

import (
	"time"
)

// SysBuildings 楼盘模型
type SysBuildings struct {
	// 主键
	ID uint `json:"id" gorm:"primaryKey;autoIncrement" comment:"主键ID"`

	// 基础信息
	Name            string `json:"name" gorm:"size:100;not null;index:idx_name" comment:"楼盘名称"`
	Developer       string `json:"developer" gorm:"size:100" comment:"开发商"`
	DetailedAddress string `json:"detailed_address" gorm:"size:500;not null" comment:"详细地址"`
	City            string `json:"city" gorm:"size:50;not null" comment:"城市"`
	District        string `json:"district" gorm:"size:50;not null" comment:"区域/区"`
	BusinessArea    string `json:"business_area" gorm:"size:100" comment:"所属商圈"`
	PropertyType    string `json:"property_type" gorm:"size:50" comment:"物业类型(住宅/商业/办公等)"`

	PropertyCompany string `json:"property_company" gorm:"size:100" comment:"物业公司"`
	Description     string `json:"description" gorm:"type:text" comment:"楼盘描述"`

	// 统计信息
	SaleCount      int `json:"sale_count" gorm:"default:0;index:idx_sale_count" comment:"在售数"`
	RentCount      int `json:"rent_count" gorm:"default:0;index:idx_rent_count" comment:"在租数"`
	SaleDealsCount int `json:"sale_deals_count" gorm:"default:0" comment:"在售成交数"`
	RentDealsCount int `json:"rent_deals_count" gorm:"default:0" comment:"在租成交数"`

	// 状态信息
	Status string `json:"status" gorm:"size:20;not null;default:'active';index:idx_status" comment:"状态(active:在租/售, inactive:停用, pending:审核中)"`
	IsHot  bool   `json:"is_hot" gorm:"default:false;index:idx_is_hot" comment:"是否顶豪楼盘"`

	// 管理信息
	CreatedBy string `json:"created_by" gorm:"size:50" comment:"创建人"`
	UpdatedBy string `json:"updated_by" gorm:"size:50" comment:"更新人"`

	// 时间戳
	CreatedAt *time.Time `json:"created_at" gorm:"autoCreateTime" comment:"创建时间"`
	UpdatedAt *time.Time `json:"updated_at" gorm:"autoUpdateTime" comment:"更新时间"`
	DeletedAt *time.Time `json:"deleted_at" gorm:"index" comment:"删除时间"`

	// 关联对象(用于GORM关联查询)
	Agents []SysAgent `json:"agents,omitempty" gorm:"many2many:sys_building_agents;" comment:"房源人列表"`
}

// TableName 设置表名
func (SysBuildings) TableName() string {
	return "sys_buildings"
}

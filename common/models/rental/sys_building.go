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
	DetailedAddress string `json:"detailedAddress" gorm:"size:500;not null;column:detailed_address" comment:"详细地址"`
	City            string `json:"city" gorm:"size:50;not null" comment:"城市"`
	District        string `json:"district" gorm:"size:50;not null" comment:"区域/区"`
	BusinessArea    string `json:"businessArea" gorm:"size:100" comment:"所属商圈"`
	SubDistrict     string `json:"subDistrict" gorm:"size:50" comment:"街道"`
	PropertyType    string `json:"propertyType" gorm:"size:50" comment:"物业类型(住宅/商业/办公等)"`

	PropertyCompany string `json:"propertyCompany" gorm:"size:100" comment:"物业公司"`
	Description     string `json:"description" gorm:"type:text" comment:"楼盘描述"`

	// 统计信息
	SaleCount      int `json:"saleCount" gorm:"default:0;index:idx_sale_count" comment:"在售数"`
	RentCount      int `json:"rentCount" gorm:"default:0;index:idx_rent_count" comment:"在租数"`
	SaleDealsCount int `json:"saleDealsCount" gorm:"default:0" comment:"在售成交数"`
	RentDealsCount int `json:"rentDealsCount" gorm:"default:0" comment:"在租成交数"`

	// 状态信息
	Status string `json:"status" gorm:"size:20;not null;default:'active';index:idx_status" comment:"状态(active:在租/售, inactive:停用, pending:审核中)"`
	IsHot  bool   `json:"isHot" gorm:"default:false;index:idx_is_hot" comment:"是否顶豪楼盘"`

	// 管理信息
	CreatedBy string `json:"createdBy" gorm:"size:50" comment:"创建人"`
	UpdatedBy string `json:"updatedBy" gorm:"size:50" comment:"更新人"`

	// 时间戳
	CreatedAt *time.Time `json:"createdAt" gorm:"autoCreateTime" comment:"创建时间"`
	UpdatedAt *time.Time `json:"updatedAt" gorm:"autoUpdateTime" comment:"更新时间"`
	DeletedAt *time.Time `json:"deletedAt" gorm:"index" comment:"删除时间"`

	
}

// TableName 设置表名
func (SysBuildings) TableName() string {
	return "sys_buildings"
}

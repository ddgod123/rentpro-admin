package rental

/*
�� 下一步开发建议
现在你可以：
完善户型API: 创建户型的增删改查接口
更新前端页面: 在户型展示页面中使用真实数据
扩展功能: 添加户型筛选、排序、搜索等功能
图片管理: 实现户型图片的上传和管理
库存管理: 实现库存的实时更新和预订功能
户型数据模型已经完全创建完成，包含了完整的字段定义、数据库表结构、示例数据和迁移文件！🎊

*/

import (
	"time"
)

// SysHouseType 户型模型
type SysHouseType struct {
	// 主键
	ID uint `json:"id" gorm:"primaryKey;autoIncrement" comment:"主键ID"`

	// 基础信息
	Name        string `json:"name" gorm:"size:100;not null;index:idx_name" comment:"户型名称"`
	Code        string `json:"code" gorm:"size:50;not null;uniqueIndex:idx_code" comment:"户型编码"`
	Description string `json:"description" gorm:"type:text" comment:"户型描述"`

	// 楼盘关联
	BuildingID uint         `json:"building_id" gorm:"not null;index:idx_building_id" comment:"所属楼盘ID"`
	Building   SysBuildings `json:"building,omitempty" gorm:"foreignKey:BuildingID" comment:"所属楼盘"`

	// 户型规格
	Area        float64 `json:"area" gorm:"type:decimal(8,2);not null;index:idx_area" comment:"建筑面积(平方米)"`
	Rooms       int     `json:"rooms" gorm:"not null;default:1" comment:"房间数"`
	Halls       int     `json:"halls" gorm:"not null;default:1" comment:"客厅数"`
	Bathrooms   int     `json:"bathrooms" gorm:"not null;default:1" comment:"卫生间数"`
	Balconies   int     `json:"balconies" gorm:"default:0" comment:"阳台数"`
	FloorHeight float64 `json:"floor_height" gorm:"type:decimal(4,2)" comment:"层高(米)"`

	// 朝向信息
	Orientation string `json:"orientation" gorm:"size:50" comment:"朝向(南北/东西/南向/北向等)"`
	View        string `json:"view" gorm:"size:100" comment:"景观(海景/山景/城市景观等)"`

	// 价格信息
	SalePrice    float64 `json:"sale_price" gorm:"type:decimal(12,2);default:0;index:idx_area" comment:"售价(元)"`
	RentPrice    float64 `json:"rent_price" gorm:"type:decimal(8,2);default:0;index:idx_rent_price" comment:"月租金(元)"`
	SalePricePer float64 `json:"sale_price_per" gorm:"type:decimal(8,2);default:0" comment:"单价(元/平方米)"`
	RentPricePer float64 `json:"rent_price_per" gorm:"type:decimal(6,2);default:0" comment:"租金单价(元/平方米/月)"`

	// 库存信息
	TotalStock    int `json:"total_stock" gorm:"not null;default:0" comment:"总库存"`
	SaleStock     int `json:"sale_stock" gorm:"not null;default:0" comment:"在售库存"`
	RentStock     int `json:"rent_stock" gorm:"not null;default:0" comment:"在租库存"`
	ReservedStock int `json:"reserved_stock" gorm:"not null;default:0" comment:"已预订库存"`

	// 状态信息
	Status     string `json:"status" gorm:"size:20;not null;default:'active';index:idx_status" comment:"状态(active:在售/租, inactive:停用, pending:审核中)"`
	SaleStatus string `json:"sale_status" gorm:"size:20;default:'available'" comment:"销售状态(available:可售, sold:已售, reserved:已预订)"`
	RentStatus string `json:"rent_status" gorm:"size:20;default:'available'" comment:"租赁状态(available:可租, rented:已租, reserved:已预订)"`
	IsHot      bool   `json:"is_hot" gorm:"default:false;index:idx_is_hot" comment:"是否热门户型"`

	// 图片信息
	MainImage string   `json:"main_image" gorm:"size:500" comment:"主图URL"`
	ImageUrls []string `json:"image_urls" gorm:"type:json" comment:"图片URL列表"`

	// 特色标签
	Tags []string `json:"tags" gorm:"type:json" comment:"特色标签(南北通透/精装修/地铁房等)"`

	// 管理信息
	CreatedBy string `json:"created_by" gorm:"size:50" comment:"创建人"`
	UpdatedBy string `json:"updated_by" gorm:"size:50" comment:"更新人"`

	// 时间戳
	CreatedAt *time.Time `json:"created_at" gorm:"autoCreateTime" comment:"创建时间"`
	UpdatedAt *time.Time `json:"updated_at" gorm:"autoUpdateTime" comment:"更新时间"`
	DeletedAt *time.Time `json:"deleted_at" gorm:"index" comment:"删除时间"`
}

// TableName 设置表名
func (SysHouseType) TableName() string {
	return "sys_house_types"
}

// GetStockStatus 获取库存状态描述
func (h *SysHouseType) GetStockStatus() string {
	if h.SaleStock > 0 && h.RentStock > 0 {
		return "可售可租"
	} else if h.SaleStock > 0 {
		return "仅可售"
	} else if h.RentPrice > 0 {
		return "仅可租"
	} else {
		return "无库存"
	}
}

// GetPriceRange 获取价格区间描述
func (h *SysHouseType) GetPriceRange() string {
	if h.SalePrice > 0 && h.RentPrice > 0 {
		return "可售可租"
	} else if h.SalePrice > 0 {
		return "仅可售"
	} else if h.RentPrice > 0 {
		return "仅可租"
	} else {
		return "价格面议"
	}
}

// IsAvailable 检查是否可售或可租
func (h *SysHouseType) IsAvailable() bool {
	return h.Status == "active" && (h.SaleStock > 0 || h.RentStock > 0)
}

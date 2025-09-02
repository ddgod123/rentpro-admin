package rental

import (
	"time"
)

// SysHouse 房屋模型
type SysHouse struct {
	// 主键
	ID uint `json:"id" gorm:"primaryKey;autoIncrement" comment:"主键ID"`

	// 基础信息
	Name string `json:"name" gorm:"size:100;not null;index:idx_name" comment:"房屋名称"`
	Code string `json:"code" gorm:"size:50;not null;uniqueIndex:idx_code" comment:"房屋编码"`

	// 楼盘关联
	BuildingID uint         `json:"buildingId" gorm:"not null;index:idx_building_id" comment:"所属楼盘ID"`
	Building   SysBuildings `json:"building,omitempty" gorm:"foreignKey:BuildingID" comment:"所属楼盘"`

	// 房屋位置
	Floor      int    `json:"floor" comment:"楼层"`
	Unit       string `json:"unit" gorm:"size:20" comment:"单元号"`
	RoomNumber string `json:"roomNumber" gorm:"size:20" comment:"房号"`

	// 房屋规格
	Area       float64 `json:"area" gorm:"type:decimal(8,2);not null;index:idx_area" comment:"建筑面积(平方米)"`
	UsableArea float64 `json:"usableArea" gorm:"type:decimal(8,2)" comment:"使用面积(平方米)"`
	Rooms      int     `json:"rooms" gorm:"not null;default:1" comment:"房间数"`
	Halls      int     `json:"halls" gorm:"not null;default:1" comment:"客厅数"`
	Bathrooms  int     `json:"bathrooms" gorm:"not null;default:1" comment:"卫生间数"`
	Balconies  int     `json:"balconies" gorm:"default:0" comment:"阳台数"`

	// 朝向信息
	Orientation string `json:"orientation" gorm:"size:50" comment:"朝向(南北/东西/南向/北向等)"`
	View        string `json:"view" gorm:"size:100" comment:"景观(海景/山景/城市景观等)"`

	// 装修信息
	Decoration string `json:"decoration" gorm:"size:50" comment:"装修情况(毛坯/简装/精装/豪装)"`

	// 价格信息
	SalePrice    float64 `json:"salePrice" gorm:"type:decimal(12,2);default:0;index:idx_sale_price" comment:"售价(元)"`
	RentPrice    float64 `json:"rentPrice" gorm:"type:decimal(8,2);default:0;index:idx_rent_price" comment:"月租金(元)"`
	SalePricePer float64 `json:"salePricePer" gorm:"type:decimal(8,2);default:0" comment:"单价(元/平方米)"`
	RentPricePer float64 `json:"rentPricePer" gorm:"type:decimal(6,2);default:0" comment:"租金单价(元/平方米/月)"`

	// 状态信息
	Status     string `json:"status" gorm:"size:20;not null;default:'available';index:idx_status" comment:"状态(available:可租/售, rented:已租, sold:已售, maintenance:维护中, inactive:停用)"`
	SaleStatus string `json:"saleStatus" gorm:"size:20;default:'available'" comment:"销售状态(available:可售, sold:已售, reserved:已预订)"`
	RentStatus string `json:"rentStatus" gorm:"size:20;default:'available'" comment:"租赁状态(available:可租, rented:已租, reserved:已预订)"`

	// 图片信息
	MainImage string   `json:"mainImage" gorm:"size:500" comment:"主图URL"`
	ImageUrls []string `json:"imageUrls" gorm:"type:json" comment:"图片URL列表"`

	// 特色标签
	Tags []string `json:"tags" gorm:"type:json" comment:"特色标签(南北通透/精装修/地铁房等)"`

	// 配套设施
	Facilities []string `json:"facilities" gorm:"type:json" comment:"配套设施(空调/暖气/家具/家电等)"`

	// 备注
	Description string `json:"description" gorm:"type:text" comment:"房屋描述"`
	Notes       string `json:"notes" gorm:"type:text" comment:"备注信息"`

	// 管理信息
	CreatedBy string `json:"createdBy" gorm:"size:50" comment:"创建人"`
	UpdatedBy string `json:"updatedBy" gorm:"size:50" comment:"更新人"`

	// 时间戳
	CreatedAt *time.Time `json:"createdAt" gorm:"autoCreateTime" comment:"创建时间"`
	UpdatedAt *time.Time `json:"updatedAt" gorm:"autoUpdateTime" comment:"更新时间"`
	DeletedAt *time.Time `json:"deletedAt" gorm:"index" comment:"删除时间"`
}

// TableName 设置表名
func (SysHouse) TableName() string {
	return "sys_houses"
}

// GetStatusText 获取状态文本描述
func (h *SysHouse) GetStatusText() string {
	switch h.Status {
	case "available":
		return "可租/售"
	case "rented":
		return "已租"
	case "sold":
		return "已售"
	case "maintenance":
		return "维护中"
	case "inactive":
		return "停用"
	default:
		return "未知"
	}
}

// GetSaleStatusText 获取销售状态文本描述
func (h *SysHouse) GetSaleStatusText() string {
	switch h.SaleStatus {
	case "available":
		return "可售"
	case "sold":
		return "已售"
	case "reserved":
		return "已预订"
	default:
		return "未知"
	}
}

// GetRentStatusText 获取租赁状态文本描述
func (h *SysHouse) GetRentStatusText() string {
	switch h.RentStatus {
	case "available":
		return "可租"
	case "rented":
		return "已租"
	case "reserved":
		return "已预订"
	default:
		return "未知"
	}
}

// GetDecorationText 获取装修情况文本描述
func (h *SysHouse) GetDecorationText() string {
	switch h.Decoration {
	case "bare":
		return "毛坯"
	case "simple":
		return "简装"
	case "fine":
		return "精装"
	case "luxury":
		return "豪装"
	default:
		return h.Decoration
	}
}

// IsAvailableForSale 判断是否可售
func (h *SysHouse) IsAvailableForSale() bool {
	return h.Status == "available" && h.SaleStatus == "available"
}

// IsAvailableForRent 判断是否可租
func (h *SysHouse) IsAvailableForRent() bool {
	return h.Status == "available" && h.RentStatus == "available"
}

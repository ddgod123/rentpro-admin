package rental

import (
	"fmt"
	"time"
)

// SysHouse 房屋模型 - 具体房屋实例
type SysHouse struct {
	// 主键
	ID uint `json:"id" gorm:"primaryKey;autoIncrement" comment:"主键ID"`

	// 基础信息
	Name string `json:"name" gorm:"size:100;not null;index:idx_name" comment:"房屋名称"`
	Code string `json:"code" gorm:"size:50;not null;uniqueIndex:idx_code" comment:"房屋编码"`

	// 关联关系
	BuildingID  uint         `json:"buildingId" gorm:"not null;index:idx_building_id" comment:"所属楼盘ID"`
	Building    SysBuildings `json:"building,omitempty" gorm:"foreignKey:BuildingID" comment:"所属楼盘"`
	HouseTypeID uint         `json:"houseTypeId" gorm:"not null;index:idx_house_type_id" comment:"所属户型ID"`
	HouseType   SysHouseType `json:"houseType,omitempty" gorm:"foreignKey:HouseTypeID" comment:"所属户型"`

	// 房屋位置
	Floor      int    `json:"floor" comment:"楼层"`
	Unit       string `json:"unit" gorm:"size:20" comment:"单元号"`
	RoomNumber string `json:"roomNumber" gorm:"size:20" comment:"房号"`

	// 实际规格信息（可能与户型标准规格有差异）
	ActualArea       float64 `json:"actualArea" gorm:"type:decimal(8,2)" comment:"实际建筑面积(平方米)"`
	ActualUsableArea float64 `json:"actualUsableArea" gorm:"type:decimal(8,2)" comment:"实际使用面积(平方米)"`

	// 实际朝向和景观（可能与户型标准不同）
	ActualOrientation string `json:"actualOrientation" gorm:"size:50" comment:"实际朝向"`
	ActualView        string `json:"actualView" gorm:"size:100" comment:"实际景观"`

	// 装修信息
	Decoration string `json:"decoration" gorm:"size:50" comment:"装修情况(毛坯/简装/精装/豪装)"`

	// 价格信息（基于户型基准价格的实际价格）
	ActualSalePrice       float64 `json:"actualSalePrice" gorm:"type:decimal(12,2);default:0;index:idx_actual_sale_price" comment:"实际售价(元)"`
	ActualRentPrice       float64 `json:"actualRentPrice" gorm:"type:decimal(8,2);default:0;index:idx_actual_rent_price" comment:"实际月租金(元)"`
	PriceAdjustment       float64 `json:"priceAdjustment" gorm:"type:decimal(8,2);default:0" comment:"价格调整金额(元)"`
	PriceAdjustmentReason string  `json:"priceAdjustmentReason" gorm:"size:200" comment:"价格调整原因"`

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

// GetFullAddress 获取房屋完整地址
func (h *SysHouse) GetFullAddress() string {
	if h.Unit != "" && h.RoomNumber != "" {
		return fmt.Sprintf("%s单元%s室", h.Unit, h.RoomNumber)
	} else if h.RoomNumber != "" {
		return h.RoomNumber
	}
	return ""
}

// GetEffectiveArea 获取有效面积（优先使用实际面积）
func (h *SysHouse) GetEffectiveArea() float64 {
	if h.ActualArea > 0 {
		return h.ActualArea
	}
	// 如果没有实际面积，从户型获取标准面积
	if h.HouseType.StandardArea > 0 {
		return h.HouseType.StandardArea
	}
	return 0
}

// GetEffectiveSalePrice 获取有效售价
func (h *SysHouse) GetEffectiveSalePrice() float64 {
	if h.ActualSalePrice > 0 {
		return h.ActualSalePrice
	}
	// 如果没有设置实际价格，使用户型基准价格
	return h.HouseType.BaseSalePrice + h.PriceAdjustment
}

// GetEffectiveRentPrice 获取有效租金
func (h *SysHouse) GetEffectiveRentPrice() float64 {
	if h.ActualRentPrice > 0 {
		return h.ActualRentPrice
	}
	// 如果没有设置实际价格，使用户型基准价格
	return h.HouseType.BaseRentPrice + h.PriceAdjustment
}

// IsCustomPricing 判断是否使用了自定义定价
func (h *SysHouse) IsCustomPricing() bool {
	return h.ActualSalePrice > 0 || h.ActualRentPrice > 0 || h.PriceAdjustment != 0
}

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
	"fmt"
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

	// 户型标准规格
	StandardArea float64 `json:"standardArea" gorm:"type:decimal(8,2);not null;index:idx_standard_area" comment:"标准建筑面积(平方米)"`
	Rooms        int     `json:"rooms" gorm:"not null;default:1" comment:"房间数"`
	Halls        int     `json:"halls" gorm:"not null;default:1" comment:"客厅数"`
	Bathrooms    int     `json:"bathrooms" gorm:"not null;default:1" comment:"卫生间数"`
	Balconies    int     `json:"balconies" gorm:"default:0" comment:"阳台数"`
	FloorHeight  float64 `json:"floorHeight" gorm:"type:decimal(4,2)" comment:"标准层高(米)"`

	// 标准朝向和景观
	StandardOrientation string `json:"standardOrientation" gorm:"size:50" comment:"标准朝向(南北/东西/南向/北向等)"`
	StandardView        string `json:"standardView" gorm:"size:100" comment:"标准景观(海景/山景/城市景观等)"`

	// 基准价格信息
	BaseSalePrice    float64 `json:"baseSalePrice" gorm:"type:decimal(12,2);default:0;index:idx_base_sale_price" comment:"基准售价(元)"`
	BaseRentPrice    float64 `json:"baseRentPrice" gorm:"type:decimal(8,2);default:0;index:idx_base_rent_price" comment:"基准月租金(元)"`
	BaseSalePricePer float64 `json:"baseSalePricePer" gorm:"type:decimal(8,2);default:0" comment:"基准单价(元/平方米)"`
	BaseRentPricePer float64 `json:"baseRentPricePer" gorm:"type:decimal(6,2);default:0" comment:"基准租金单价(元/平方米/月)"`

	// 库存统计（自动计算）
	TotalStock     int `json:"totalStock" gorm:"default:0" comment:"总库存"`
	AvailableStock int `json:"availableStock" gorm:"default:0" comment:"可用库存"`
	SoldStock      int `json:"soldStock" gorm:"default:0" comment:"已售库存"`
	RentedStock    int `json:"rentedStock" gorm:"default:0" comment:"已租库存"`
	ReservedStock  int `json:"reservedStock" gorm:"default:0" comment:"已预订库存"`

	// 户型状态
	Status string `json:"status" gorm:"size:20;not null;default:'active';index:idx_status" comment:"状态(active:在售/租, inactive:停用, pending:审核中)"`
	IsHot  bool   `json:"isHot" gorm:"default:false;index:idx_is_hot" comment:"是否热门户型"`

	// 户型展示图片
	MainImage    string   `json:"mainImage" gorm:"size:500" comment:"主图URL"`
	FloorPlanUrl string   `json:"floorPlanUrl" gorm:"size:500" comment:"户型图URL"`
	ImageUrls    []string `json:"imageUrls" gorm:"type:json" comment:"图片URL列表"`

	// 特色标签
	Tags []string `json:"tags" gorm:"type:json" comment:"特色标签(南北通透/精装修/地铁房等)"`

	// 管理信息
	CreatedBy string `json:"createdBy" gorm:"size:50" comment:"创建人"`
	UpdatedBy string `json:"updatedBy" gorm:"size:50" comment:"更新人"`

	// 时间戳
	CreatedAt *time.Time `json:"createdAt" gorm:"autoCreateTime" comment:"创建时间"`
	UpdatedAt *time.Time `json:"updatedAt" gorm:"autoUpdateTime" comment:"更新时间"`
	DeletedAt *time.Time `json:"deletedAt" gorm:"index" comment:"删除时间"`

	// 关联的房屋列表
	Houses []SysHouse `json:"houses,omitempty" gorm:"foreignKey:HouseTypeID" comment:"关联房屋列表"`
}

// TableName 设置表名
func (SysHouseType) TableName() string {
	return "sys_house_types"
}

// GetStockStatus 获取库存状态描述
func (h *SysHouseType) GetStockStatus() string {
	if h.AvailableStock > 0 {
		if h.BaseSalePrice > 0 && h.BaseRentPrice > 0 {
			return "可售可租"
		} else if h.BaseSalePrice > 0 {
			return "仅可售"
		} else if h.BaseRentPrice > 0 {
			return "仅可租"
		}
		return "有库存"
	}
	return "无库存"
}

// GetPriceRange 获取基准价格区间描述
func (h *SysHouseType) GetPriceRange() string {
	if h.BaseSalePrice > 0 && h.BaseRentPrice > 0 {
		return fmt.Sprintf("售价: %.0f万, 租金: %.0f元/月", h.BaseSalePrice/10000, h.BaseRentPrice)
	} else if h.BaseSalePrice > 0 {
		return fmt.Sprintf("售价: %.0f万", h.BaseSalePrice/10000)
	} else if h.BaseRentPrice > 0 {
		return fmt.Sprintf("租金: %.0f元/月", h.BaseRentPrice)
	}
	return "价格面议"
}

// IsAvailable 检查是否可售或可租
func (h *SysHouseType) IsAvailable() bool {
	return h.Status == "active" && h.AvailableStock > 0
}

// UpdateStockFromHouses 根据关联房屋更新库存统计
func (h *SysHouseType) UpdateStockFromHouses() {
	if len(h.Houses) == 0 {
		return
	}

	h.TotalStock = len(h.Houses)
	h.AvailableStock = 0
	h.SoldStock = 0
	h.RentedStock = 0
	h.ReservedStock = 0

	for _, house := range h.Houses {
		switch house.Status {
		case "available":
			h.AvailableStock++
		case "sold":
			h.SoldStock++
		case "rented":
			h.RentedStock++
		}

		if house.SaleStatus == "reserved" || house.RentStatus == "reserved" {
			h.ReservedStock++
		}
	}
}

// GetHouseLayout 获取户型布局描述
func (h *SysHouseType) GetHouseLayout() string {
	return fmt.Sprintf("%d室%d厅%d卫", h.Rooms, h.Halls, h.Bathrooms)
}

// CalculateBasePricePer 计算基准单价
func (h *SysHouseType) CalculateBasePricePer() {
	if h.StandardArea > 0 {
		if h.BaseSalePrice > 0 {
			h.BaseSalePricePer = h.BaseSalePrice / h.StandardArea
		}
		if h.BaseRentPrice > 0 {
			h.BaseRentPricePer = h.BaseRentPrice / h.StandardArea
		}
	}
}

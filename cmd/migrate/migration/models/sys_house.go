package models

import (
	"rentPro/rentpro-admin/common/models/rental"
	"time"
)

// SysHouse 房源模型
type SysHouse struct {
	// 主键
	ID uint `json:"id" gorm:"primaryKey;autoIncrement;type:int unsigned" comment:"主键ID"`

	// 关联信息
	BuildingID uint `json:"building_id" gorm:"not null;index:idx_building_id;type:int unsigned" comment:"所属楼盘ID"`
	AgentID    uint `json:"agent_id" gorm:"index:idx_agent_id;type:int unsigned" comment:"负责经纪人ID"`

	// 基础信息
	HouseCode   string `json:"house_code" gorm:"size:50;not null;uniqueIndex:idx_house_code" comment:"房源编号"`
	Floor       int    `json:"floor" gorm:"not null" comment:"楼层"`
	Unit        string `json:"unit" gorm:"size:20" comment:"单元号"`
	Room        string `json:"room" gorm:"size:20;not null" comment:"房间号"`
	FullAddress string `json:"full_address" gorm:"size:500;not null" comment:"完整地址"`

	// 房屋属性
	HouseType   string  `json:"house_type" gorm:"size:50;not null" comment:"房型(1室1厅/2室1厅等)"`
	Area        float64 `json:"area" gorm:"type:decimal(8,2);not null" comment:"建筑面积(平方米)"`
	UsableArea  float64 `json:"usable_area" gorm:"type:decimal(8,2)" comment:"使用面积(平方米)"`
	Orientation string  `json:"orientation" gorm:"size:50" comment:"朝向(南向/北向/东南向等)"`
	Decoration  string  `json:"decoration" gorm:"size:50" comment:"装修情况(毛坯/简装/精装/豪装)"`

	// 租赁信息
	RentPrice    int    `json:"rent_price" gorm:"not null;index:idx_rent_price" comment:"月租金(元)"`
	ServiceFee   int    `json:"service_fee" gorm:"default:0" comment:"服务费(元)"`
	Deposit      int    `json:"deposit" gorm:"not null" comment:"押金(元)"`
	PaymentCycle string `json:"payment_cycle" gorm:"size:20;default:'月付'" comment:"付款周期(月付/季付/年付)"`

	// 配套设施
	HasElevator     bool `json:"has_elevator" gorm:"default:false" comment:"是否有电梯"`
	HasParking      bool `json:"has_parking" gorm:"default:false" comment:"是否有停车位"`
	HasBalcony      bool `json:"has_balcony" gorm:"default:false" comment:"是否有阳台"`
	HasAircon       bool `json:"has_aircon" gorm:"default:false" comment:"是否有空调"`
	HasWifi         bool `json:"has_wifi" gorm:"default:false" comment:"是否有WIFI"`
	HasWasher       bool `json:"has_washer" gorm:"default:false" comment:"是否有洗衣机"`
	HasRefrigerator bool `json:"has_refrigerator" gorm:"default:false" comment:"是否有冰箱"`

	// 状态信息
	RentalStatus  string `json:"rental_status" gorm:"size:20;not null;default:'available';index:idx_rental_status" comment:"出租状态(available:可租/rented:已出租/maintenance:维护中/offline:下架)"`
	IsRecommended bool   `json:"is_recommended" gorm:"default:false;index:idx_is_recommended" comment:"是否推荐房源"`
	ViewCount     int    `json:"view_count" gorm:"default:0" comment:"浏览次数"`

	// 描述信息
	Title       string `json:"title" gorm:"size:200;not null" comment:"房源标题"`
	Description string `json:"description" gorm:"type:text" comment:"房源描述"`
	Images      string `json:"images" gorm:"type:text" comment:"房源图片(JSON格式存储图片URL列表)"`
	Tags        string `json:"tags" gorm:"size:500" comment:"房源标签(JSON格式)"`

	// 管理信息
	CreatedBy string `json:"created_by" gorm:"size:50" comment:"创建人"`
	UpdatedBy string `json:"updated_by" gorm:"size:50" comment:"更新人"`

	// 时间戳
	CreatedAt *time.Time `json:"created_at" gorm:"autoCreateTime" comment:"创建时间"`
	UpdatedAt *time.Time `json:"updated_at" gorm:"autoUpdateTime" comment:"更新时间"`
	DeletedAt *time.Time `json:"deleted_at" gorm:"index" comment:"删除时间"`

	// 关联对象(用于GORM关联查询)
	Building *rental.SysBuildings `json:"building,omitempty" gorm:"-" comment:"所属楼盘"`
	Agent    *rental.SysAgent     `json:"agent,omitempty" gorm:"-" comment:"负责经纪人"`
}

// TableName 设置表名
func (SysHouse) TableName() string {
	return "sys_house"
}

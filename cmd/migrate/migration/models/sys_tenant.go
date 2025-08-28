package models

import (
	"time"
)

// SysTenant 租客模型
type SysTenant struct {
	// 主键
	ID uint `json:"id" gorm:"primaryKey;autoIncrement;type:int unsigned" comment:"主键ID"`

	// 基础信息
	TenantCode string     `json:"tenant_code" gorm:"size:50;not null;uniqueIndex:idx_tenant_code" comment:"租客编号"`
	Name       string     `json:"name" gorm:"size:50;not null;index:idx_name" comment:"姓名"`
	Gender     string     `json:"gender" gorm:"size:10" comment:"性别(男/女)"`
	Birthday   *time.Time `json:"birthday" gorm:"type:date" comment:"出生日期"`
	IDCard     string     `json:"id_card" gorm:"size:18;index:idx_id_card" comment:"身份证号"`

	// 联系信息
	Phone            string `json:"phone" gorm:"size:20;not null;index:idx_phone" comment:"手机号"`
	Email            string `json:"email" gorm:"size:100;index:idx_email" comment:"邮箱"`
	WechatID         string `json:"wechat_id" gorm:"size:100" comment:"微信号"`
	EmergencyContact string `json:"emergency_contact" gorm:"size:50" comment:"紧急联系人"`
	EmergencyPhone   string `json:"emergency_phone" gorm:"size:20" comment:"紧急联系人电话"`
	Avatar           string `json:"avatar" gorm:"size:500" comment:"头像URL"`

	// 工作信息
	Occupation    string `json:"occupation" gorm:"size:100" comment:"职业"`
	Company       string `json:"company" gorm:"size:200" comment:"工作单位"`
	WorkAddress   string `json:"work_address" gorm:"size:300" comment:"工作地址"`
	MonthlyIncome int    `json:"monthly_income" gorm:"comment:" comment:"月收入(元)"`

	// 租房偏好
	PreferredArea      string  `json:"preferred_area" gorm:"size:200" comment:"偏好区域"`
	PreferredHouseType string  `json:"preferred_house_type" gorm:"size:100" comment:"偏好房型"`
	MinRent            int     `json:"min_rent" gorm:"default:0" comment:"最低租金(元)"`
	MaxRent            int     `json:"max_rent" gorm:"default:0" comment:"最高租金(元)"`
	MinArea            float64 `json:"min_area" gorm:"type:decimal(8,2);default:0" comment:"最小面积(平方米)"`
	MaxArea            float64 `json:"max_area" gorm:"type:decimal(8,2);default:0" comment:"最大面积(平方米)"`
	PreferredFloor     string  `json:"preferred_floor" gorm:"size:50" comment:"偏好楼层(低楼层/中楼层/高楼层)"`
	RequiredFacilities string  `json:"required_facilities" gorm:"type:text" comment:"必需设施(JSON格式)"`

	// 租赁历史统计
	RentalCount     int     `json:"rental_count" gorm:"default:0" comment:"租房次数"`
	TotalRentPaid   float64 `json:"total_rent_paid" gorm:"type:decimal(12,2);default:0" comment:"累计租金支付(元)"`
	AverageRentDays int     `json:"average_rent_days" gorm:"default:0" comment:"平均租期(天)"`
	CreditScore     int     `json:"credit_score" gorm:"default:100" comment:"信用评分(0-100)"`

	// 状态信息
	TenantStatus   string     `json:"tenant_status" gorm:"size:20;not null;default:'active';index:idx_tenant_status" comment:"租客状态(active:活跃/inactive:非活跃/blacklist:黑名单)"`
	IsVIP          bool       `json:"is_vip" gorm:"default:false;index:idx_is_vip" comment:"是否VIP客户"`
	RegisterSource string     `json:"register_source" gorm:"size:50" comment:"注册来源(官网/APP/经纪人推荐等)"`
	LastActiveAt   *time.Time `json:"last_active_at" gorm:"" comment:"最后活跃时间"`

	// 个人说明
	PersonalNote string `json:"personal_note" gorm:"type:text" comment:"个人说明/备注"`
	SpecialNeeds string `json:"special_needs" gorm:"type:text" comment:"特殊需求"`

	// 管理信息
	CreatedBy string `json:"created_by" gorm:"size:50" comment:"创建人"`
	UpdatedBy string `json:"updated_by" gorm:"size:50" comment:"更新人"`

	// 时间戳
	CreatedAt *time.Time `json:"created_at" gorm:"autoCreateTime" comment:"创建时间"`
	UpdatedAt *time.Time `json:"updated_at" gorm:"autoUpdateTime" comment:"更新时间"`
	DeletedAt *time.Time `json:"deleted_at" gorm:"index" comment:"删除时间"`

	// 关联对象(用于GORM关联查询)
	RentalRecords []SysRentalRecord `json:"rental_records,omitempty" gorm:"-" comment:"租赁记录列表"`
}

// TableName 设置表名
func (SysTenant) TableName() string {
	return "sys_tenant"
}

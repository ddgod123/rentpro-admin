package models

import (
	"time"
)

// SysAgent 经纪人模型
type SysAgent struct {
	// 主键
	ID uint `json:"id" gorm:"primaryKey;autoIncrement;type:int unsigned" comment:"主键ID"`

	// 基础信息
	AgentCode string     `json:"agent_code" gorm:"size:50;not null;uniqueIndex:idx_agent_code" comment:"经纪人编号"`
	Name      string     `json:"name" gorm:"size:50;not null;index:idx_name" comment:"姓名"`
	Gender    string     `json:"gender" gorm:"size:10" comment:"性别(男/女)"`
	Birthday  *time.Time `json:"birthday" gorm:"type:date" comment:"出生日期"`
	IDCard    string     `json:"id_card" gorm:"size:18;index:idx_id_card" comment:"身份证号"`

	// 联系信息
	Phone   string `json:"phone" gorm:"size:20;not null;index:idx_phone" comment:"手机号"`
	Email   string `json:"email" gorm:"size:100;index:idx_email" comment:"邮箱"`
	Address string `json:"address" gorm:"size:300" comment:"居住地址"`
	Avatar  string `json:"avatar" gorm:"size:500" comment:"头像URL"`

	// 工作信息
	JobTitle   string     `json:"job_title" gorm:"size:50;not null" comment:"职位(初级经纪人/高级经纪人/店长等)"`
	Department string     `json:"department" gorm:"size:100" comment:"所属部门"`
	HireDate   *time.Time `json:"hire_date" gorm:"type:date" comment:"入职日期"`
	WorkYears  int        `json:"work_years" gorm:"default:0" comment:"从业年限"`
	LicenseNo  string     `json:"license_no" gorm:"size:100" comment:"执业证书号"`

	// 业务信息
	ServiceArea    string  `json:"service_area" gorm:"size:200" comment:"服务区域"`
	Specialization string  `json:"specialization" gorm:"size:200" comment:"专业领域"`
	Commission     float64 `json:"commission" gorm:"type:decimal(5,2);default:3.00" comment:"佣金比例(%)"`

	// 统计信息
	TotalDeals     int     `json:"total_deals" gorm:"default:0" comment:"总成交量"`
	MonthlyDeals   int     `json:"monthly_deals" gorm:"default:0" comment:"本月成交量"`
	TotalRevenue   float64 `json:"total_revenue" gorm:"type:decimal(12,2);default:0" comment:"总业绩(元)"`
	MonthlyRevenue float64 `json:"monthly_revenue" gorm:"type:decimal(12,2);default:0" comment:"本月业绩(元)"`
	Rating         float64 `json:"rating" gorm:"type:decimal(3,2);default:5.00" comment:"用户评分(1-5分)"`
	ReviewCount    int     `json:"review_count" gorm:"default:0" comment:"评价数量"`

	// 状态信息
	Status      string     `json:"status" gorm:"size:20;not null;default:'active';index:idx_status" comment:"状态(active:在职/leave:请假/resigned:离职)"`
	IsTopAgent  bool       `json:"is_top_agent" gorm:"default:false;index:idx_is_top_agent" comment:"是否金牌经纪人"`
	IsOnline    bool       `json:"is_online" gorm:"default:false" comment:"是否在线"`
	LastLoginAt *time.Time `json:"last_login_at" gorm:"" comment:"最后登录时间"`

	// 个人介绍
	Introduction string `json:"introduction" gorm:"type:text" comment:"个人介绍"`
	Achievements string `json:"achievements" gorm:"type:text" comment:"主要成就"`

	// 管理信息
	CreatedBy string `json:"created_by" gorm:"size:50" comment:"创建人"`
	UpdatedBy string `json:"updated_by" gorm:"size:50" comment:"更新人"`

	// 时间戳
	CreatedAt *time.Time `json:"created_at" gorm:"autoCreateTime" comment:"创建时间"`
	UpdatedAt *time.Time `json:"updated_at" gorm:"autoUpdateTime" comment:"更新时间"`
	DeletedAt *time.Time `json:"deleted_at" gorm:"index" comment:"删除时间"`

	// 关联对象(用于GORM关联查询)
	Houses []SysHouse `json:"houses,omitempty" gorm:"-" comment:"管理的房源列表"`
}

// TableName 设置表名
func (SysAgent) TableName() string {
	return "sys_agent"
}

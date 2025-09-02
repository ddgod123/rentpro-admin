package rental

import (
	"time"
)

// SysAgent 经纪人模型
type SysAgent struct {
	// 主键
	ID uint `json:"id" gorm:"primaryKey;autoIncrement" comment:"主键ID"`

	// 基础信息
	Name    string `json:"name" gorm:"size:100;not null;index:idx_name" comment:"经纪人姓名"`
	Phone   string `json:"phone" gorm:"size:20;not null;uniqueIndex:idx_phone" comment:"联系电话"`
	IDCard  string `json:"idCard" gorm:"size:18;uniqueIndex:idx_id_card" comment:"身份证号"`
	Email   string `json:"email" gorm:"size:100;uniqueIndex:idx_email" comment:"邮箱"`
	Address string `json:"address" gorm:"size:500" comment:"联系地址"`

	// 所属公司
	CompanyID   uint   `json:"companyId" gorm:"index:idx_company_id" comment:"所属公司ID"`
	CompanyName string `json:"companyName" gorm:"size:100" comment:"所属公司名称"`

	// 证书信息
	CertificationNumber string     `json:"certificationNumber" gorm:"size:50;uniqueIndex:idx_cert_number" comment:"从业资格证书编号"`
	CertificationDate   *time.Time `json:"certificationDate" gorm:"comment:"资格证书获得日期"`
	CertificationImage  string     `json:"certificationImage" gorm:"size:500" comment:"资格证书图片URL"`

	// 专业信息
	Specialization string `json:"specialization" gorm:"size:100" comment:"专业领域(住宅/商业/办公等)"`
	Experience     int    `json:"experience" gorm:"default:0" comment:"从业经验(年)"`

	// 业绩统计
	TotalDeals      int     `json:"totalDeals" gorm:"default:0;index:idx_total_deals" comment:"总成交数"`
	TotalCommission float64 `json:"totalCommission" gorm:"type:decimal(12,2);default:0" comment:"总佣金收入"`
	AverageRating   float64 `json:"averageRating" gorm:"type:decimal(3,2);default:0;index:idx_avg_rating" comment:"平均评分"`

	// 状态信息
	Status string `json:"status" gorm:"size:20;not null;default:'active';index:idx_status" comment:"状态(active:正常, inactive:停用, suspended:暂停)"`

	// 备注
	Notes string `json:"notes" gorm:"type:text" comment:"备注信息"`

	// 管理信息
	CreatedBy string `json:"createdBy" gorm:"size:50" comment:"创建人"`
	UpdatedBy string `json:"updatedBy" gorm:"size:50" comment:"更新人"`

	// 时间戳
	CreatedAt *time.Time `json:"createdAt" gorm:"autoCreateTime" comment:"创建时间"`
	UpdatedAt *time.Time `json:"updatedAt" gorm:"autoUpdateTime" comment:"更新时间"`
	DeletedAt *time.Time `json:"deletedAt" gorm:"index" comment:"删除时间"`
}

// TableName 设置表名
func (SysAgent) TableName() string {
	return "sys_agents"
}

// GetStatusText 获取状态文本描述
func (a *SysAgent) GetStatusText() string {
	switch a.Status {
	case "active":
		return "正常"
	case "inactive":
		return "停用"
	case "suspended":
		return "暂停"
	default:
		return "未知"
	}
}

// GetSpecializationText 获取专业领域文本描述
func (a *SysAgent) GetSpecializationText() string {
	switch a.Specialization {
	case "residential":
		return "住宅"
	case "commercial":
		return "商业"
	case "office":
		return "办公"
	default:
		return a.Specialization
	}
}

// GetExperienceText 获取经验描述
func (a *SysAgent) GetExperienceText() string {
	if a.Experience == 0 {
		return "新手"
	} else if a.Experience < 3 {
		return "1-3年经验"
	} else if a.Experience < 5 {
		return "3-5年经验"
	} else if a.Experience < 10 {
		return "5-10年经验"
	} else {
		return "10年以上经验"
	}
}

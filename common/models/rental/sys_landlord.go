package rental

import (
	"time"
)

// SysLandlord 房东模型
type SysLandlord struct {
	// 主键
	ID uint `json:"id" gorm:"primaryKey;autoIncrement" comment:"主键ID"`

	// 基础信息
	Name             string `json:"name" gorm:"size:100;not null;index:idx_name" comment:"房东姓名"`
	Phone            string `json:"phone" gorm:"size:20;not null;uniqueIndex:idx_phone" comment:"联系电话"`
	IDCard           string `json:"idCard" gorm:"size:18;uniqueIndex:idx_id_card" comment:"身份证号"`
	Email            string `json:"email" gorm:"size:100;uniqueIndex:idx_email" comment:"邮箱"`
	Address          string `json:"address" gorm:"size:500" comment:"联系地址"`
	EmergencyContact string `json:"emergencyContact" gorm:"size:100" comment:"紧急联系人"`
	EmergencyPhone   string `json:"emergencyPhone" gorm:"size:20" comment:"紧急联系电话"`

	// 公司信息（如果是企业房东）
	CompanyName     string `json:"companyName" gorm:"size:100" comment:"公司名称"`
	CompanyAddress  string `json:"companyAddress" gorm:"size:500" comment:"公司地址"`
	BusinessLicense string `json:"businessLicense" gorm:"size:100" comment:"营业执照号"`

	// 房东类型和状态
	Type   string `json:"type" gorm:"size:20;not null;default:'individual';index:idx_type" comment:"房东类型(individual:个人, company:企业)"`
	Status string `json:"status" gorm:"size:20;not null;default:'active';index:idx_status" comment:"状态(active:正常, inactive:停用, blacklisted:黑名单)"`

	// 房产信息
	PropertyCount int     `json:"propertyCount" gorm:"default:0;index:idx_property_count" comment:"房产数量"`
	TotalArea     float64 `json:"totalArea" gorm:"type:decimal(10,2);default:0" comment:"总面积(平方米)"`

	// 收益信息
	TotalIncome   float64 `json:"totalIncome" gorm:"type:decimal(12,2);default:0" comment:"总收入"`
	AverageIncome float64 `json:"averageIncome" gorm:"type:decimal(10,2);default:0" comment:"平均月收入"`

	// 信用信息
	CreditScore   int  `json:"creditScore" gorm:"default:100" comment:"信用评分(0-100)"`
	IsVIP         bool `json:"isVIP" gorm:"default:false;index:idx_is_vip" comment:"是否VIP房东"`
	IsBlacklisted bool `json:"isBlacklisted" gorm:"default:false;index:idx_is_blacklisted" comment:"是否黑名单"`

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
func (SysLandlord) TableName() string {
	return "sys_landlords"
}

// GetStatusText 获取状态文本描述
func (l *SysLandlord) GetStatusText() string {
	switch l.Status {
	case "active":
		return "正常"
	case "inactive":
		return "停用"
	case "blacklisted":
		return "黑名单"
	default:
		return "未知"
	}
}

// GetTypeText 获取类型文本描述
func (l *SysLandlord) GetTypeText() string {
	switch l.Type {
	case "individual":
		return "个人"
	case "company":
		return "企业"
	default:
		return "未知"
	}
}

// IsIndividual 判断是否为个人房东
func (l *SysLandlord) IsIndividual() bool {
	return l.Type == "individual"
}

// IsCompany 判断是否为企业房东
func (l *SysLandlord) IsCompany() bool {
	return l.Type == "company"
}

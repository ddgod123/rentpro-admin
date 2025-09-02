package rental

import (
	"time"
)

// SysContract 合同模型
type SysContract struct {
	// 主键
	ID uint `json:"id" gorm:"primaryKey;autoIncrement" comment:"主键ID"`

	// 合同基本信息
	ContractNumber string `json:"contractNumber" gorm:"size:50;not null;uniqueIndex:idx_contract_number" comment:"合同编号"`
	Title          string `json:"title" gorm:"size:200;not null" comment:"合同标题"`
	Type           string `json:"type" gorm:"size:20;not null;index:idx_type" comment:"合同类型(rent:租赁, sale:买卖)"`

	// 关联信息
	PropertyID   uint   `json:"propertyId" gorm:"not null;index:idx_property_id" comment:"房源ID"`
	PropertyType string `json:"propertyType" gorm:"size:20;not null" comment:"房源类型(building:楼盘, house:房屋)"`
	TenantID     uint   `json:"tenantId" gorm:"not null;index:idx_tenant_id" comment:"租户ID"`
	LandlordID   uint   `json:"landlordId" gorm:"not null;index:idx_landlord_id" comment:"房东ID"`
	AgentID      uint   `json:"agentId" gorm:"index:idx_agent_id" comment:"经纪人ID"`

	// 时间信息
	StartDate     *time.Time `json:"startDate" gorm:"not null" comment:"合同开始日期"`
	EndDate       *time.Time `json:"endDate" gorm:"not null" comment:"合同结束日期"`
	SigningDate   *time.Time `json:"signingDate" gorm:"not null" comment:"签约日期"`
	EffectiveDate *time.Time `json:"effectiveDate" comment:"生效日期"`

	// 金额信息
	RentAmount float64 `json:"rentAmount" gorm:"type:decimal(10,2);not null" comment:"租金金额"`
	Deposit    float64 `json:"deposit" gorm:"type:decimal(10,2)" comment:"押金"`
	Commission float64 `json:"commission" gorm:"type:decimal(10,2)" comment:"佣金"`
	OtherFees  float64 `json:"otherFees" gorm:"type:decimal(10,2)" comment:"其他费用"`

	// 支付信息
	PaymentCycle    string     `json:"paymentCycle" gorm:"size:20" comment:"支付周期(monthly:月付, quarterly:季付, yearly:年付)"`
	NextPaymentDate *time.Time `json:"nextPaymentDate" comment:"下次支付日期"`

	// 合同状态
	Status string `json:"status" gorm:"size:20;not null;default:'pending';index:idx_status" comment:"合同状态(pending:待生效, active:生效中, expired:已过期, terminated:已终止, cancelled:已取消)"`

	// 房屋信息
	Address string  `json:"address" gorm:"size:500" comment:"房屋地址"`
	Area    float64 `json:"area" gorm:"type:decimal(8,2)" comment:"房屋面积(平方米)"`

	// 附件信息
	ContractFile string `json:"contractFile" gorm:"size:500" comment:"合同文件URL"`
	Attachments  string `json:"attachments" gorm:"type:text" comment:"附件信息(JSON格式)"`

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
func (SysContract) TableName() string {
	return "sys_contracts"
}

// GetStatusText 获取状态文本描述
func (c *SysContract) GetStatusText() string {
	switch c.Status {
	case "pending":
		return "待生效"
	case "active":
		return "生效中"
	case "expired":
		return "已过期"
	case "terminated":
		return "已终止"
	case "cancelled":
		return "已取消"
	default:
		return "未知"
	}
}

// GetTypeText 获取类型文本描述
func (c *SysContract) GetTypeText() string {
	switch c.Type {
	case "rent":
		return "租赁"
	case "sale":
		return "买卖"
	default:
		return "未知"
	}
}

// GetPaymentCycleText 获取支付周期文本描述
func (c *SysContract) GetPaymentCycleText() string {
	switch c.PaymentCycle {
	case "monthly":
		return "月付"
	case "quarterly":
		return "季付"
	case "yearly":
		return "年付"
	default:
		return "未知"
	}
}

// IsRent 判断是否为租赁合同
func (c *SysContract) IsRent() bool {
	return c.Type == "rent"
}

// IsSale 判断是否为买卖合同
func (c *SysContract) IsSale() bool {
	return c.Type == "sale"
}

// IsActive 判断合同是否生效中
func (c *SysContract) IsActive() bool {
	return c.Status == "active"
}

// IsExpired 判断合同是否已过期
func (c *SysContract) IsExpired() bool {
	return c.Status == "expired"
}

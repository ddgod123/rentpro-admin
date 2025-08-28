package models

import (
	"time"
)

// SysRentalRecord 租赁记录模型
type SysRentalRecord struct {
	// 主键
	ID uint `json:"id" gorm:"primaryKey;autoIncrement;type:int unsigned" comment:"主键ID"`

	// 关联信息
	TenantID uint `json:"tenant_id" gorm:"not null;index:idx_tenant_id;type:int unsigned" comment:"租客ID"`
	HouseID  uint `json:"house_id" gorm:"not null;index:idx_house_id;type:int unsigned" comment:"房源ID"`
	AgentID  uint `json:"agent_id" gorm:"index:idx_agent_id;type:int unsigned" comment:"经纪人ID"`

	// 合同信息
	ContractCode  string     `json:"contract_code" gorm:"size:100;not null;uniqueIndex:idx_contract_code" comment:"合同编号"`
	ContractType  string     `json:"contract_type" gorm:"size:20;not null" comment:"合同类型(租赁/续租/转租)"`
	StartDate     *time.Time `json:"start_date" gorm:"type:date;not null" comment:"租赁开始日期"`
	EndDate       *time.Time `json:"end_date" gorm:"type:date;not null" comment:"租赁结束日期"`
	ActualEndDate *time.Time `json:"actual_end_date" gorm:"type:date" comment:"实际结束日期"`
	RentalDays    int        `json:"rental_days" gorm:"not null" comment:"租期天数"`

	// 租金信息
	MonthlyRent   int     `json:"monthly_rent" gorm:"not null" comment:"月租金(元)"`
	Deposit       int     `json:"deposit" gorm:"not null" comment:"押金(元)"`
	ServiceFee    int     `json:"service_fee" gorm:"default:0" comment:"服务费(元)"`
	Commission    float64 `json:"commission" gorm:"type:decimal(10,2);default:0" comment:"佣金(元)"`
	TotalAmount   int     `json:"total_amount" gorm:"not null" comment:"合同总金额(元)"`
	PaymentCycle  string  `json:"payment_cycle" gorm:"size:20;not null" comment:"付款周期(月付/季付/年付)"`
	PaymentMethod string  `json:"payment_method" gorm:"size:50" comment:"付款方式(现金/转账/支付宝等)"`

	// 状态信息
	RentalStatus  string `json:"rental_status" gorm:"size:20;not null;default:'active';index:idx_rental_status" comment:"租赁状态(active:生效中/expired:已到期/terminated:已终止/pending:待生效)"`
	PaymentStatus string `json:"payment_status" gorm:"size:20;not null;default:'current';index:idx_payment_status" comment:"付款状态(current:正常/overdue:逾期/prepaid:预付)"`

	// 评价信息
	TenantRating   float64 `json:"tenant_rating" gorm:"type:decimal(3,2);default:0" comment:"租客评分(1-5分)"`
	HouseRating    float64 `json:"house_rating" gorm:"type:decimal(3,2);default:0" comment:"房源评分(1-5分)"`
	AgentRating    float64 `json:"agent_rating" gorm:"type:decimal(3,2);default:0" comment:"经纪人评分(1-5分)"`
	TenantReview   string  `json:"tenant_review" gorm:"type:text" comment:"租客评价"`
	LandlordReview string  `json:"landlord_review" gorm:"type:text" comment:"房东评价"`

	// 退租信息
	CheckOutDate    *time.Time `json:"check_out_date" gorm:"type:date" comment:"退房日期"`
	DepositRefund   int        `json:"deposit_refund" gorm:"default:0" comment:"押金退还金额(元)"`
	DeductionAmount int        `json:"deduction_amount" gorm:"default:0" comment:"扣除金额(元)"`
	DeductionReason string     `json:"deduction_reason" gorm:"type:text" comment:"扣除原因"`

	// 备注信息
	ContractNote string `json:"contract_note" gorm:"type:text" comment:"合同备注"`
	SpecialTerms string `json:"special_terms" gorm:"type:text" comment:"特殊条款"`

	// 管理信息
	CreatedBy string `json:"created_by" gorm:"size:50" comment:"创建人"`
	UpdatedBy string `json:"updated_by" gorm:"size:50" comment:"更新人"`

	// 时间戳
	CreatedAt *time.Time `json:"created_at" gorm:"autoCreateTime" comment:"创建时间"`
	UpdatedAt *time.Time `json:"updated_at" gorm:"autoUpdateTime" comment:"更新时间"`
	DeletedAt *time.Time `json:"deleted_at" gorm:"index" comment:"删除时间"`

	// 关联对象(用于GORM关联查询)
	Tenant *SysTenant `json:"tenant,omitempty" gorm:"-" comment:"租客信息"`
	House  *SysHouse  `json:"house,omitempty" gorm:"-" comment:"房源信息"`
	Agent  *SysAgent  `json:"agent,omitempty" gorm:"-" comment:"经纪人信息"`
}

// TableName 设置表名
func (SysRentalRecord) TableName() string {
	return "sys_rental_record"
}

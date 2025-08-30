package system

import (
	"time"

	"gorm.io/gorm"
)

// SysPost 系统岗位模型
// 管理系统中的岗位信息
type SysPost struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	PostCode  string         `gorm:"size:64;not null;unique;index:idx_post_code" json:"post_code" comment:"岗位编码"`
	PostName  string         `gorm:"size:128;not null" json:"post_name" comment:"岗位名称"`
	Sort      int            `gorm:"default:1" json:"sort" comment:"排序"`
	Status    string         `gorm:"size:1;default:'0'" json:"status" comment:"状态 0:正常 1:停用"`
	Remark    string         `gorm:"size:255" json:"remark" comment:"备注"`
	CreatedAt time.Time      `json:"created_at" comment:"创建时间"`
	UpdatedAt time.Time      `json:"updated_at" comment:"更新时间"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-" comment:"删除时间"`

	// 关联关系
	Users []SysUser `gorm:"foreignKey:PostID" json:"users,omitempty"`
}

// TableName 设置表名
func (SysPost) TableName() string {
	return "sys_post"
}

// IsActive 检查岗位是否为活跃状态
func (p *SysPost) IsActive() bool {
	return p.Status == "0" && p.DeletedAt.Time.IsZero()
}

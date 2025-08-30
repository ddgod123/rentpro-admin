package system

import (
	"time"

	"gorm.io/gorm"
)

// SysDept 系统部门模型
// 管理组织架构和部门信息
type SysDept struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	ParentID  uint           `gorm:"default:0;index:idx_parent_id" json:"parent_id" comment:"父部门ID"`
	DeptPath  string         `gorm:"size:255" json:"dept_path" comment:"部门路径"`
	DeptName  string         `gorm:"size:128;not null" json:"dept_name" comment:"部门名称"`
	Sort      int            `gorm:"default:1" json:"sort" comment:"排序"`
	Leader    string         `gorm:"size:128" json:"leader" comment:"负责人"`
	Phone     string         `gorm:"size:32" json:"phone" comment:"联系电话"`
	Email     string         `gorm:"size:128" json:"email" comment:"邮箱"`
	Status    string         `gorm:"size:1;default:'0'" json:"status" comment:"状态 0:正常 1:停用"`
	CreatedAt time.Time      `json:"created_at" comment:"创建时间"`
	UpdatedAt time.Time      `json:"updated_at" comment:"更新时间"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-" comment:"删除时间"`

	// 关联关系
	Children []SysDept `gorm:"-" json:"children,omitempty"`
	Parent   *SysDept  `gorm:"-" json:"parent,omitempty"`
	Users    []SysUser `gorm:"foreignKey:DeptID" json:"users,omitempty"`
}

// TableName 设置表名
func (SysDept) TableName() string {
	return "sys_dept"
}

// IsActive 检查部门是否为活跃状态
func (d *SysDept) IsActive() bool {
	return d.Status == "0" && d.DeletedAt.Time.IsZero()
}

// GetFullPath 获取完整部门路径
func (d *SysDept) GetFullPath() string {
	if d.DeptPath != "" {
		return d.DeptPath
	}

	if d.Parent != nil {
		return d.Parent.GetFullPath() + "/" + d.DeptName
	}
	return d.DeptName
}

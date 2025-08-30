package system

import (
	"time"

	"gorm.io/gorm"
)

// SysRole 系统角色模型
// 定义系统中的角色和权限组
type SysRole struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	Name      string         `gorm:"size:128;not null;unique;index:idx_name" json:"name" comment:"角色名称"`
	Key       string         `gorm:"size:128;not null;unique;index:idx_key" json:"key" comment:"角色标识"`
	Status    int            `gorm:"default:1" json:"status" comment:"状态 1:启用 2:禁用"`
	Sort      int            `gorm:"default:1" json:"sort" comment:"排序"`
	Flag      string         `gorm:"size:128" json:"flag" comment:"角色标志"`
	Remark    string         `gorm:"size:255" json:"remark" comment:"备注"`
	Admin     bool           `gorm:"default:false" json:"admin" comment:"是否为管理员角色"`
	DataScope string         `gorm:"size:128;default:'1'" json:"data_scope" comment:"数据权限范围"`
	Params    string         `gorm:"size:255" json:"params" comment:"角色参数"`
	CreatedAt time.Time      `json:"created_at" comment:"创建时间"`
	UpdatedAt time.Time      `json:"updated_at" comment:"更新时间"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-" comment:"删除时间"`

	// 关联关系
	Users []SysUser `gorm:"foreignKey:RoleID" json:"users,omitempty"`
	Menus []SysMenu `gorm:"many2many:sys_role_menu;" json:"menus,omitempty"`
}

// TableName 设置表名
func (SysRole) TableName() string {
	return "sys_role"
}

// IsActive 检查角色是否为活跃状态
func (r *SysRole) IsActive() bool {
	return r.Status == 1 && r.DeletedAt.Time.IsZero()
}

// IsAdmin 检查是否为管理员角色
func (r *SysRole) IsAdmin() bool {
	return r.Admin || r.Key == "admin" || r.Key == "super_admin"
}

// HasPermission 检查角色是否有指定权限（通过菜单权限）
func (r *SysRole) HasPermission(permission string) bool {
	// 管理员默认拥有所有权限
	if r.IsAdmin() {
		return true
	}

	// 检查菜单权限
	for _, menu := range r.Menus {
		if menu.Permission == permission && menu.IsActive() {
			return true
		}
	}

	return false
}

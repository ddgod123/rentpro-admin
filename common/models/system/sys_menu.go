package system

import (
	"time"

	"gorm.io/gorm"
)

// SysMenu 系统菜单/权限模型
// 管理系统菜单和权限点
type SysMenu struct {
	ID         uint           `gorm:"primarykey" json:"id"`
	Name       string         `gorm:"size:128;not null" json:"name" comment:"菜单名称"`
	Title      string         `gorm:"size:128" json:"title" comment:"菜单标题"`
	Icon       string         `gorm:"size:128" json:"icon" comment:"菜单图标"`
	Path       string         `gorm:"size:128" json:"path" comment:"路由地址"`
	Redirect   string         `gorm:"size:128" json:"redirect" comment:"重定向地址"`
	Component  string         `gorm:"size:128" json:"component" comment:"组件路径"`
	Permission string         `gorm:"size:255" json:"permission" comment:"权限标识"`
	ParentID   uint           `gorm:"default:0;index:idx_parent_id" json:"parent_id" comment:"父菜单ID"`
	Type       string         `gorm:"size:1;default:'M'" json:"type" comment:"菜单类型 M:菜单 C:目录 F:按钮"`
	Sort       int            `gorm:"default:1" json:"sort" comment:"排序"`
	Visible    string         `gorm:"size:1;default:'0'" json:"visible" comment:"是否显示 0:显示 1:隐藏"`
	IsFrame    string         `gorm:"size:1;default:'1'" json:"is_frame" comment:"是否外链 0:是 1:否"`
	IsCache    string         `gorm:"size:1;default:'0'" json:"is_cache" comment:"是否缓存 0:缓存 1:不缓存"`
	MenuType   string         `gorm:"size:1;default:''" json:"menu_type" comment:"菜单类型 1:左侧菜单 2:顶部菜单 3:按钮"`
	Status     string         `gorm:"size:1;default:'0'" json:"status" comment:"状态 0:正常 1:停用"`
	Perms      string         `gorm:"size:100" json:"perms" comment:"权限字符串"`
	CreatedAt  time.Time      `json:"created_at" comment:"创建时间"`
	UpdatedAt  time.Time      `json:"updated_at" comment:"更新时间"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-" comment:"删除时间"`

	// 关联关系
	Children []SysMenu `gorm:"-" json:"children,omitempty"`
	Parent   *SysMenu  `gorm:"-" json:"parent,omitempty"`
	Roles    []SysRole `gorm:"many2many:sys_role_menu;" json:"roles,omitempty"`
}

// TableName 设置表名
func (SysMenu) TableName() string {
	return "sys_menu"
}

// IsActive 检查菜单是否为活跃状态
func (m *SysMenu) IsActive() bool {
	return m.Status == "0" && m.DeletedAt.Time.IsZero()
}

// IsVisible 检查菜单是否可见
func (m *SysMenu) IsVisible() bool {
	return m.Visible == "0"
}

// IsMenu 是否为菜单类型
func (m *SysMenu) IsMenu() bool {
	return m.Type == "M"
}

// IsDirectory 是否为目录类型
func (m *SysMenu) IsDirectory() bool {
	return m.Type == "C"
}

// IsButton 是否为按钮类型
func (m *SysMenu) IsButton() bool {
	return m.Type == "F"
}

// GetFullPath 获取完整路径
func (m *SysMenu) GetFullPath() string {
	if m.Parent != nil {
		return m.Parent.GetFullPath() + "/" + m.Path
	}
	return m.Path
}

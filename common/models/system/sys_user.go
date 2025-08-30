package system

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// SysUser 系统用户模型
// 管理系统中所有用户的基本信息和权限
type SysUser struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	Username    string         `gorm:"size:64;not null;unique;index:idx_username" json:"username" comment:"用户名"`
	Password    string         `gorm:"size:128;not null" json:"-" comment:"密码"`
	NickName    string         `gorm:"size:128" json:"nick_name" comment:"昵称"`
	Avatar      string         `gorm:"size:255" json:"avatar" comment:"头像"`
	Email       string         `gorm:"size:128;index:idx_email" json:"email" comment:"邮箱"`
	Phone       string         `gorm:"size:32;index:idx_phone" json:"phone" comment:"手机号"`
	Status      int            `gorm:"default:1" json:"status" comment:"状态 1:启用 2:禁用"`
	IsAdmin     bool           `gorm:"default:false" json:"is_admin" comment:"是否为管理员"`
	Remark      string         `gorm:"size:255" json:"remark" comment:"备注"`
	DeptID      uint           `gorm:"default:0" json:"dept_id" comment:"部门ID"`
	PostID      uint           `gorm:"default:0" json:"post_id" comment:"岗位ID"`
	RoleID      uint           `gorm:"default:0" json:"role_id" comment:"角色ID"`
	Salt        string         `gorm:"size:255" json:"-" comment:"加密盐"`
	LastLoginIP string         `gorm:"size:128" json:"last_login_ip" comment:"最后登录IP"`
	LastLoginAt *time.Time     `json:"last_login_at" comment:"最后登录时间"`
	CreatedAt   time.Time      `json:"created_at" comment:"创建时间"`
	UpdatedAt   time.Time      `json:"updated_at" comment:"更新时间"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-" comment:"删除时间"`

	// 关联关系
	Role *SysRole `gorm:"foreignKey:RoleID" json:"role,omitempty"`
	Dept *SysDept `gorm:"foreignKey:DeptID" json:"dept,omitempty"`
	Post *SysPost `gorm:"foreignKey:PostID" json:"post,omitempty"`
}

// TableName 设置表名
func (SysUser) TableName() string {
	return "sys_user"
}

// BeforeCreate 创建前的钩子函数
func (u *SysUser) BeforeCreate(tx *gorm.DB) error {
	// 如果没有设置密码，使用默认密码
	if u.Password == "" {
		u.Password = "123456"
	}

	// 加密密码
	return u.Encrypt()
}

// BeforeUpdate 更新前的钩子函数
func (u *SysUser) BeforeUpdate(tx *gorm.DB) error {
	// 如果密码被修改，重新加密
	if tx.Statement.Changed("Password") && u.Password != "" {
		return u.Encrypt()
	}
	return nil
}

// Encrypt 加密密码
func (u *SysUser) Encrypt() error {
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hash)
	return nil
}

// ComparePassword 验证密码
func (u *SysUser) ComparePassword(password string) bool {
	// 首先尝试使用bcrypt验证（新密码）
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err == nil {
		return true
	}

	// 如果bcrypt验证失败，尝试明文比较（旧密码）
	// 这是为了兼容数据库中可能存在的明文密码
	if u.Password == password {
		// 如果明文匹配，重新加密密码以更新为加密格式
		if encryptErr := u.Encrypt(); encryptErr == nil {
			// 这里应该更新数据库中的密码，但在当前上下文中无法访问数据库
			// 在实际应用中，可能需要在登录成功后更新密码
		}
		return true
	}

	return false
}

// IsActive 检查用户是否为活跃状态
func (u *SysUser) IsActive() bool {
	return u.Status == 1 && u.DeletedAt.Time.IsZero()
}

// GetDisplayName 获取显示名称
func (u *SysUser) GetDisplayName() string {
	if u.NickName != "" {
		return u.NickName
	}
	return u.Username
}

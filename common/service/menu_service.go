package service

import (
	"fmt"
	"rentPro/rentpro-admin/common/models/system"
	"sort"

	"gorm.io/gorm"
)

// MenuService 菜单服务
type MenuService struct {
	DB *gorm.DB
}

// NewMenuService 创建菜单服务实例
func NewMenuService(db *gorm.DB) *MenuService {
	return &MenuService{
		DB: db,
	}
}

// GetMenusByRoleID 根据角色ID获取菜单列表
func (s *MenuService) GetMenusByRoleID(roleID uint) ([]system.SysMenu, error) {
	var menus []system.SysMenu

	// 通过角色菜单关联表查询菜单
	err := s.DB.Table("sys_menu").
		Joins("JOIN sys_role_menu ON sys_menu.id = sys_role_menu.sys_menu_id").
		Where("sys_role_menu.sys_role_id = ? AND sys_menu.visible = '0' AND sys_menu.status = '0'", roleID).
		Order("sys_menu.sort ASC").
		Find(&menus).Error

	return menus, err
}

// GetMenusByUserID 根据用户ID获取菜单列表
func (s *MenuService) GetMenusByUserID(userID uint) ([]system.SysMenu, error) {
	var menus []system.SysMenu

	// 通过用户角色关联表查询菜单
	err := s.DB.Table("sys_menu").
		Joins("JOIN sys_role_menu ON sys_menu.id = sys_role_menu.sys_menu_id").
		Joins("JOIN sys_user_role ON sys_role_menu.sys_role_id = sys_user_role.sys_role_id").
		Where("sys_user_role.sys_user_id = ? AND sys_menu.visible = '0' AND sys_menu.status = '0'", userID).
		Order("sys_menu.sort ASC").
		Find(&menus).Error

	return menus, err
}

// BuildMenuTree 构建菜单树结构
func (s *MenuService) BuildMenuTree(menus []system.SysMenu) []system.SysMenu {
	menuMap := make(map[uint]*system.SysMenu)
	var rootMenus []system.SysMenu

	// 创建菜单映射
	for i := range menus {
		menuMap[menus[i].ID] = &menus[i]
	}

	// 构建树结构
	for _, menu := range menus {
		if menu.ParentID == 0 {
			// 根菜单
			rootMenus = append(rootMenus, menu)
		} else {
			// 子菜单，添加到父菜单的children中
			if parent, exists := menuMap[menu.ParentID]; exists {
				parent.Children = append(parent.Children, menu)
			}
		}
	}

	// 对菜单进行排序
	s.sortMenus(&rootMenus)

	return rootMenus
}

// sortMenus 递归排序菜单
func (s *MenuService) sortMenus(menus *[]system.SysMenu) {
	sort.Slice(*menus, func(i, j int) bool {
		return (*menus)[i].Sort < (*menus)[j].Sort
	})

	for i := range *menus {
		if len((*menus)[i].Children) > 0 {
			s.sortMenus(&(*menus)[i].Children)
		}
	}
}

// GetMenuPermissions 获取菜单权限列表
func (s *MenuService) GetMenuPermissions(menus []system.SysMenu) []string {
	permissions := make([]string, 0)
	for _, menu := range menus {
		if menu.Permission != "" {
			permissions = append(permissions, menu.Permission)
		}
	}
	return permissions
}

// GetAllMenus 获取所有菜单
func (s *MenuService) GetAllMenus() ([]system.SysMenu, error) {
	var menus []system.SysMenu
	err := s.DB.Where("visible = '0' AND status = '0'").Order("sort ASC").Find(&menus).Error
	return menus, err
}

// GetMenuByID 根据ID获取菜单
func (s *MenuService) GetMenuByID(id uint) (*system.SysMenu, error) {
	var menu system.SysMenu
	err := s.DB.First(&menu, id).Error
	if err != nil {
		return nil, err
	}
	return &menu, nil
}

// CreateMenu 创建菜单
func (s *MenuService) CreateMenu(menu *system.SysMenu) error {
	return s.DB.Create(menu).Error
}

// UpdateMenu 更新菜单
func (s *MenuService) UpdateMenu(menu *system.SysMenu) error {
	return s.DB.Save(menu).Error
}

// DeleteMenu 删除菜单
func (s *MenuService) DeleteMenu(id uint) error {
	// 检查是否有子菜单
	var count int64
	if err := s.DB.Model(&system.SysMenu{}).Where("parent_id = ?", id).Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		return fmt.Errorf("存在子菜单，无法删除")
	}

	// 删除菜单
	if err := s.DB.Delete(&system.SysMenu{}, id).Error; err != nil {
		return err
	}

	// 删除角色菜单关联
	return s.DB.Exec("DELETE FROM sys_role_menu WHERE sys_menu_id = ?", id).Error
}

// GetMenuRoutes 获取菜单路由
func (s *MenuService) GetMenuRoutes(userID uint) ([]system.SysMenu, error) {
	return s.GetMenusByUserID(userID)
}

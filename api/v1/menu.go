package v1

import (
	"net/http"
	"rentPro/rentpro-admin/common/models/base"
	"rentPro/rentpro-admin/common/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

// MenuController 菜单控制器
type MenuController struct {
	menuService *service.MenuService
}

// NewMenuController 创建菜单控制器实例
func NewMenuController(menuService *service.MenuService) *MenuController {
	return &MenuController{
		menuService: menuService,
	}
}

// GetUserMenus 获取用户菜单
func (c *MenuController) GetUserMenus(ctx *gin.Context) {
	// 从JWT中获取用户ID
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"code": 401,
			"msg":  "未授权",
		})
		return
	}

	// 获取用户菜单
	menus, err := c.menuService.GetMenusByUserID(userID.(uint))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "获取菜单失败: " + err.Error(),
		})
		return
	}

	// 构建菜单树
	menuTree := c.menuService.BuildMenuTree(menus)

	// 获取权限列表
	permissions := c.menuService.GetMenuPermissions(menus)

	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "获取成功",
		"data": gin.H{
			"menus":       menuTree,
			"permissions": permissions,
		},
	})
}

// GetRoleMenus 获取角色菜单
func (c *MenuController) GetRoleMenus(ctx *gin.Context) {
	roleIDStr := ctx.Param("roleId")
	roleID, err := strconv.ParseUint(roleIDStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "角色ID格式错误",
		})
		return
	}

	// 获取角色菜单
	menus, err := c.menuService.GetMenusByRoleID(uint(roleID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "获取菜单失败: " + err.Error(),
		})
		return
	}

	// 构建菜单树
	menuTree := c.menuService.BuildMenuTree(menus)

	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "获取成功",
		"data": menuTree,
	})
}

// GetAllMenus 获取所有菜单（管理员用）
func (c *MenuController) GetAllMenus(ctx *gin.Context) {
	// 获取所有菜单
	menus, err := c.menuService.GetAllMenus()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "获取菜单失败: " + err.Error(),
		})
		return
	}

	// 构建菜单树
	menuTree := c.menuService.BuildMenuTree(menus)

	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "获取成功",
		"data": menuTree,
	})
}

// CreateMenu 创建菜单
func (c *MenuController) CreateMenu(ctx *gin.Context) {
	var menu system.SysMenu
	if err := ctx.ShouldBindJSON(&menu); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "请求参数错误: " + err.Error(),
		})
		return
	}

	// 创建菜单
	if err := c.menuService.DB.Create(&menu).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "创建菜单失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "创建成功",
		"data": menu,
	})
}

// UpdateMenu 更新菜单
func (c *MenuController) UpdateMenu(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "菜单ID格式错误",
		})
		return
	}

	var menu system.SysMenu
	if err := ctx.ShouldBindJSON(&menu); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "请求参数错误: " + err.Error(),
		})
		return
	}

	menu.ID = uint(id)

	// 更新菜单
	if err := c.menuService.DB.Save(&menu).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "更新菜单失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "更新成功",
		"data": menu,
	})
}

// DeleteMenu 删除菜单
func (c *MenuController) DeleteMenu(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "菜单ID格式错误",
		})
		return
	}

	// 检查是否有子菜单
	var count int64
	if err := c.menuService.DB.Model(&system.SysMenu{}).Where("parent_id = ?", id).Count(&count).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "检查子菜单失败: " + err.Error(),
		})
		return
	}

	if count > 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "存在子菜单，无法删除",
		})
		return
	}

	// 删除菜单
	if err := c.menuService.DB.Delete(&system.SysMenu{}, id).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "删除菜单失败: " + err.Error(),
		})
		return
	}

	// 删除角色菜单关联
	if err := c.menuService.DB.Exec("DELETE FROM sys_role_menu WHERE sys_menu_id = ?", id).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "删除角色菜单关联失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "删除成功",
	})
}

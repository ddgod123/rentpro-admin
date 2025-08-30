package v1

import (
	"rentPro/rentpro-admin/common/service"

	"github.com/gin-gonic/gin"
)

// RegisterMenuRoutes 注册菜单相关路由
func RegisterMenuRoutes(r *gin.RouterGroup, menuService *service.MenuService) {
	menuController := NewMenuController(menuService)
	
	// 菜单管理路由组
	menuGroup := r.Group("/menu")
	{
		// 获取用户菜单
		menuGroup.GET("/user", menuController.GetUserMenus)
		
		// 获取角色菜单
		menuGroup.GET("/role/:roleId", menuController.GetRoleMenus)
		
		// 获取所有菜单（管理员用）
		menuGroup.GET("/all", menuController.GetAllMenus)
		
		// 创建菜单
		menuGroup.POST("", menuController.CreateMenu)
		
		// 更新菜单
		menuGroup.PUT("/:id", menuController.UpdateMenu)
		
		// 删除菜单
		menuGroup.DELETE("/:id", menuController.DeleteMenu)
	}
}

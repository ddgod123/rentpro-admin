package api

import (
	"net/http"
	"strconv"

	"rentPro/rentpro-admin/common/database"
	"rentPro/rentpro-admin/common/models/system"
	"rentPro/rentpro-admin/common/service"

	"github.com/gin-gonic/gin"
)

// MenuController 菜单控制器
type MenuController struct {
	menuService *service.MenuService
}

// NewMenuController 创建菜单控制器实例
func NewMenuController() *MenuController {
	return &MenuController{
		menuService: service.NewMenuService(database.DB),
	}
}

// GetUserMenus 获取用户菜单
// @Summary 获取用户菜单
// @Description 获取当前登录用户的菜单权限
// @Tags 菜单
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} Response{data=[]service.MenuRoute}
// @Failure 401 {object} Response
// @Router /api/menu/user-menus [get]
func (c *MenuController) GetUserMenus(ctx *gin.Context) {
	// 从上下文获取用户ID
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, Response{
			Code: 401,
			Msg:  "未授权访问",
		})
		return
	}

	uid, ok := userID.(uint)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, Response{
			Code: 401,
			Msg:  "用户信息错误",
		})
		return
	}

	menus, err := c.menuService.GetMenusByUserID(uid)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{
			Code: 500,
			Msg:  "获取菜单失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, Response{
		Code: 200,
		Msg:  "获取菜单成功",
		Data: menus,
	})
}

// GetAllMenus 获取所有菜单（管理员用）
// @Summary 获取所有菜单
// @Description 获取系统中所有菜单（需要管理员权限）
// @Tags 菜单
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} Response{data=[]service.MenuRoute}
// @Failure 401 {object} Response
// @Failure 403 {object} Response
// @Router /api/menu/all [get]
func (c *MenuController) GetAllMenus(ctx *gin.Context) {
	menus, err := c.menuService.GetAllMenus()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{
			Code: 500,
			Msg:  "获取菜单失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, Response{
		Code: 200,
		Msg:  "获取菜单成功",
		Data: menus,
	})
}

// GetMenuByID 根据ID获取菜单
// @Summary 获取菜单详情
// @Description 根据菜单ID获取菜单详细信息
// @Tags 菜单
// @Accept json
// @Produce json
// @Param id path int true "菜单ID"
// @Security ApiKeyAuth
// @Success 200 {object} Response{data=service.MenuRoute}
// @Failure 400 {object} Response
// @Failure 404 {object} Response
// @Router /api/menu/{id} [get]
func (c *MenuController) GetMenuByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, Response{
			Code: 400,
			Msg:  "无效的菜单ID",
		})
		return
	}

	menu, err := c.menuService.GetMenuByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, Response{
			Code: 404,
			Msg:  "菜单不存在: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, Response{
		Code: 200,
		Msg:  "获取菜单成功",
		Data: menu,
	})
}

// CreateMenu 创建菜单
// @Summary 创建菜单
// @Description 创建新的菜单项
// @Tags 菜单
// @Accept json
// @Produce json
// @Param menu body system.SysMenu true "菜单信息"
// @Security ApiKeyAuth
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Router /api/menu [post]
func (c *MenuController) CreateMenu(ctx *gin.Context) {
	var menu system.SysMenu
	if err := ctx.ShouldBindJSON(&menu); err != nil {
		ctx.JSON(http.StatusBadRequest, Response{
			Code: 400,
			Msg:  "请求参数错误: " + err.Error(),
		})
		return
	}

	if err := c.menuService.CreateMenu(&menu); err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{
			Code: 500,
			Msg:  "创建菜单失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, Response{
		Code: 200,
		Msg:  "创建菜单成功",
		Data: menu,
	})
}

// UpdateMenu 更新菜单
// @Summary 更新菜单
// @Description 更新菜单信息
// @Tags 菜单
// @Accept json
// @Produce json
// @Param id path int true "菜单ID"
// @Param menu body system.SysMenu true "菜单信息"
// @Security ApiKeyAuth
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Failure 404 {object} Response
// @Router /api/menu/{id} [put]
func (c *MenuController) UpdateMenu(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, Response{
			Code: 400,
			Msg:  "无效的菜单ID",
		})
		return
	}

	var menu system.SysMenu
	if err := ctx.ShouldBindJSON(&menu); err != nil {
		ctx.JSON(http.StatusBadRequest, Response{
			Code: 400,
			Msg:  "请求参数错误: " + err.Error(),
		})
		return
	}

	menu.ID = uint(id)
	if err := c.menuService.UpdateMenu(&menu); err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{
			Code: 500,
			Msg:  "更新菜单失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, Response{
		Code: 200,
		Msg:  "更新菜单成功",
		Data: menu,
	})
}

// DeleteMenu 删除菜单
// @Summary 删除菜单
// @Description 删除指定菜单
// @Tags 菜单
// @Accept json
// @Produce json
// @Param id path int true "菜单ID"
// @Security ApiKeyAuth
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Failure 404 {object} Response
// @Router /api/menu/{id} [delete]
func (c *MenuController) DeleteMenu(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, Response{
			Code: 400,
			Msg:  "无效的菜单ID",
		})
		return
	}

	if err := c.menuService.DeleteMenu(uint(id)); err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{
			Code: 500,
			Msg:  "删除菜单失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, Response{
		Code: 200,
		Msg:  "删除菜单成功",
	})
}

// GetMenuRoutes 获取菜单路由配置
// @Summary 获取菜单路由
// @Description 获取用户菜单的路由配置（用于前端路由生成）
// @Tags 菜单
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} Response{data=[]service.MenuRoute}
// @Failure 401 {object} Response
// @Router /api/menu/routes [get]
func (c *MenuController) GetMenuRoutes(ctx *gin.Context) {
	// 从上下文获取用户ID
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, Response{
			Code: 401,
			Msg:  "未授权访问",
		})
		return
	}

	uid, ok := userID.(uint)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, Response{
			Code: 401,
			Msg:  "用户信息错误",
		})
		return
	}

	routes, err := c.menuService.GetMenuRoutes(uid)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{
			Code: 500,
			Msg:  "获取路由失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, Response{
		Code: 200,
		Msg:  "获取路由成功",
		Data: routes,
	})
}

package api

import (
	"log"
	"net/http"
	"strings"
	"time"

	"rentPro/rentpro-admin/common/database"
	"rentPro/rentpro-admin/common/service"

	"github.com/gin-gonic/gin"
)

// AuthController 认证控制器
type AuthController struct {
	authService *service.AuthService
}

// NewAuthController 创建认证控制器实例
func NewAuthController() *AuthController {
	// 从配置获取JWT设置
	jwtSecret := "go-admin" // 应该从配置文件读取
	jwtExpiry := time.Hour  // 1小时过期

	authService := service.NewAuthService(database.DB, jwtSecret, jwtExpiry)

	return &AuthController{
		authService: authService,
	}
}

// Response 通用响应结构
type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

// Login 用户登录
// @Summary 用户登录
// @Description 用户登录接口
// @Tags 认证
// @Accept json
// @Produce json
// @Param request body service.LoginRequest true "登录请求"
// @Success 200 {object} Response{data=service.LoginResponse}
// @Failure 400 {object} Response
// @Router /api/auth/login [post]
func (c *AuthController) Login(ctx *gin.Context) {
	var req service.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, Response{
			Code: 400,
			Msg:  "请求参数错误: " + err.Error(),
		})
		return
	}

	resp, err := c.authService.Login(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, Response{
			Code: 400,
			Msg:  err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, Response{
		Code: 200,
		Msg:  "登录成功",
		Data: resp,
	})
}

// GetUserInfo 获取当前用户信息
// @Summary 获取用户信息
// @Description 获取当前登录用户的详细信息
// @Tags 认证
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} Response{data=service.UserInfoResponse}
// @Failure 401 {object} Response
// @Router /api/auth/user-info [get]
func (c *AuthController) GetUserInfo(ctx *gin.Context) {
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

	userInfo, err := c.authService.GetUserInfo(uid)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{
			Code: 500,
			Msg:  "获取用户信息失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, Response{
		Code: 200,
		Msg:  "获取成功",
		Data: userInfo,
	})
}

// Logout 用户退出
// @Summary 用户退出
// @Description 用户退出登录
// @Tags 认证
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} Response
// @Router /api/auth/logout [post]
func (c *AuthController) Logout(ctx *gin.Context) {
	// 从请求头获取token
	authHeader := ctx.GetHeader("Authorization")
	if authHeader != "" {
		// 提取token
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) == 2 && parts[0] == "Bearer" {
			token := parts[1]
			// 将token加入黑名单
			if err := c.authService.BlacklistToken(token); err != nil {
				log.Printf("将token加入黑名单失败: %v", err)
			} else {
				log.Printf("用户退出，token已加入黑名单: %s", token[:10]+"...")
			}
		}
	}

	ctx.JSON(http.StatusOK, Response{
		Code: 200,
		Msg:  "退出成功",
	})
}

// JWTAuthMiddleware JWT认证中间件
func (c *AuthController) JWTAuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 从请求头获取token
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.JSON(http.StatusUnauthorized, Response{
				Code: 401,
				Msg:  "请求头中auth为空",
			})
			ctx.Abort()
			return
		}

		// 检查token格式
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			ctx.JSON(http.StatusUnauthorized, Response{
				Code: 401,
				Msg:  "请求头中auth格式有误",
			})
			ctx.Abort()
			return
		}

		// 验证token
		claims, err := c.authService.ValidateToken(parts[1])
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, Response{
				Code: 401,
				Msg:  "Token验证失败: " + err.Error(),
			})
			ctx.Abort()
			return
		}

		// 将用户信息存储到上下文
		ctx.Set("user_id", claims.UserID)
		ctx.Set("username", claims.Username)
		ctx.Set("role_id", claims.RoleID)
		ctx.Set("role_key", claims.RoleKey)
		ctx.Set("is_admin", claims.IsAdmin)

		ctx.Next()
	}
}

// PermissionMiddleware 权限验证中间件
func (c *AuthController) PermissionMiddleware(permission string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 获取用户ID
		userID, exists := ctx.Get("user_id")
		if !exists {
			ctx.JSON(http.StatusUnauthorized, Response{
				Code: 401,
				Msg:  "未授权访问",
			})
			ctx.Abort()
			return
		}

		uid, ok := userID.(uint)
		if !ok {
			ctx.JSON(http.StatusUnauthorized, Response{
				Code: 401,
				Msg:  "用户信息错误",
			})
			ctx.Abort()
			return
		}

		// 检查权限
		hasPermission, err := c.authService.HasPermission(uid, permission)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, Response{
				Code: 500,
				Msg:  "权限检查失败: " + err.Error(),
			})
			ctx.Abort()
			return
		}

		if !hasPermission {
			ctx.JSON(http.StatusForbidden, Response{
				Code: 403,
				Msg:  "权限不足",
			})
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}

// AdminMiddleware 管理员权限验证中间件
func (c *AuthController) AdminMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 检查是否为管理员
		isAdmin, exists := ctx.Get("is_admin")
		if !exists || !isAdmin.(bool) {
			ctx.JSON(http.StatusForbidden, Response{
				Code: 403,
				Msg:  "需要管理员权限",
			})
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}

// GetMenus 获取用户菜单
// @Summary 获取用户菜单
// @Description 获取当前用户的菜单列表
// @Tags 认证
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} Response{data=[]service.MenuResponse}
// @Failure 401 {object} Response
// @Router /api/auth/menus [get]
func (c *AuthController) GetMenus(ctx *gin.Context) {
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

	userInfo, err := c.authService.GetUserInfo(uid)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{
			Code: 500,
			Msg:  "获取菜单失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, Response{
		Code: 200,
		Msg:  "获取成功",
		Data: userInfo.Menus,
	})
}

// CheckPermission 检查权限
// @Summary 检查用户权限
// @Description 检查当前用户是否有指定权限
// @Tags 认证
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param permission query string true "权限标识"
// @Success 200 {object} Response{data=bool}
// @Failure 401 {object} Response
// @Router /api/auth/check-permission [get]
func (c *AuthController) CheckPermission(ctx *gin.Context) {
	permission := ctx.Query("permission")
	if permission == "" {
		ctx.JSON(http.StatusBadRequest, Response{
			Code: 400,
			Msg:  "权限标识不能为空",
		})
		return
	}

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

	hasPermission, err := c.authService.HasPermission(uid, permission)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{
			Code: 500,
			Msg:  "权限检查失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, Response{
		Code: 200,
		Msg:  "检查完成",
		Data: hasPermission,
	})
}

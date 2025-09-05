package routes

import (
	"net/http"
	"strings"
	"time"

	"rentPro/rentpro-admin/common/database"
	"rentPro/rentpro-admin/common/models/system"
	"rentPro/rentpro-admin/common/utils"

	"github.com/gin-gonic/gin"
)

// SetupAuthRoutes 设置认证相关路由
func SetupAuthRoutes(api *gin.RouterGroup) {
	// 用户登录
	api.POST("/auth/login", func(c *gin.Context) {
		// 解析请求体
		var loginData struct {
			Username string `json:"username" binding:"required"`
			Password string `json:"password" binding:"required"`
		}

		if err := c.ShouldBindJSON(&loginData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "请求参数错误",
				"error":   err.Error(),
			})
			return
		}

		// 验证用户凭据
		var user system.SysUser
		result := database.DB.Where("username = ?", loginData.Username).First(&user)
		if result.Error != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "用户名或密码错误",
			})
			return
		}

		// 验证密码
		if !user.ComparePassword(loginData.Password) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "用户名或密码错误",
			})
			return
		}

		// 检查用户状态
		if !user.IsActive() {
			c.JSON(http.StatusForbidden, gin.H{
				"code":    403,
				"message": "用户已被禁用",
			})
			return
		}

		// 生成JWT token
		jwtConfig := utils.JWTConfig{
			Secret:  "rentpro-admin-secret-key",
			Timeout: 86400, // 24 hours in seconds
		}
		jwtInstance := &utils.JWT{
			Config: jwtConfig,
		}
		token, err := jwtInstance.GenerateToken(user.ID, user.Username)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "生成token失败",
				"error":   err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "登录成功",
			"data": gin.H{
				"token": token,
				"user": gin.H{
					"id":       user.ID,
					"username": user.Username,
					"nickname": user.NickName,
					"role_id":  user.RoleID,
				},
			},
		})
	})

	// 用户退出
	api.POST("/auth/logout", func(c *gin.Context) {
		// 退出登录逻辑：前端清除token即可，后端无需特殊处理
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "退出成功",
		})
	})

	// 获取用户信息
	api.GET("/auth/userinfo", func(c *gin.Context) {
		// 从请求头获取token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "未提供认证信息",
			})
			return
		}

		// 提取token
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "认证格式错误",
			})
			return
		}

		// 解析token
		jwtInstance := &utils.JWT{
			Config: utils.JWTConfig{
				Secret:  "rentpro-admin-secret-key",
				Timeout: 86400, // 24 hours in seconds
			},
		}
		claims, err := jwtInstance.ParseToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "token无效",
				"error":   err.Error(),
			})
			return
		}

		// 获取用户信息
		var user system.SysUser
		result := database.DB.Where("id = ?", claims.UserID).First(&user)
		if result.Error != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"code":    404,
				"message": "用户不存在",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "获取成功",
			"data": gin.H{
				"id":       user.ID,
				"username": user.Username,
				"nickname": user.NickName,
				"role_id":  user.RoleID,
			},
		})
	})

	// 验证token
	api.GET("/auth/check", func(c *gin.Context) {
		// 从请求头获取token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "未提供认证信息",
			})
			return
		}

		// 提取token
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "认证格式错误",
			})
			return
		}

		// 解析token
		jwtInstance := &utils.JWT{
			Config: utils.JWTConfig{
				Secret:  "rentpro-admin-secret-key",
				Timeout: 86400, // 24 hours in seconds
			},
		}
		claims, err := jwtInstance.ParseToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "token无效",
				"error":   err.Error(),
			})
			return
		}

		// 检查token是否过期
		if claims.ExpiresAt != nil && time.Now().Unix() > claims.ExpiresAt.Unix() {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "token已过期",
			})
			return
		}

		// 设置用户信息到上下文
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)

		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "token有效",
			"data": gin.H{
				"user_id":  claims.UserID,
				"username": claims.Username,
			},
		})
	})
}

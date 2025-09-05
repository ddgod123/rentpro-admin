package routes

import (
	"net/http"

	"rentPro/rentpro-admin/common/database"

	"github.com/gin-gonic/gin"
)

// SetupUserRoutes 设置用户管理相关路由
func SetupUserRoutes(api *gin.RouterGroup) {
	// 获取用户列表
	api.GET("/users", func(c *gin.Context) {
		// 从数据库查询用户列表
		var users []map[string]interface{}
		result := database.DB.Table("sys_users").Where("deleted_at IS NULL").Find(&users)

		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "查询用户列表失败",
				"error":   result.Error.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "获取用户列表成功",
			"data":    users,
		})
	})

	// 获取单个用户信息
	api.GET("/users/:id", func(c *gin.Context) {
		id := c.Param("id")

		var user map[string]interface{}
		result := database.DB.Table("sys_users").Where("id = ? AND deleted_at IS NULL", id).First(&user)

		if result.Error != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"code":    404,
				"message": "用户不存在",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "获取用户信息成功",
			"data":    user,
		})
	})

	// 创建用户
	api.POST("/users", func(c *gin.Context) {
		var userData struct {
			Username string `json:"username" binding:"required"`
			Password string `json:"password" binding:"required"`
			Nickname string `json:"nickname"`
			RoleID   uint64 `json:"roleId"`
			Status   string `json:"status"`
		}

		if err := c.ShouldBindJSON(&userData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "请求参数错误",
				"error":   err.Error(),
			})
			return
		}

		// 检查用户名是否已存在
		var count int64
		database.DB.Table("sys_users").Where("username = ?", userData.Username).Count(&count)
		if count > 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "用户名已存在",
			})
			return
		}

		// 插入数据库
		result := database.DB.Exec(
			"INSERT INTO sys_users (username, password, nickname, role_id, status, created_at, updated_at) VALUES (?, ?, ?, ?, ?, NOW(), NOW())",
			userData.Username,
			userData.Password,
			userData.Nickname,
			userData.RoleID,
			userData.Status,
		)

		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "创建用户失败",
				"error":   result.Error.Error(),
			})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"code":    201,
			"message": "创建用户成功",
		})
	})

	// 更新用户信息
	api.PUT("/users/:id", func(c *gin.Context) {
		id := c.Param("id")

		var userData struct {
			Username string `json:"username"`
			Nickname string `json:"nickname"`
			RoleID   uint64 `json:"roleId"`
			Status   string `json:"status"`
		}

		if err := c.ShouldBindJSON(&userData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "请求参数错误",
				"error":   err.Error(),
			})
			return
		}

		// 更新数据库
		updateData := make(map[string]interface{})
		if userData.Username != "" {
			updateData["username"] = userData.Username
		}
		if userData.Nickname != "" {
			updateData["nickname"] = userData.Nickname
		}
		if userData.RoleID > 0 {
			updateData["role_id"] = userData.RoleID
		}
		if userData.Status != "" {
			updateData["status"] = userData.Status
		}

		result := database.DB.Table("sys_users").Where("id = ?", id).Updates(updateData)

		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "更新用户失败",
				"error":   result.Error.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "更新用户成功",
		})
	})

	// 删除用户
	api.DELETE("/users/:id", func(c *gin.Context) {
		id := c.Param("id")

		result := database.DB.Exec("UPDATE sys_users SET deleted_at = NOW() WHERE id = ?", id)

		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "删除用户失败",
				"error":   result.Error.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "删除用户成功",
		})
	})
}

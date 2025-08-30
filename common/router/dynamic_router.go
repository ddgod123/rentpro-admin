package router

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"rentPro/rentpro-admin/common/api"
	"rentPro/rentpro-admin/common/database"
	"rentPro/rentpro-admin/common/models/rental"
	"rentPro/rentpro-admin/common/models/system"
	"rentPro/rentpro-admin/common/service"

	"github.com/gin-gonic/gin"
)

// DynamicRouter 动态路由生成器
type DynamicRouter struct {
	menuService    *service.MenuService
	authController *api.AuthController
}

// NewDynamicRouter 创建动态路由生成器实例
func NewDynamicRouter(menuService *service.MenuService, authController *api.AuthController) *DynamicRouter {
	return &DynamicRouter{
		menuService:    menuService,
		authController: authController,
	}
}

// RegisterDynamicRoutes 注册动态路由
func (r *DynamicRouter) RegisterDynamicRoutes(router *gin.Engine) {
	// 创建菜单控制器
	menuController := api.NewMenuController()

	// API根路径分组
	apiGroup := router.Group("/api")

	// 菜单管理路由组
	menuGroup := apiGroup.Group("/menu")
	menuGroup.Use(r.authController.JWTAuthMiddleware()) // 需要JWT认证
	{
		// 获取用户菜单
		menuGroup.GET("/user-menus", menuController.GetUserMenus)
		// 获取所有菜单（管理员用）
		menuGroup.GET("/all", menuController.GetAllMenus)
		// 获取菜单路由配置
		menuGroup.GET("/routes", menuController.GetMenuRoutes)
		// 根据ID获取菜单
		menuGroup.GET("/:id", menuController.GetMenuByID)
		// 创建菜单
		menuGroup.POST("", menuController.CreateMenu)
		// 更新菜单
		menuGroup.PUT("/:id", menuController.UpdateMenu)
		// 删除菜单
		menuGroup.DELETE("/:id", menuController.DeleteMenu)
	}

	// 添加测试路由
	apiGroup.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "动态路由测试成功"})
	})

	// 动态生成业务路由
	r.generateBusinessRoutes(apiGroup)
}

// generateBusinessRoutes 生成业务路由
func (r *DynamicRouter) generateBusinessRoutes(apiGroup *gin.RouterGroup) {
	// 获取所有菜单
	menus, err := r.menuService.GetAllMenus()
	if err != nil {
		fmt.Printf("获取菜单失败: %v\n", err)
		return
	}

	fmt.Printf("=== 动态路由生成开始 ===\n")
	fmt.Printf("获取到 %d 个菜单\n", len(menus))

	// 为每个菜单生成对应的路由
	for _, menu := range menus {
		fmt.Printf("处理菜单: %s, 路径: %s, 权限: %s\n", menu.Name, menu.Path, menu.Permission)
		r.generateMenuRoutes(apiGroup, menu)
	}

	fmt.Printf("=== 动态路由生成完成 ===\n")
}

// generateMenuRoutes 为单个菜单生成路由
func (r *DynamicRouter) generateMenuRoutes(apiGroup *gin.RouterGroup, menu system.SysMenu) {
	// 只处理有权限标识的菜单
	if menu.Permission == "" {
		return
	}

	// 根据菜单路径生成路由组
	path := strings.TrimPrefix(menu.Path, "/")
	if path == "" {
		return
	}

	// 创建路由组
	routeGroup := apiGroup.Group(path)
	routeGroup.Use(r.authController.JWTAuthMiddleware()) // 需要JWT认证

	// 如果有权限标识，添加权限中间件
	if menu.Permission != "" {
		routeGroup.Use(r.authController.PermissionMiddleware(menu.Permission))
	}

	// 根据菜单类型生成不同的路由
	switch menu.Name {
	case "Dashboard":
		r.generateDashboardRoutes(routeGroup)
	case "User":
		r.generateUserRoutes(routeGroup)
	case "Role":
		r.generateRoleRoutes(routeGroup)
	case "Building":
		r.generateBuildingRoutes(routeGroup)
	case "Room":
		r.generateRoomRoutes(routeGroup)
	case "Tenant":
		r.generateTenantRoutes(routeGroup)
	case "Contract":
		r.generateContractRoutes(routeGroup)
	default:
		// 默认路由处理
		r.generateDefaultRoutes(routeGroup, menu)
	}
}

// generateDashboardRoutes 生成仪表板路由
func (r *DynamicRouter) generateDashboardRoutes(routeGroup *gin.RouterGroup) {
	routeGroup.GET("/stats", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "获取仪表板统计数据",
			"data": gin.H{
				"totalUsers":     100,
				"totalBuildings": 50,
				"totalRooms":     200,
				"totalContracts": 150,
			},
		})
	})

	routeGroup.GET("/chart", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "获取图表数据",
			"data": gin.H{
				"monthlyRent": []gin.H{
					{"month": "1月", "amount": 50000},
					{"month": "2月", "amount": 55000},
					{"month": "3月", "amount": 60000},
				},
			},
		})
	})
}

// generateUserRoutes 生成用户管理路由
func (r *DynamicRouter) generateUserRoutes(routeGroup *gin.RouterGroup) {
	routeGroup.GET("/list", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "获取用户列表",
			"data": []gin.H{
				{"id": 1, "username": "admin", "nickname": "管理员", "email": "admin@example.com"},
				{"id": 2, "username": "user1", "nickname": "用户1", "email": "user1@example.com"},
			},
		})
	})

	routeGroup.POST("/create", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "创建用户成功",
		})
	})

	routeGroup.PUT("/:id", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "更新用户成功",
		})
	})

	routeGroup.DELETE("/:id", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "删除用户成功",
		})
	})
}

// generateRoleRoutes 生成角色管理路由
func (r *DynamicRouter) generateRoleRoutes(routeGroup *gin.RouterGroup) {
	routeGroup.GET("/list", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "获取角色列表",
			"data": []gin.H{
				{"id": 1, "name": "超级管理员", "code": "super_admin"},
				{"id": 2, "name": "普通用户", "code": "user"},
			},
		})
	})

	routeGroup.POST("/create", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "创建角色成功",
		})
	})

	routeGroup.PUT("/:id", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "更新角色成功",
		})
	})

	routeGroup.DELETE("/:id", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "删除角色成功",
		})
	})
}

// generateBuildingRoutes 生成楼盘管理路由
func (r *DynamicRouter) generateBuildingRoutes(routeGroup *gin.RouterGroup) {
	// 获取楼盘列表（支持分页和搜索）
	routeGroup.GET("", func(c *gin.Context) {
		// 获取查询参数
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
		name := c.Query("name")
		city := c.Query("city")
		status := c.Query("status")

		// 构建查询条件
		query := database.DB.Model(&rental.SysBuildings{})

		if name != "" {
			query = query.Where("name LIKE ?", "%"+name+"%")
		}
		if city != "" {
			query = query.Where("city = ?", city)
		}
		if status != "" {
			query = query.Where("status = ?", status)
		}

		// 获取总数
		var total int64
		query.Count(&total)

		// 分页查询
		var buildings []rental.SysBuildings
		offset := (page - 1) * pageSize
		if err := query.Offset(offset).Limit(pageSize).Find(&buildings).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code": 500,
				"msg":  "获取楼盘列表失败: " + err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "获取楼盘列表成功",
			"data": gin.H{
				"list":     buildings,
				"total":    total,
				"page":     page,
				"pageSize": pageSize,
			},
		})
	})

	// 获取楼盘列表（简化版，用于下拉选择等）
	routeGroup.GET("/list", func(c *gin.Context) {
		var buildings []rental.SysBuildings
		if err := database.DB.Find(&buildings).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code": 500,
				"msg":  "获取楼盘列表失败: " + err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "获取楼盘列表成功",
			"data": buildings,
		})
	})

	routeGroup.POST("/create", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "创建楼盘成功",
		})
	})

	routeGroup.PUT("/:id", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "更新楼盘成功",
		})
	})

	routeGroup.DELETE("/:id", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "删除楼盘成功",
		})
	})
}

// generateRoomRoutes 生成房源管理路由
func (r *DynamicRouter) generateRoomRoutes(routeGroup *gin.RouterGroup) {
	routeGroup.GET("/list", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "获取房源列表",
			"data": []gin.H{
				{"id": 1, "roomNumber": "A101", "building": "阳光花园", "area": 80, "status": "已租"},
				{"id": 2, "roomNumber": "A102", "building": "阳光花园", "area": 90, "status": "空置"},
			},
		})
	})

	routeGroup.POST("/create", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "创建房源成功",
		})
	})

	routeGroup.PUT("/:id", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "更新房源成功",
		})
	})

	routeGroup.DELETE("/:id", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "删除房源成功",
		})
	})
}

// generateTenantRoutes 生成租客管理路由
func (r *DynamicRouter) generateTenantRoutes(routeGroup *gin.RouterGroup) {
	routeGroup.GET("/list", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "获取租客列表",
			"data": []gin.H{
				{"id": 1, "name": "张三", "phone": "13800138001", "room": "A101"},
				{"id": 2, "name": "李四", "phone": "13800138002", "room": "B201"},
			},
		})
	})

	routeGroup.POST("/create", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "创建租客成功",
		})
	})

	routeGroup.PUT("/:id", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "更新租客成功",
		})
	})

	routeGroup.DELETE("/:id", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "删除租客成功",
		})
	})
}

// generateContractRoutes 生成合同管理路由
func (r *DynamicRouter) generateContractRoutes(routeGroup *gin.RouterGroup) {
	routeGroup.GET("/list", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "获取合同列表",
			"data": []gin.H{
				{"id": 1, "contractNo": "CT001", "tenant": "张三", "room": "A101", "startDate": "2024-01-01", "endDate": "2024-12-31"},
				{"id": 2, "contractNo": "CT002", "tenant": "李四", "room": "B201", "startDate": "2024-02-01", "endDate": "2025-01-31"},
			},
		})
	})

	routeGroup.POST("/create", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "创建合同成功",
		})
	})

	routeGroup.PUT("/:id", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "更新合同成功",
		})
	})

	routeGroup.DELETE("/:id", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "删除合同成功",
		})
	})
}

// generateDefaultRoutes 生成默认路由
func (r *DynamicRouter) generateDefaultRoutes(routeGroup *gin.RouterGroup, menu system.SysMenu) {
	routeGroup.GET("/list", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  fmt.Sprintf("获取%s列表", menu.Title),
			"data": []gin.H{},
		})
	})

	routeGroup.POST("/create", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  fmt.Sprintf("创建%s成功", menu.Title),
		})
	})

	routeGroup.PUT("/:id", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  fmt.Sprintf("更新%s成功", menu.Title),
		})
	})

	routeGroup.DELETE("/:id", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  fmt.Sprintf("删除%s成功", menu.Title),
		})
	})
}
